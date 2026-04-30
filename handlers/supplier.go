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

	middlewares.SimpleLog(c, "update", "supplier", supplier.ID, "编辑供应商: "+supplier.Name)
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
	middlewares.SimpleLog(c, "delete", "supplier", supplier.ID, "删除供应商: "+supplier.Name)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetAllSuppliers(c *gin.Context) {
	var suppliers []models.Supplier
	database.DB.Where("status = 1").Order("name ASC").Find(&suppliers)
	c.JSON(http.StatusOK, gin.H{"data": suppliers})
}

func GetSupplierOrders(c *gin.Context) {
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

	var orders []models.PurchaseOrder
	database.DB.Preload("Supplier").Where("supplier_id = ?", supplier.ID).Order("id DESC").Limit(20).Find(&orders)

	// Get payable info
	type PayableInfo struct {
		TotalAmount float64 `json:"total_amount"`
		PaidAmount  float64 `json:"paid_amount"`
		DueAmount   float64 `json:"due_amount"`
	}
	var payable PayableInfo
	database.DB.Model(&models.AccountPayable{}).
		Select("COALESCE(SUM(total_amount),0) as total_amount, COALESCE(SUM(paid_amount),0) as paid_amount, COALESCE(SUM(due_amount),0) as due_amount").
		Where("supplier_id = ?", supplier.ID).
		Scan(&payable)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"supplier": supplier,
			"orders":   orders,
			"payable":  payable,
			"order_count": len(orders),
		},
	})
}

func ImportSuppliers(c *gin.Context) {
	var req struct {
		Data []struct {
			Code          string `json:"编码"`
			Name          string `json:"名称"`
			ContactPerson string `json:"联系人"`
			Phone         string `json:"电话"`
			Address       string `json:"地址"`
		} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	var success, fail int
	for _, item := range req.Data {
		if item.Code == "" || item.Name == "" {
			fail++
			continue
		}
		var count int64
		database.DB.Model(&models.Supplier{}).Where("code = ?", item.Code).Count(&count)
		if count > 0 {
			fail++
			continue
		}
		supplier := models.Supplier{
			Code:          item.Code,
			Name:          item.Name,
			ContactPerson: item.ContactPerson,
			Phone:         item.Phone,
			Address:       item.Address,
			Status:        1,
		}
		if err := database.DB.Create(&supplier).Error; err != nil {
			fail++
			continue
		}
		success++
	}

	middlewares.SimpleLog(c, "import", "supplier", 0, "导入供应商: 成功"+fmt.Sprintf("%d", success)+"条, 失败"+fmt.Sprintf("%d", fail)+"条")
	c.JSON(http.StatusOK, gin.H{"message": "导入完成", "count": success, "fail": fail})
}
