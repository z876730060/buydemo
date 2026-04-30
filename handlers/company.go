package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/middlewares"
	"github.com/z876730060/buydemo/models"
)

func GetCompanies(c *gin.Context) {
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
	query := database.DB.Model(&models.Company{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR contact LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var companies []models.Company
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&companies)

	c.JSON(http.StatusOK, gin.H{
		"data":  companies,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var company models.Company
	if err := database.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公司不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": company})
}

func CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	if company.Code == "" || company.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "编码和名称不能为空"})
		return
	}

	if err := database.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败，编码可能已存在"})
		return
	}

	middlewares.SimpleLog(c, "create", "company", company.ID, "新增公司: "+company.Name)
	c.JSON(http.StatusOK, gin.H{"data": company, "message": "创建成功"})
}

func UpdateCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var company models.Company
	if err := database.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公司不存在"})
		return
	}

	var input models.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}

	database.DB.Model(&company).Updates(map[string]interface{}{
		"code":    input.Code,
		"name":    input.Name,
		"contact": input.Contact,
		"phone":   input.Phone,
		"email":   input.Email,
		"address": input.Address,
		"tax_no":  input.TaxNo,
		"status":  input.Status,
		"remark":  input.Remark,
	})

	middlewares.SimpleLog(c, "update", "company", company.ID, "编辑公司: "+company.Name)
	c.JSON(http.StatusOK, gin.H{"data": company, "message": "更新成功"})
}

func DeleteCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var company models.Company
	if err := database.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公司不存在"})
		return
	}

	var warehouseCount int64
	database.DB.Model(&models.Warehouse{}).Where("company_id = ?", company.ID).Count(&warehouseCount)
	if warehouseCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该公司下还有仓库，无法删除"})
		return
	}

	database.DB.Delete(&company)
	middlewares.SimpleLog(c, "delete", "company", company.ID, "删除公司: "+company.Name)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetAllCompanies(c *gin.Context) {
	var companies []models.Company
	database.DB.Where("status = 1").Order("name ASC").Find(&companies)
	c.JSON(http.StatusOK, gin.H{"data": companies})
}
