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

// ========== Accounts Payable ==========

func GetAccountsPayable(c *gin.Context) {
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
	query := database.DB.Model(&models.AccountPayable{}).Preload("Supplier")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Joins("JOIN suppliers ON suppliers.id = account_payables.supplier_id").
			Where("account_payables.order_no LIKE ? OR suppliers.name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var list []models.AccountPayable
	query.Preload("Supplier").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"data":  list,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func PayAccountPayable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		PaymentMethod string  `json:"payment_method"`
		ReferenceNo   string  `json:"reference_no"`
		Remark        string  `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	var payable models.AccountPayable
	if err := database.DB.First(&payable, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "应付账款记录不存在"})
		return
	}

	if payable.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该应付账款已结清"})
		return
	}

	if req.Amount > payable.DueAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("付款金额（%.2f）不能超过应付金额（%.2f）", req.Amount, payable.DueAmount)})
		return
	}

	userID := c.GetUint("user_id")
	now := time.Now()

	tx := database.DB.Begin()

	// Update payable
	newPaid := payable.PaidAmount + req.Amount
	newDue := payable.TotalAmount - newPaid
	newStatus := "pending"
	if newDue <= 0 {
		newStatus = "paid"
	} else if newPaid > 0 {
		newStatus = "partial"
	}
	updates := map[string]interface{}{
		"paid_amount": newPaid,
		"due_amount":  newDue,
		"status":      newStatus,
	}
	if newStatus == "paid" {
		updates["paid_at"] = &now
	}
	tx.Model(&payable).Updates(updates)

	// Create payment record
	payment := models.PaymentRecord{
		PayableID:     &payable.ID,
		Type:          "pay",
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		ReferenceNo:   req.ReferenceNo,
		Remark:        req.Remark,
		OperatedBy:    userID,
		OperatedAt:    now,
	}
	tx.Create(&payment)

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"data": payable, "message": "付款成功"})
}

// ========== Accounts Receivable ==========

func GetAccountsReceivable(c *gin.Context) {
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
	query := database.DB.Model(&models.AccountReceivable{}).Preload("Customer")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Joins("JOIN customers ON customers.id = account_receivables.customer_id").
			Where("account_receivables.order_no LIKE ? OR customers.name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var list []models.AccountReceivable
	query.Preload("Customer").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"data":  list,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func ReceiveAccountReceivable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		PaymentMethod string  `json:"payment_method"`
		ReferenceNo   string  `json:"reference_no"`
		Remark        string  `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	var receivable models.AccountReceivable
	if err := database.DB.First(&receivable, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "应收账款记录不存在"})
		return
	}

	if receivable.Status == "received" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该应收账款已结清"})
		return
	}

	if req.Amount > receivable.DueAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("收款金额（%.2f）不能超过应收金额（%.2f）", req.Amount, receivable.DueAmount)})
		return
	}

	userID := c.GetUint("user_id")
	now := time.Now()

	tx := database.DB.Begin()

	newReceived := receivable.ReceivedAmount + req.Amount
	newDue := receivable.TotalAmount - newReceived
	newStatus := "pending"
	if newDue <= 0 {
		newStatus = "received"
	} else if newReceived > 0 {
		newStatus = "partial"
	}
	updates := map[string]interface{}{
		"received_amount": newReceived,
		"due_amount":      newDue,
		"status":          newStatus,
	}
	if newStatus == "received" {
		updates["received_at"] = &now
	}
	tx.Model(&receivable).Updates(updates)

	// Create payment record
	payment := models.PaymentRecord{
		ReceivableID:  &receivable.ID,
		Type:          "receive",
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		ReferenceNo:   req.ReferenceNo,
		Remark:        req.Remark,
		OperatedBy:    userID,
		OperatedAt:    now,
	}
	tx.Create(&payment)

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"data": receivable, "message": "收款成功"})
}

// ========== Expenses ==========

