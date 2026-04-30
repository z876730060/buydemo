package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/middlewares"
	"github.com/z876730060/buydemo/models"
)

func GetInventories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query := database.DB.Model(&models.Inventory{}).Joins("Product")
	if keyword != "" {
		query = query.Where("Product.name LIKE ? OR Product.code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var inventories []models.Inventory
	query.Preload("Product").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&inventories)

	c.JSON(http.StatusOK, gin.H{
		"data":  inventories,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetInventoryLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	productID := c.Query("product_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query := database.DB.Model(&models.InventoryLog{}).Preload("Product")
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	query.Count(&total)

	var logs []models.InventoryLog
	query.Preload("Product").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetLowStock(c *gin.Context) {
	threshold, _ := strconv.ParseFloat(c.DefaultQuery("threshold", "10"), 64)
	if threshold <= 0 {
		threshold = 10
	}

	var inventories []models.Inventory
	database.DB.Preload("Product").
		Where("quantity <= ?", threshold).
		Where("quantity > 0").
		Order("quantity ASC").
		Find(&inventories)

	c.JSON(http.StatusOK, gin.H{"data": inventories})
}

func GetDashboardStats(c *gin.Context) {
	// Total suppliers
	var supplierCount int64
	database.DB.Model(&models.Supplier{}).Where("status = 1").Count(&supplierCount)

	// Total customers
	var customerCount int64
	database.DB.Model(&models.Customer{}).Where("status = 1").Count(&customerCount)

	// Total products
	var productCount int64
	database.DB.Model(&models.Product{}).Where("status = 1").Count(&productCount)

	// Today's purchase orders
	var todayOrderCount int64
	database.DB.Model(&models.PurchaseOrder{}).Where("created_at >= date('now')").Count(&todayOrderCount)

	// Low stock count
	var lowStockCount int64
	database.DB.Model(&models.Inventory{}).Where("quantity <= 10 AND quantity > 0").Count(&lowStockCount)

	// Out of stock count
	var outOfStockCount int64
	database.DB.Model(&models.Inventory{}).Where("quantity <= 0").Count(&outOfStockCount)

	// Purchase orders by status
	type StatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var statusCounts []StatusCount
	database.DB.Model(&models.PurchaseOrder{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&statusCounts)

	// Sales orders by status
	var salesStatusCounts []StatusCount
	database.DB.Model(&models.SalesOrder{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&salesStatusCounts)

	// Today's sales orders
	var todaySalesCount int64
	database.DB.Model(&models.SalesOrder{}).Where("created_at >= date('now')").Count(&todaySalesCount)

	// Total accounts summary
	var apAmount float64
	database.DB.Model(&models.AccountPayable{}).
		Select("COALESCE(SUM(due_amount),0)").
		Where("status != 'paid'").
		Scan(&apAmount)

	var arAmount float64
	database.DB.Model(&models.AccountReceivable{}).
		Select("COALESCE(SUM(due_amount),0)").
		Where("status != 'received'").
		Scan(&arAmount)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"supplier_count":      supplierCount,
			"customer_count":      customerCount,
			"product_count":       productCount,
			"today_order_count":   todayOrderCount,
			"today_sales_count":   todaySalesCount,
			"low_stock_count":     lowStockCount,
			"out_of_stock_count":  outOfStockCount,
			"order_status":        statusCounts,
			"sales_order_status":  salesStatusCounts,
			"ap_amount":           apAmount,
			"ar_amount":           arAmount,
		},
	})
}

// AdjustInventory manually adjusts inventory quantity
func AdjustInventory(c *gin.Context) {
	var req struct {
		ProductID uint    `json:"product_id" binding:"required"`
		Quantity  float64 `json:"quantity" binding:"required"` // positive = add, negative = remove
		Remark    string  `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	if req.Quantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "调整数量不能为0"})
		return
	}

	adjType := "in"
	if req.Quantity < 0 {
		adjType = "out"
	}

	tx := database.DB.Begin()

	var inv models.Inventory
	if err := tx.Where("product_id = ?", req.ProductID).First(&inv).Error; err != nil {
		// Create inventory if not exists
		if req.Quantity < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "库存不足，无法出库"})
			return
		}
		inv = models.Inventory{ProductID: req.ProductID, Quantity: req.Quantity}
		tx.Create(&inv)
	} else {
		if inv.Quantity+req.Quantity < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "库存不足，调整后库存为负"})
			return
		}
		tx.Model(&inv).Update("quantity", inv.Quantity+req.Quantity)
	}

	// Reload to get updated balance
	var invAfter models.Inventory
	tx.Where("product_id = ?", req.ProductID).First(&invAfter)

	// Create inventory log
	tx.Create(&models.InventoryLog{
		ProductID:     req.ProductID,
		Type:          adjType,
		Quantity:      req.Quantity,
		Balance:       invAfter.Quantity,
		ReferenceType: "adjust",
		Remark:        req.Remark,
	})

	tx.Commit()

	middlewares.SimpleLog(c, "adjust", "inventory", req.ProductID, "手动调整库存: "+product.Name+" 数量:"+fmt.Sprintf("%.2f", req.Quantity))
	c.JSON(http.StatusOK, gin.H{"data": invAfter, "message": "库存调整成功"})
}
