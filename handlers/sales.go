package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/models"
)

type CreateSalesOrderRequest struct {
	CustomerID uint                     `json:"customer_id" binding:"required"`
	Remark     string                   `json:"remark"`
	Items      []CreateSalesOrderItemReq `json:"items" binding:"required,min=1"`
}

type CreateSalesOrderItemReq struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"required,gte=0"`
}

func generateSalesOrderNo() string {
	now := time.Now()
	var count int64
	database.DB.Model(&models.SalesOrder{}).Where("created_at >= ? AND created_at < ?",
		now.Format("2006-01-02"), now.AddDate(0, 0, 1).Format("2006-01-02")).Count(&count)
	return fmt.Sprintf("SO%s%04d", now.Format("20060102"), count+1)
}

func GetSalesOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query := database.DB.Model(&models.SalesOrder{}).Preload("Customer")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Joins("JOIN customers ON customers.id = sales_orders.customer_id").
			Where("sales_orders.order_no LIKE ? OR customers.name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var orders []models.SalesOrder
	query.Preload("Customer").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"data":  orders,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.SalesOrder
	if err := database.DB.Preload("Customer").Preload("Creator").
		Preload("Items.Product").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "销售单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func CreateSalesOrder(c *gin.Context) {
	var req CreateSalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	// Verify customer
	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "客户不存在"})
		return
	}

	userID := c.GetUint("user_id")

	// Calculate total
	var totalAmount float64
	var items []models.SalesOrderItem
	for _, item := range req.Items {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("商品ID %d 不存在", item.ProductID)})
			return
		}
		amount := item.Quantity * item.UnitPrice
		totalAmount += amount
		items = append(items, models.SalesOrderItem{
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      amount,
			DeliveredQty: 0,
		})
	}

	order := models.SalesOrder{
		OrderNo:     generateSalesOrderNo(),
		CustomerID:  req.CustomerID,
		TotalAmount: totalAmount,
		Status:      "draft",
		Remark:      req.Remark,
		CreatedBy:   userID,
		Items:       items,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建销售单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "创建成功"})
}

func UpdateSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.SalesOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "销售单不存在"})
		return
	}

	if order.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能编辑草稿状态的销售单"})
		return
	}

	var req CreateSalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	// Delete old items
	database.DB.Where("order_id = ?", order.ID).Delete(&models.SalesOrderItem{})

	// Recreate items
	var totalAmount float64
	for _, item := range req.Items {
		amount := item.Quantity * item.UnitPrice
		totalAmount += amount
		database.DB.Create(&models.SalesOrderItem{
			OrderID:     order.ID,
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      amount,
			DeliveredQty: 0,
		})
	}

	database.DB.Model(&order).Updates(map[string]interface{}{
		"customer_id":  req.CustomerID,
		"total_amount": totalAmount,
		"remark":       req.Remark,
	})

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "更新成功"})
}

func ApproveSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.SalesOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "销售单不存在"})
		return
	}

	if order.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能审核草稿状态的销售单"})
		return
	}

	now := time.Now()
	database.DB.Model(&order).Updates(map[string]interface{}{
		"status":      "approved",
		"approved_at": &now,
	})

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "审核通过"})
}

func DeliverSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.SalesOrder
	if err := database.DB.Preload("Items.Product").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "销售单不存在"})
		return
	}

	if order.Status != "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能出库已审核的销售单"})
		return
	}

	// Check inventory sufficiency
	for _, item := range order.Items {
		remaining := item.Quantity - item.DeliveredQty
		if remaining <= 0 {
			continue
		}
		var inv models.Inventory
		if err := database.DB.Where("product_id = ?", item.ProductID).First(&inv).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("商品 %s 无库存记录", item.Product.Name)})
			return
		}
		if inv.Quantity < remaining {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("商品 %s 库存不足（当前: %.2f, 需出库: %.2f）",
					item.Product.Name, inv.Quantity, remaining),
			})
			return
		}
	}

	// Process delivery
	tx := database.DB.Begin()

	for _, item := range order.Items {
		remaining := item.Quantity - item.DeliveredQty
		if remaining <= 0 {
			continue
		}

		// Update item delivered qty
		tx.Model(&item).Update("delivered_qty", item.Quantity)

		// Reduce inventory
		var inv models.Inventory
		tx.Where("product_id = ?", item.ProductID).First(&inv)
		tx.Model(&inv).Update("quantity", inv.Quantity-remaining)

		// Create inventory log
		var invAfter models.Inventory
		tx.Where("product_id = ?", item.ProductID).First(&invAfter)

		tx.Create(&models.InventoryLog{
			ProductID:     item.ProductID,
			Type:          "out",
			Quantity:      remaining,
			Balance:       invAfter.Quantity,
			ReferenceType: "sales",
			ReferenceID:   order.ID,
			Remark:        fmt.Sprintf("销售出库：%s", order.OrderNo),
		})
	}

	now := time.Now()
	tx.Model(&order).Updates(map[string]interface{}{
		"status":       "delivered",
		"delivered_at": &now,
	})

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "出库成功"})
}

func CancelSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.SalesOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "销售单不存在"})
		return
	}

	if order.Status != "draft" && order.Status != "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该状态的销售单无法取消"})
		return
	}

	database.DB.Model(&order).Update("status", "cancelled")
	c.JSON(http.StatusOK, gin.H{"data": order, "message": "已取消"})
}
