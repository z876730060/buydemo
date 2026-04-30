package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/middlewares"
	"github.com/z876730060/buydemo/models"
)

// ========== Customers ==========

func GetCustomers(c *gin.Context) {
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
	query := database.DB.Model(&models.Customer{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR contact_person LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var customers []models.Customer
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&customers)

	c.JSON(http.StatusOK, gin.H{
		"data":  customers,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var customer models.Customer
	if err := database.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	if customer.Code == "" || customer.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "编码和名称不能为空"})
		return
	}

	if err := database.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败，编码可能已存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer, "message": "创建成功"})
}

func UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var customer models.Customer
	if err := database.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户不存在"})
		return
	}

	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}

	database.DB.Model(&customer).Updates(map[string]interface{}{
		"code":           input.Code,
		"name":           input.Name,
		"contact_person": input.ContactPerson,
		"phone":          input.Phone,
		"address":        input.Address,
		"status":         input.Status,
	})

	middlewares.SimpleLog(c, "update", "customer", customer.ID, "编辑客户: "+customer.Name)
	c.JSON(http.StatusOK, gin.H{"data": customer, "message": "更新成功"})
}

func DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var customer models.Customer
	if err := database.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户不存在"})
		return
	}

	database.DB.Delete(&customer)
	middlewares.SimpleLog(c, "delete", "customer", customer.ID, "删除客户: "+customer.Name)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetAllCustomers(c *gin.Context) {
	var customers []models.Customer
	database.DB.Where("status = 1").Order("name ASC").Find(&customers)
	c.JSON(http.StatusOK, gin.H{"data": customers})
}