func GetExpenses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query := database.DB.Model(&models.Expense{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if startDate != "" {
		query = query.Where("occurred_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("occurred_at <= ?", endDate+" 23:59:59")
	}
	query.Count(&total)

	var expenses []models.Expense
	query.Preload("Creator").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&expenses)

	c.JSON(http.StatusOK, gin.H{
		"data":  expenses,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func CreateExpense(c *gin.Context) {
	var req struct {
		Category    string  `json:"category" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
		OccurredAt  string  `json:"occurred_at"` // optional: 2006-01-02
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误: " + err.Error()})
		return
	}

	occurredAt := time.Now()
	if req.OccurredAt != "" {
		if t, err := time.Parse("2006-01-02", req.OccurredAt); err == nil {
			occurredAt = t
		}
	}

	expense := models.Expense{
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
		OccurredAt:  occurredAt,
		CreatedBy:   c.GetUint("user_id"),
	}

	if err := database.DB.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建费用记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expense, "message": "创建成功"})
}

func UpdateExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var expense models.Expense
	if err := database.DB.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "费用记录不存在"})
		return
	}

	var req struct {
		Category    string  `json:"category"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		OccurredAt  string  `json:"occurred_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Amount > 0 {
		updates["amount"] = req.Amount
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.OccurredAt != "" {
		if t, err := time.Parse("2006-01-02", req.OccurredAt); err == nil {
			updates["occurred_at"] = t
		}
	}
	if len(updates) > 0 {
		database.DB.Model(&expense).Updates(updates)
	}

	c.JSON(http.StatusOK, gin.H{"data": expense, "message": "更新成功"})
}

func DeleteExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var expense models.Expense
	if err := database.DB.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "费用记录不存在"})
		return
	}

	database.DB.Delete(&expense)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ========== Payment Records ==========

func GetPaymentRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	ptype := c.Query("type") // pay / receive

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query := database.DB.Model(&models.PaymentRecord{}).Preload("Operator")
	if ptype != "" {
		query = query.Where("type = ?", ptype)
	}
	query.Count(&total)

	var records []models.PaymentRecord
	query.Preload("Operator").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"data":  records,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// ========== Payment Methods ==========

func GetPaymentMethods(c *gin.Context) {
	methods := []string{"现金", "银行转账", "微信支付", "支付宝", "支票", "其他"}
	c.JSON(http.StatusOK, gin.H{"data": methods})
}

func GetExpenseCategories(c *gin.Context) {
	categories := []string{"办公费", "运输费", "税费", "水电费", "租金", "工资", "采购成本", "其他"}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// ========== Financial Summary ==========

func GetFinancialSummary(c *gin.Context) {
	type SummaryResult struct {
		TotalAmount   float64 `json:"total_amount"`
		PendingAmount float64 `json:"pending_amount"`
		PaidAmount    float64 `json:"paid_amount"`
		Count         int64   `json:"count"`
	}

	// Accounts Payable Summary
	var apSummary SummaryResult
	database.DB.Model(&models.AccountPayable{}).
		Select("COALESCE(SUM(total_amount),0) as total_amount, COALESCE(SUM(CASE WHEN status != 'paid' THEN due_amount ELSE 0 END),0) as pending_amount, COALESCE(SUM(paid_amount),0) as paid_amount, COUNT(*) as count").
		Scan(&apSummary)

	// Accounts Receivable Summary
	var arSummary SummaryResult
	database.DB.Model(&models.AccountReceivable{}).
		Select("COALESCE(SUM(total_amount),0) as total_amount, COALESCE(SUM(CASE WHEN status != 'received' THEN due_amount ELSE 0 END),0) as pending_amount, COALESCE(SUM(received_amount),0) as paid_amount, COUNT(*) as count").
		Scan(&arSummary)

	// Expense Summary
	var expenseTotal float64
	database.DB.Model(&models.Expense{}).Select("COALESCE(SUM(amount),0)").Scan(&expenseTotal)

	// Monthly income/expense for current year
	year := time.Now().Format("2006")
	type MonthlyData struct {
		Month  string  `json:"month"`
		Income float64 `json:"income"`
		Expense float64 `json:"expense"`
	}
	var monthlyData []MonthlyData

	// Income from payment records (receive)
	database.DB.Model(&models.PaymentRecord{}).
		Select("strftime('%m', operated_at) as month, COALESCE(SUM(amount),0) as income, 0 as expense").
		Where("type = 'receive' AND strftime('%Y', operated_at) = ?", year).
		Group("month").Order("month").Find(&monthlyData)

	// Expense from payment records (pay) + expense table
	type PayExpense struct {
		Month  string  `json:"month"`
		Amount float64 `json:"amount"`
	}
	var payExpenses []PayExpense
	database.DB.Model(&models.PaymentRecord{}).
		Select("strftime('%m', operated_at) as month, COALESCE(SUM(amount),0) as amount").
		Where("type = 'pay' AND strftime('%Y', operated_at) = ?", year).
		Group("month").Order("month").Find(&payExpenses)

	var expenseByMonth []PayExpense
	database.DB.Model(&models.Expense{}).
		Select("strftime('%m', occurred_at) as month, COALESCE(SUM(amount),0) as amount").
		Where("strftime('%Y', occurred_at) = ?", year).
		Group("month").Order("month").Find(&expenseByMonth)

	// Merge
	monthMap := make(map[string]*MonthlyData)
	for i := 1; i <= 12; i++ {
		m := fmt.Sprintf("%02d", i)
		monthMap[m] = &MonthlyData{Month: m, Income: 0, Expense: 0}
	}
	for _, d := range monthlyData {
		if v, ok := monthMap[d.Month]; ok {
			v.Income = d.Income
		}
	}
	for _, d := range payExpenses {
		if v, ok := monthMap[d.Month]; ok {
			v.Expense += d.Amount
		}
	}
	for _, d := range expenseByMonth {
		if v, ok := monthMap[d.Month]; ok {
			v.Expense += d.Amount
		}
	}
	var result []MonthlyData
	for i := 1; i <= 12; i++ {
		m := fmt.Sprintf("%02d", i)
		result = append(result, *monthMap[m])
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"accounts_payable":   apSummary,
			"accounts_receivable": arSummary,
			"total_expense":      expenseTotal,
			"monthly_data":       result,
		},
	})
}

// ========== Auto-generate from orders ==========

// GenerateAccountPayable creates account payable record when purchase order is received
func GenerateAccountPayable(order models.PurchaseOrder) {
	// Check if already exists
	var count int64
	database.DB.Model(&models.AccountPayable{}).Where("purchase_order_id = ?", order.ID).Count(&count)
	if count > 0 {
		return
	}

	payable := models.AccountPayable{
		OrderNo:         order.OrderNo,
		PurchaseOrderID: order.ID,
		SupplierID:      order.SupplierID,
		TotalAmount:     order.TotalAmount,
		PaidAmount:      0,
		DueAmount:       order.TotalAmount,
		Status:          "pending",
	}
	database.DB.Create(&payable)
}

// GenerateAccountReceivable creates account receivable record when sales order is delivered
func GenerateAccountReceivable(order models.SalesOrder) {
	var count int64
	database.DB.Model(&models.AccountReceivable{}).Where("sales_order_id = ?", order.ID).Count(&count)
	if count > 0 {
		return
	}

	receivable := models.AccountReceivable{
		OrderNo:        order.OrderNo,
		SalesOrderID:   order.ID,
		CustomerID:     order.CustomerID,
		TotalAmount:    order.TotalAmount,
		ReceivedAmount: 0,
		DueAmount:      order.TotalAmount,
		Status:         "pending",
	}
	database.DB.Create(&receivable)
}
