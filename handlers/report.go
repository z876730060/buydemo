package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/middlewares"
	"github.com/z876730060/buydemo/models"
)

// ========== Reports ==========

// GetPurchaseReport returns purchase order summary report
func GetPurchaseReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type PurchaseSummary struct {
		Date       string  `json:"date"`
		OrderCount int64   `json:"order_count"`
		TotalAmount float64 `json:"total_amount"`
		ReceivedAmount float64 `json:"received_amount"`
	}

	var results []PurchaseSummary
	database.DB.Model(&models.PurchaseOrder{}).
		Select("date(created_at) as date, count(*) as order_count, COALESCE(SUM(total_amount),0) as total_amount, COALESCE(SUM(CASE WHEN status='received' THEN total_amount ELSE 0 END),0) as received_amount").
		Where("created_at >= ? AND created_at < ?", startDate, endDate+" 23:59:59").
		Group("date(created_at)").Order("date(created_at)").Find(&results)

	// Totals
	var totalCount int64
	var totalAmount float64
	database.DB.Model(&models.PurchaseOrder{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate+" 23:59:59").
		Count(&totalCount)
	database.DB.Model(&models.PurchaseOrder{}).
		Select("COALESCE(SUM(total_amount),0)").
		Where("created_at >= ? AND created_at < ?", startDate, endDate+" 23:59:59").
		Scan(&totalAmount)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"details":    results,
			"total_count": totalCount,
			"total_amount": totalAmount,
			"start_date":  startDate,
			"end_date":    endDate,
		},
	})
}

// GetSalesReport returns sales order summary report
func GetSalesReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type SalesSummary struct {
		Date        string  `json:"date"`
		OrderCount  int64   `json:"order_count"`
		TotalAmount float64 `json:"total_amount"`
		DeliveredAmount float64 `json:"delivered_amount"`
	}

	var results []SalesSummary
	database.DB.Model(&models.SalesOrder{}).
		Select("date(created_at) as date, count(*) as order_count, COALESCE(SUM(total_amount),0) as total_amount, COALESCE(SUM(CASE WHEN status='delivered' THEN total_amount ELSE 0 END),0) as delivered_amount").
		Where("created_at >= ? AND created_at < ?", startDate, endDate+" 23:59:59").
		Group("date(created_at)").Order("date(created_at)").Find(&results)

	var totalCount int64
	var totalAmount float64
	database.DB.Model(&models.SalesOrder{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate+" 23:59:59").
		Count(&totalCount)
	database.DB.Model(&models.SalesOrder{}).
		Select("COALESCE(SUM(total_amount),0)").
		Where("created_at >= ? AND created_at < ?", startDate, endDate+" 23:59:59").
		Scan(&totalAmount)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"details":      results,
			"total_count":  totalCount,
			"total_amount": totalAmount,
			"start_date":   startDate,
			"end_date":     endDate,
		},
	})
}

// GetInventoryReport returns inventory status report
func GetInventoryReport(c *gin.Context) {
	type InventoryItem struct {
		ProductCode  string  `json:"product_code"`
		ProductName  string  `json:"product_name"`
		Category     string  `json:"category"`
		Unit         string  `json:"unit"`
		Spec         string  `json:"spec"`
		Quantity     float64 `json:"quantity"`
		Price        float64 `json:"price"`
		StockValue   float64 `json:"stock_value"`
		Status       string  `json:"status"`
	}

	var items []InventoryItem
	database.DB.Model(&models.Inventory{}).
		Select("products.code as product_code, products.name as product_name, products.category, products.unit, products.spec, inventories.quantity, products.price, (inventories.quantity * products.price) as stock_value").
		Joins("JOIN products ON products.id = inventories.product_id").
		Order("products.category ASC, products.code ASC").
		Find(&items)

	// Add status
	for i := range items {
		if items[i].Quantity <= 0 {
			items[i].Status = "out_of_stock"
		} else if items[i].Quantity <= 10 {
			items[i].Status = "low_stock"
		} else {
			items[i].Status = "normal"
		}
	}

	// Summary
	var totalValue float64
	var totalQty float64
	var categoryCount int64
	database.DB.Model(&models.Product{}).Distinct("category").Where("category != ''").Count(&categoryCount)
	for _, item := range items {
		totalValue += item.StockValue
		totalQty += item.Quantity
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"details":       items,
			"total_value":   totalValue,
			"total_quantity": totalQty,
			"category_count": categoryCount,
		},
	})
}

// ========== Data Export (CSV) ==========

