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

type CreateOrderRequest struct {
	SupplierID uint                   `json:"supplier_id" binding:"required"`
	Remark     string                 `json:"remark"`
	Items      []CreateOrderItemReq   `json:"items" binding:"required,min=1"`
}

type CreateOrderItemReq struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"required,gte=0"`
}

func generateOrderNo() string {
	now := time.Now()
	var count int64
	database.DB.Model(&models.PurchaseOrder{}).Where("created_at >= ? AND created_at < ?",
		now.Format("2006-01-02"), now.AddDate(0, 0, 1).Format("2006-01-02")).Count(&count)
	return fmt.Sprintf("PO%s%04d", now.Format("20060102"), count+1)
}

func GetPurchaseOrders(c *gin.Context) {
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
	query := database.DB.Model(&models.PurchaseOrder{}).Preload("Supplier")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Joins("JOIN suppliers ON suppliers.id = purchase_orders.supplier_id").
			Where("purchase_orders.order_no LIKE ? OR suppliers.name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var orders []models.PurchaseOrder
	query.Preload("Supplier").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"data":  orders,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetPurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.PurchaseOrder
	if err := database.DB.Preload("Supplier").Preload("Creator").
		Preload("Items.Product").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "采购单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func CreatePurchaseOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	// Verify supplier
	var supplier models.Supplier
	if err := database.DB.First(&supplier, req.SupplierID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "供应商不存在"})
		return
	}

	userID := c.GetUint("user_id")

	// Calculate total
	var totalAmount float64
	var items []models.PurchaseOrderItem
	for _, item := range req.Items {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("商品ID %d 不存在", item.ProductID)})
			return
		}
		amount := item.Quantity * item.UnitPrice
		totalAmount += amount
		items = append(items, models.PurchaseOrderItem{
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      amount,
			ReceivedQty: 0,
		})
	}

	order := models.PurchaseOrder{
		OrderNo:     generateOrderNo(),
		SupplierID:  req.SupplierID,
		TotalAmount: totalAmount,
		Status:      "draft",
		Remark:      req.Remark,
		CreatedBy:   userID,
		Items:       items,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建采购单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "创建成功"})
}

func UpdatePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.PurchaseOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "采购单不存在"})
		return
	}

	if order.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能编辑草稿状态的采购单"})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	// Delete old items
	database.DB.Where("order_id = ?", order.ID).Delete(&models.PurchaseOrderItem{})

	// Recreate items
	var totalAmount float64
	for _, item := range req.Items {
		amount := item.Quantity * item.UnitPrice
		totalAmount += amount
		database.DB.Create(&models.PurchaseOrderItem{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			Amount:     amount,
			ReceivedQty: 0,
		})
	}

	database.DB.Model(&order).Updates(map[string]interface{}{
		"supplier_id":  req.SupplierID,
		"total_amount": totalAmount,
		"remark":       req.Remark,
	})

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "更新成功"})
}

func ApproveOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.PurchaseOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "采购单不存在"})
		return
	}

	if order.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能审核草稿状态的采购单"})
		return
	}

	now := time.Now()
	database.DB.Model(&order).Updates(map[string]interface{}{
		"status":      "approved",
		"approved_at": &now,
	})

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "审核通过"})
}

func ReceiveOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.PurchaseOrder
	if err := database.DB.Preload("Items.Product").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "采购单不存在"})
		return
	}

	if order.Status != "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能入库已审核的采购单"})
		return
	}

	// Update received quantity and inventory
	tx := database.DB.Begin()

	for _, item := range order.Items {
		remaining := item.Quantity - item.ReceivedQty
		if remaining <= 0 {
			continue
		}

		// Update item received qty
		tx.Model(&item).Update("received_qty", item.Quantity)

		// Update inventory
		var inv models.Inventory
		if err := tx.Where("product_id = ?", item.ProductID).First(&inv).Error; err != nil {
			tx.Create(&models.Inventory{
				ProductID: item.ProductID,
				Quantity:  remaining,
			})
		} else {
			tx.Model(&inv).Update("quantity", inv.Quantity+remaining)
		}

		// Create inventory log
		var invAfter models.Inventory
		tx.Where("product_id = ?", item.ProductID).First(&invAfter)

		tx.Create(&models.InventoryLog{
			ProductID:     item.ProductID,
			Type:          "in",
			Quantity:      remaining,
			Balance:       invAfter.Quantity,
			ReferenceType: "purchase",
			ReferenceID:   order.ID,
			Remark:        fmt.Sprintf("采购入库：%s", order.OrderNo),
		})
	}

	now := time.Now()
	tx.Model(&order).Updates(map[string]interface{}{
		"status":      "received",
		"received_at": &now,
	})

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "入库成功"})
}

func CancelOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var order models.PurchaseOrder
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "采购单不存在"})
		return
	}

	if order.Status != "draft" && order.Status != "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该状态的采购单无法取消"})
		return
	}

	database.DB.Model(&order).Update("status", "cancelled")
	c.JSON(http.StatusOK, gin.H{"data": order, "message": "已取消"})
}
