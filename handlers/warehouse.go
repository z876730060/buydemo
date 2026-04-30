package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/middlewares"
	"github.com/z876730060/buydemo/models"
)

func GetWarehouses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	companyID, _ := strconv.Atoi(c.Query("company_id"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query := database.DB.Model(&models.Warehouse{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if companyID > 0 {
		query = query.Where("company_id = ?", companyID)
	}
	query.Count(&total)

	var warehouses []models.Warehouse
	query.Preload("Company").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&warehouses)

	c.JSON(http.StatusOK, gin.H{
		"data":  warehouses,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetWarehouse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var warehouse models.Warehouse
	if err := database.DB.Preload("Company").First(&warehouse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": warehouse})
}

func CreateWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	if warehouse.Code == "" || warehouse.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "编码和名称不能为空"})
		return
	}

	if warehouse.IsDefault {
		database.DB.Model(&models.Warehouse{}).Where("company_id = ?", warehouse.CompanyID).Update("is_default", false)
	}

	if err := database.DB.Create(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败，编码可能已存在"})
		return
	}

	middlewares.SimpleLog(c, "create", "warehouse", warehouse.ID, "新增仓库: "+warehouse.Name)
	c.JSON(http.StatusOK, gin.H{"data": warehouse, "message": "创建成功"})
}

func UpdateWarehouse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var warehouse models.Warehouse
	if err := database.DB.First(&warehouse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	var input models.Warehouse
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}

	if input.IsDefault {
		database.DB.Model(&models.Warehouse{}).Where("company_id = ?", warehouse.CompanyID).Update("is_default", false)
	}

	database.DB.Model(&warehouse).Updates(map[string]interface{}{
		"code":       input.Code,
		"name":       input.Name,
		"company_id": input.CompanyID,
		"address":    input.Address,
		"contact":    input.Contact,
		"phone":      input.Phone,
		"is_default": input.IsDefault,
		"status":     input.Status,
		"remark":     input.Remark,
	})

	middlewares.SimpleLog(c, "update", "warehouse", warehouse.ID, "编辑仓库: "+warehouse.Name)
	c.JSON(http.StatusOK, gin.H{"data": warehouse, "message": "更新成功"})
}

func DeleteWarehouse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var warehouse models.Warehouse
	if err := database.DB.First(&warehouse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	var inventoryCount int64
	database.DB.Model(&models.Inventory{}).Where("warehouse_id = ?", warehouse.ID).Count(&inventoryCount)
	if inventoryCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该仓库下还有库存，无法删除"})
		return
	}

	database.DB.Delete(&warehouse)
	middlewares.SimpleLog(c, "delete", "warehouse", warehouse.ID, "删除仓库: "+warehouse.Name)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetAllWarehouses(c *gin.Context) {
	companyID, _ := strconv.Atoi(c.Query("company_id"))
	var warehouses []models.Warehouse
	query := database.DB.Where("status = 1")
	if companyID > 0 {
		query = query.Where("company_id = ?", companyID)
	}
	query.Preload("Company").Order("name ASC").Find(&warehouses)
	c.JSON(http.StatusOK, gin.H{"data": warehouses})
}