func ExportSuppliers(c *gin.Context) {
	var suppliers []models.Supplier
	database.DB.Order("code ASC").Find(&suppliers)

	csv := "\xEF\xBB\xBF"
	csv += "编码,名称,联系人,电话,地址,状态\n"
	for _, s := range suppliers {
		status := "启用"
		if s.Status != 1 {
			status = "禁用"
		}
		csv += s.Code + "," + s.Name + "," + s.ContactPerson + "," + s.Phone + "," + s.Address + "," + status + "\n"
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=suppliers.csv")
	c.String(http.StatusOK, csv)
}

func ExportProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Order("category ASC, code ASC").Find(&products)

	csv := "\xEF\xBB\xBF"
	csv += "编码,名称,分类,单位,规格,参考单价,状态\n"
	for _, p := range products {
		status := "启用"
		if p.Status != 1 {
			status = "禁用"
		}
		csv += p.Code + "," + p.Name + "," + p.Category + "," + p.Unit + "," + p.Spec + "," + fmt.Sprintf("%.2f", p.Price) + "," + status + "\n"
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=products.csv")
	c.String(http.StatusOK, csv)
}

func ExportPurchaseOrders(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -3, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	var orders []models.PurchaseOrder
	database.DB.Preload("Supplier").Preload("Items.Product").
		Where("created_at >= ? AND created_at < ?", startDate, endDate+" 23:59:59").
		Order("id DESC").Find(&orders)

	csv := "\xEF\xBB\xBF"
	csv += "采购单号,供应商,总金额,状态,创建时间,备注\n"
	for _, o := range orders {
		status := o.Status
		statusMap := map[string]string{"draft": "草稿", "approved": "已审核", "received": "已入库", "cancelled": "已取消"}
		if v, ok := statusMap[status]; ok {
			status = v
		}
		csv += o.OrderNo + "," + o.Supplier.Name + "," + fmt.Sprintf("%.2f", o.TotalAmount) + "," + status + "," + o.CreatedAt.Format("2006-01-02 15:04") + "," + o.Remark + "\n"
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=purchase_orders.csv")
	c.String(http.StatusOK, csv)
}

func ExportSalesOrders(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -3, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	var orders []models.SalesOrder
	database.DB.Preload("Customer").Preload("Items.Product").
		Where("created_at >= ? AND created_at < ?", startDate, endDate+" 23:59:59").
		Order("id DESC").Find(&orders)

	csv := "\xEF\xBB\xBF"
	csv += "销售单号,客户,总金额,状态,创建时间,备注\n"
	for _, o := range orders {
		status := o.Status
		statusMap := map[string]string{"draft": "草稿", "approved": "已审核", "delivered": "已出库", "cancelled": "已取消"}
		if v, ok := statusMap[status]; ok {
			status = v
		}
		csv += o.OrderNo + "," + o.Customer.Name + "," + fmt.Sprintf("%.2f", o.TotalAmount) + "," + status + "," + o.CreatedAt.Format("2006-01-02 15:04") + "," + o.Remark + "\n"
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=sales_orders.csv")
	c.String(http.StatusOK, csv)
}

func ExportInventory(c *gin.Context) {
	type InvItem struct {
		ProductCode string  `json:"product_code"`
		ProductName string  `json:"product_name"`
		Category    string  `json:"category"`
		Unit        string  `json:"unit"`
		Spec        string  `json:"spec"`
		Quantity    float64 `json:"quantity"`
		Price       float64 `json:"price"`
		StockValue  float64 `json:"stock_value"`
	}

	var items []InvItem
	database.DB.Model(&models.Inventory{}).
		Select("products.code as product_code, products.name as product_name, products.category, products.unit, products.spec, inventories.quantity, products.price, (inventories.quantity * products.price) as stock_value").
		Joins("JOIN products ON products.id = inventories.product_id").
		Order("products.category ASC, products.code ASC").
		Find(&items)

	csv := "\xEF\xBB\xBF"
	csv += "编码,名称,分类,单位,规格,库存数量,单价,库存价值\n"
	for _, item := range items {
		csv += item.ProductCode + "," + item.ProductName + "," + item.Category + "," + item.Unit + "," + item.Spec + "," +
			fmt.Sprintf("%.2f", item.Quantity) + "," + fmt.Sprintf("%.2f", item.Price) + "," + fmt.Sprintf("%.2f", item.StockValue) + "\n"
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=inventory.csv")
	c.String(http.StatusOK, csv)
}

// Log export operations and serve
func LoggedExport(handler gin.HandlerFunc, target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		middlewares.SimpleLog(c, "export", target, 0, "导出"+target+"数据")
		handler(c)
	}
}
