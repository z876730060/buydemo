package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/middlewares"
	"github.com/z876730060/buydemo/models"
)

func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query := database.DB.Model(&models.Product{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR spec LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query.Count(&total)

	var products []models.Product
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	if product.Code == "" || product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "编码和名称不能为空"})
		return
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败，编码可能已存在"})
		return
	}

	// Auto create inventory record
	var count int64
	database.DB.Model(&models.Inventory{}).Where("product_id = ?", product.ID).Count(&count)
	if count == 0 {
		database.DB.Create(&models.Inventory{
			ProductID: product.ID,
			Quantity:  0,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": product, "message": "创建成功"})
}

func UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}

	database.DB.Model(&product).Updates(map[string]interface{}{
		"code":     input.Code,
		"name":     input.Name,
		"category": input.Category,
		"unit":     input.Unit,
		"spec":     input.Spec,
		"price":    input.Price,
		"status":   input.Status,
	})

	middlewares.SimpleLog(c, "update", "product", product.ID, "编辑商品: "+product.Name)
	c.JSON(http.StatusOK, gin.H{"data": product, "message": "更新成功"})
}

func DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	database.DB.Delete(&product)
	database.DB.Where("product_id = ?", product.ID).Delete(&models.Inventory{})
	middlewares.SimpleLog(c, "delete", "product", product.ID, "删除商品: "+product.Name)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetAllProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Where("status = 1").Order("name ASC").Find(&products)
	c.JSON(http.StatusOK, gin.H{"data": products})
}

func GetCategories(c *gin.Context) {
	var categories []string
	database.DB.Model(&models.Product{}).Distinct("category").Where("category != ''").Pluck("category", &categories)
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func GetProductDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	// Get inventory
	var inv models.Inventory
	database.DB.Where("product_id = ?", product.ID).First(&inv)

	// Get recent purchase items
	type PurchaseItemVO struct {
		OrderNo    string  `json:"order_no"`
		Quantity   float64 `json:"quantity"`
		UnitPrice  float64 `json:"unit_price"`
		Status     string  `json:"status"`
		CreatedAt  string  `json:"created_at"`
	}
	var purchaseItems []PurchaseItemVO
	database.DB.Model(&models.PurchaseOrderItem{}).
		Select("purchase_orders.order_no, purchase_order_items.quantity, purchase_order_items.unit_price, purchase_orders.status, purchase_orders.created_at").
		Joins("JOIN purchase_orders ON purchase_orders.id = purchase_order_items.order_id").
		Where("purchase_order_items.product_id = ?", product.ID).
		Order("purchase_orders.id DESC").
		Limit(20).
		Find(&purchaseItems)

	// Get recent sales items
	type SalesItemVO struct {
		OrderNo   string  `json:"order_no"`
		Quantity  float64 `json:"quantity"`
		UnitPrice float64 `json:"unit_price"`
		Status    string  `json:"status"`
		CreatedAt string  `json:"created_at"`
	}
	var salesItems []SalesItemVO
	database.DB.Model(&models.SalesOrderItem{}).
		Select("sales_orders.order_no, sales_order_items.quantity, sales_order_items.unit_price, sales_orders.status, sales_orders.created_at").
		Joins("JOIN sales_orders ON sales_orders.id = sales_order_items.order_id").
		Where("sales_order_items.product_id = ?", product.ID).
		Order("sales_orders.id DESC").
		Limit(20).
		Find(&salesItems)

	// Get inventory logs
	type LogVO struct {
		Type      string  `json:"type"`
		Quantity  float64 `json:"quantity"`
		Balance   float64 `json:"balance"`
		RefType   string  `json:"reference_type"`
		Remark    string  `json:"remark"`
		CreatedAt string  `json:"created_at"`
	}
	var logs []LogVO
	database.DB.Model(&models.InventoryLog{}).
		Select("type, quantity, balance, reference_type, remark, created_at").
		Where("product_id = ?", product.ID).
		Order("id DESC").
		Limit(50).
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"product":       product,
			"inventory":     inv,
			"purchase_items": purchaseItems,
			"sales_items":   salesItems,
			"inventory_logs": logs,
		},
	})
}
