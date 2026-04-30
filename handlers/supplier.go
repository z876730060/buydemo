package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/models"
)

func GetSuppliers(c *gin.Context) {
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
	query := database.DB.Model(&models.Supplier{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR contact_person LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var suppliers []models.Supplier
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&suppliers)

	c.JSON(http.StatusOK, gin.H{
		"data":  suppliers,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var supplier models.Supplier
	if err := database.DB.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": supplier})
}

func CreateSupplier(c *gin.Context) {
	var supplier models.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	if supplier.Code == "" || supplier.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "编码和名称不能为空"})
		return
	}

	if err := database.DB.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败，编码可能已存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": supplier, "message": "创建成功"})
}

func UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var supplier models.Supplier
	if err := database.DB.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
		return
	}

	var input models.Supplier
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}

	database.DB.Model(&supplier).Updates(map[string]interface{}{
		"code":           input.Code,
		"name":           input.Name,
		"contact_person": input.ContactPerson,
		"phone":          input.Phone,
		"address":        input.Address,
		"status":         input.Status,
	})

	c.JSON(http.StatusOK, gin.H{"data": supplier, "message": "更新成功"})
}

func DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var supplier models.Supplier
	if err := database.DB.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
		return
	}

	database.DB.Delete(&supplier)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetAllSuppliers(c *gin.Context) {
	var suppliers []models.Supplier
	database.DB.Where("status = 1").Order("name ASC").Find(&suppliers)
	c.JSON(http.StatusOK, gin.H{"data": suppliers})
}
