package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/middlewares"
	"github.com/z876730060/buydemo/models"
	"github.com/google/uuid"
)

func generateContractNo() string {
	now := time.Now()
	var count int64
	database.DB.Model(&models.Contract{}).Where("created_at >= ? AND created_at < ?",
		now.Format("2006-01-02"), now.AddDate(0, 0, 1).Format("2006-01-02")).Count(&count)
	return fmt.Sprintf("CT%s%04d", now.Format("20060102"), count+1)
}

func GetContracts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	ctype := c.Query("type")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query := database.DB.Model(&models.Contract{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if ctype != "" {
		query = query.Where("type = ?", ctype)
	}
	if keyword != "" {
		query = query.Where("contract_no LIKE ? OR title LIKE ? OR party_b_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var contracts []models.Contract
	query.Preload("Creator").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&contracts)

	c.JSON(http.StatusOK, gin.H{
		"data":  contracts,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func GetContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var contract models.Contract
	if err := database.DB.Preload("Creator").First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "合同不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contract})
}

type CreateContractRequest struct {
	Title       string  `json:"title" binding:"required"`
	Type        string  `json:"type" binding:"required"` // purchase / sales
	PartyBType  string  `json:"party_b_type" binding:"required"` // supplier / customer
	PartyBID    uint    `json:"party_b_id" binding:"required"`
	Content     string  `json:"content"`
	TotalAmount float64 `json:"total_amount"`
	Remark      string  `json:"remark"`
}

func CreateContract(c *gin.Context) {
	var req CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	var partyBName string
	var partyBContact string
	if req.PartyBType == "supplier" {
		var s models.Supplier
		if err := database.DB.First(&s, req.PartyBID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "供应商不存在"})
			return
		}
		partyBName = s.Name
		partyBContact = s.ContactPerson + " " + s.Phone
	} else if req.PartyBType == "customer" {
		var cus models.Customer
		if err := database.DB.First(&cus, req.PartyBID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "客户不存在"})
			return
		}
		partyBName = cus.Name
		partyBContact = cus.ContactPerson + " " + cus.Phone
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "乙方类型必须是 supplier 或 customer"})
		return
	}

	contract := models.Contract{
		ContractNo:    generateContractNo(),
		Title:         req.Title,
		Type:          req.Type,
		PartyBType:    req.PartyBType,
		PartyBID:      req.PartyBID,
		PartyBName:    partyBName,
		PartyBContact: partyBContact,
		Content:       req.Content,
		TotalAmount:   req.TotalAmount,
		Status:        "draft",
		SignToken:     uuid.New().String(),
		CreatedBy:     c.GetUint("user_id"),
		Remark:        req.Remark,
	}

	if err := database.DB.Create(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建合同失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contract, "message": "创建成功"})
}

func UpdateContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "合同不存在"})
		return
	}

	if contract.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能编辑草稿状态的合同"})
		return
	}

	var req CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	partyBName := contract.PartyBName
	partyBContact := contract.PartyBContact
	if req.PartyBID != contract.PartyBID || req.PartyBType != contract.PartyBType {
		if req.PartyBType == "supplier" {
			var s models.Supplier
			if err := database.DB.First(&s, req.PartyBID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "供应商不存在"})
				return
			}
			partyBName = s.Name
			partyBContact = s.ContactPerson + " " + s.Phone
		} else {
			var cus models.Customer
			if err := database.DB.First(&cus, req.PartyBID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "客户不存在"})
				return
			}
			partyBName = cus.Name
			partyBContact = cus.ContactPerson + " " + cus.Phone
		}
	}

	database.DB.Model(&contract).Updates(map[string]interface{}{
		"title":          req.Title,
		"type":           req.Type,
		"party_b_type":   req.PartyBType,
		"party_b_id":     req.PartyBID,
		"party_b_name":   partyBName,
		"party_b_contact": partyBContact,
		"content":        req.Content,
		"total_amount":   req.TotalAmount,
		"remark":         req.Remark,
	})

	middlewares.SimpleLog(c, "update", "contract", contract.ID, "编辑合同: "+contract.ContractNo)
	c.JSON(http.StatusOK, gin.H{"data": contract, "message": "更新成功"})
}

func DeleteContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "合同不存在"})
		return
	}

	if contract.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能删除草稿状态的合同"})
		return
	}

	database.DB.Delete(&contract)
	middlewares.SimpleLog(c, "delete", "contract", contract.ID, "删除合同: "+contract.ContractNo)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func SignContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		Party     string `json:"party" binding:"required"`
		Signature string `json:"signature" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "合同不存在"})
		return
	}

	if contract.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已取消的合同无法签字"})
		return
	}

	now := time.Now()

	if req.Party == "a" {
		if contract.Status != "draft" && contract.Status != "sent" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "当前状态不允许甲方签字"})
			return
		}
		database.DB.Model(&contract).Updates(map[string]interface{}{
			"signed_by_a": req.Signature,
			"signed_at_a": &now,
			"status":      "sent",
		})
		middlewares.SimpleLog(c, "sign", "contract", contract.ID, "甲方签字: "+contract.ContractNo)
	} else if req.Party == "b" {
		if contract.Status != "sent" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先让甲方签字后再由乙方签字"})
			return
		}
		newStatus := "completed"
		if contract.SignedByA == "" {
			newStatus = "signed_by_b"
		}
		database.DB.Model(&contract).Updates(map[string]interface{}{
			"signed_by_b": req.Signature,
			"signed_at_b": &now,
			"status":      newStatus,
		})
		middlewares.SimpleLog(c, "sign", "contract", contract.ID, "乙方签字: "+contract.ContractNo)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "party 必须是 a 或 b"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contract, "message": "签字成功"})
}

func CancelContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "合同不存在"})
		return
	}

	if contract.Status == "completed" || contract.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已完成或已取消的合同无法取消"})
		return
	}

	database.DB.Model(&contract).Update("status", "cancelled")
	middlewares.SimpleLog(c, "cancel", "contract", contract.ID, "取消合同: "+contract.ContractNo)
	c.JSON(http.StatusOK, gin.H{"data": contract, "message": "已取消"})
}

func GetContractParties(c *gin.Context) {
	ptype := c.Query("type")
	if ptype == "supplier" {
		var suppliers []models.Supplier
		database.DB.Where("status = 1").Order("name ASC").Find(&suppliers)
		var result []gin.H
		for _, s := range suppliers {
			result = append(result, gin.H{"id": s.ID, "name": s.Name, "contact": s.ContactPerson + " " + s.Phone})
		}
		c.JSON(http.StatusOK, gin.H{"data": result})
	} else if ptype == "customer" {
		var customers []models.Customer
		database.DB.Where("status = 1").Order("name ASC").Find(&customers)
		var result []gin.H
		for _, cus := range customers {
			result = append(result, gin.H{"id": cus.ID, "name": cus.Name, "contact": cus.ContactPerson + " " + cus.Phone})
		}
		c.JSON(http.StatusOK, gin.H{"data": result})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type 必须是 supplier 或 customer"})
	}
}

// GenerateSignToken (re)generates a sign token for sharing
func GenerateSignToken(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "合同不存在"})
		return
	}

	if contract.Status == "completed" || contract.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已完成或已取消的合同无法生成签署链接"})
		return
	}

	token := uuid.New().String()
	database.DB.Model(&contract).Update("sign_token", token)

	signURL := fmt.Sprintf("%s/contract-sign?token=%s", getBaseURL(c), token)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"sign_token": token,
			"sign_url":   signURL,
		},
		"message": "签署链接已生成",
	})
}

func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// ========== Public (unauthorized) endpoints ==========

// GetPublicContract returns contract info by sign token (no auth)
func GetPublicContract(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少签署令牌"})
		return
	}

	var contract models.Contract
	if err := database.DB.Where("sign_token = ?", token).First(&contract).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "合同不存在或链接已失效"})
		return
	}

	if contract.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该合同已取消"})
		return
	}

	// Return only public-safe fields
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":            contract.ID,
			"contract_no":   contract.ContractNo,
			"title":         contract.Title,
			"type":          contract.Type,
			"party_b_name":  contract.PartyBName,
			"party_b_contact": contract.PartyBContact,
			"content":       contract.Content,
			"total_amount":  contract.TotalAmount,
			"status":        contract.Status,
			"signed_by_a":   contract.SignedByA,
			"signed_by_b":   contract.SignedByB,
			"signed_at_a":   contract.SignedAtA,
			"signed_at_b":   contract.SignedAtB,
		},
	})
}

// SignPublicContract signs a contract via public token (no auth, party b only)
func SignPublicContract(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少签署令牌"})
		return
	}

	var req struct {
		Signature string `json:"signature" binding:"required"`
		Name      string `json:"name"` // signer name for record
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	var contract models.Contract
	if err := database.DB.Where("sign_token = ?", token).First(&contract).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "合同不存在或链接已失效"})
		return
	}

	if contract.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该合同已取消"})
		return
	}

	if contract.Status != "sent" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "合同尚未由甲方签署或已完成签署"})
		return
	}

	if contract.SignedByB != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "乙方已完成签署"})
		return
	}

	now := time.Now()
	newStatus := "completed"

	database.DB.Model(&contract).Updates(map[string]interface{}{
		"signed_by_b": req.Signature,
		"signed_at_b": &now,
		"status":      newStatus,
	})

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"contract_no":  contract.ContractNo,
		"title":        contract.Title,
		"status":       newStatus,
		"signed_at_b":  &now,
	}, "message": "签署成功，感谢您的签字！"})
}
