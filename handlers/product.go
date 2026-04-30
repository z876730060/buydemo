package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
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
