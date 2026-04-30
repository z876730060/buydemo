package router

import (
	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/handlers"
	"github.com/z876730060/buydemo/middlewares"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	api := r.Group("/api")
	{
		// Public routes
		api.POST("/auth/login", handlers.Login)

		// Protected routes
		auth := api.Group("")
		auth.Use(middlewares.AuthRequired())
		{
			auth.GET("/auth/me", handlers.Me)
			auth.POST("/auth/change-password", handlers.ChangePassword)

			// Suppliers
			auth.GET("/suppliers", handlers.GetSuppliers)
			auth.GET("/suppliers/all", handlers.GetAllSuppliers)
			auth.GET("/suppliers/:id", handlers.GetSupplier)
			auth.GET("/suppliers/:id/orders", handlers.GetSupplierOrders)
			auth.POST("/suppliers", middlewares.LogOperation("create", "supplier", 0, "新增供应商"), handlers.CreateSupplier)
			auth.PUT("/suppliers/:id", handlers.UpdateSupplier)
			auth.DELETE("/suppliers/:id", handlers.DeleteSupplier)
			auth.POST("/suppliers/import", middlewares.LogOperation("import", "supplier", 0, "导入供应商"), handlers.ImportSuppliers)

			// Products
			auth.GET("/products", handlers.GetProducts)
			auth.GET("/products/all", handlers.GetAllProducts)
			auth.GET("/products/categories", handlers.GetCategories)
			auth.GET("/products/:id", handlers.GetProduct)
			auth.GET("/products/:id/detail", handlers.GetProductDetail)
			auth.POST("/products", middlewares.LogOperation("create", "product", 0, "新增商品"), handlers.CreateProduct)
			auth.PUT("/products/:id", handlers.UpdateProduct)
			auth.DELETE("/products/:id", handlers.DeleteProduct)
			auth.POST("/products/import", middlewares.LogOperation("import", "product", 0, "导入商品"), handlers.ImportProducts)

			// Purchase Orders
			auth.GET("/purchase-orders", handlers.GetPurchaseOrders)
			auth.GET("/purchase-orders/:id", handlers.GetPurchaseOrder)
			auth.POST("/purchase-orders", middlewares.LogOperation("create", "purchase_order", 0, "新建采购单"), handlers.CreatePurchaseOrder)
			auth.PUT("/purchase-orders/:id", handlers.UpdatePurchaseOrder)
			auth.POST("/purchase-orders/:id/approve", middlewares.LogOperation("approve", "purchase_order", 0, "审核采购单"), handlers.ApproveOrder)
			auth.POST("/purchase-orders/:id/receive", middlewares.LogOperation("receive", "purchase_order", 0, "采购入库"), handlers.ReceiveOrder)
			auth.POST("/purchase-orders/:id/cancel", middlewares.LogOperation("cancel", "purchase_order", 0, "取消采购单"), handlers.CancelOrder)

			// Customers
			auth.GET("/customers", handlers.GetCustomers)
			auth.GET("/customers/all", handlers.GetAllCustomers)
			auth.GET("/customers/:id", handlers.GetCustomer)
			auth.GET("/customers/:id/orders", handlers.GetCustomerOrders)
			auth.POST("/customers", middlewares.LogOperation("create", "customer", 0, "新增客户"), handlers.CreateCustomer)
			auth.PUT("/customers/:id", handlers.UpdateCustomer)
			auth.DELETE("/customers/:id", handlers.DeleteCustomer)

			// Sales Orders
			auth.GET("/sales-orders", handlers.GetSalesOrders)
			auth.GET("/sales-orders/:id", handlers.GetSalesOrder)
			auth.POST("/sales-orders", middlewares.LogOperation("create", "sales_order", 0, "新建销售单"), handlers.CreateSalesOrder)
			auth.PUT("/sales-orders/:id", handlers.UpdateSalesOrder)
			auth.POST("/sales-orders/:id/approve", middlewares.LogOperation("approve", "sales_order", 0, "审核销售单"), handlers.ApproveSalesOrder)
			auth.POST("/sales-orders/:id/deliver", middlewares.LogOperation("deliver", "sales_order", 0, "销售出库"), handlers.DeliverSalesOrder)
			auth.POST("/sales-orders/:id/cancel", middlewares.LogOperation("cancel", "sales_order", 0, "取消销售单"), handlers.CancelSalesOrder)

			// Inventory
			auth.GET("/inventories", handlers.GetInventories)
			auth.GET("/inventories/logs", handlers.GetInventoryLogs)
			auth.GET("/inventories/low-stock", handlers.GetLowStock)
			auth.POST("/inventories/adjust", middlewares.LogOperation("adjust", "inventory", 0, "库存调整"), handlers.AdjustInventory)

			// Users (admin only)
			auth.GET("/users", handlers.GetUsers)
			auth.POST("/users", middlewares.LogOperation("create", "user", 0, "新增用户"), handlers.CreateUser)
			auth.PUT("/users/:id", handlers.UpdateUser)
			auth.DELETE("/users/:id", handlers.DeleteUser)
			auth.POST("/users/:id/reset-password", middlewares.LogOperation("reset_password", "user", 0, "重置密码"), handlers.ResetPassword)

			// Dashboard
			auth.GET("/dashboard/stats", handlers.GetDashboardStats)

			// Finance - Accounts Payable
			auth.GET("/finance/payable", handlers.GetAccountsPayable)
			auth.POST("/finance/payable/:id/pay", middlewares.LogOperation("pay", "account_payable", 0, "付款"), handlers.PayAccountPayable)

			// Finance - Accounts Receivable
			auth.GET("/finance/receivable", handlers.GetAccountsReceivable)
			auth.POST("/finance/receivable/:id/receive", middlewares.LogOperation("receive", "account_receivable", 0, "收款"), handlers.ReceiveAccountReceivable)

			// Finance - Expenses
			auth.GET("/finance/expenses", handlers.GetExpenses)
			auth.POST("/finance/expenses", middlewares.LogOperation("create", "expense", 0, "新增费用"), handlers.CreateExpense)
			auth.PUT("/finance/expenses/:id", handlers.UpdateExpense)
			auth.DELETE("/finance/expenses/:id", handlers.DeleteExpense)

			// Finance - Payment Records
			auth.GET("/finance/payments", handlers.GetPaymentRecords)

			// Finance - Summary
			auth.GET("/finance/summary", handlers.GetFinancialSummary)

			// Finance - Options
			auth.GET("/finance/payment-methods", handlers.GetPaymentMethods)
			auth.GET("/finance/expense-categories", handlers.GetExpenseCategories)

			// Operation Logs
			auth.GET("/operation-logs", handlers.GetOperationLogs)
			auth.GET("/operation-logs/export", handlers.LoggedExport(handlers.ExportOperationLogs, "操作日志"))

			// Reports
			auth.GET("/reports/purchase", handlers.GetPurchaseReport)
			auth.GET("/reports/sales", handlers.GetSalesReport)
			auth.GET("/reports/inventory", handlers.GetInventoryReport)

			// Data Export
			auth.GET("/export/suppliers", handlers.LoggedExport(handlers.ExportSuppliers, "供应商"))
			auth.GET("/export/products", handlers.LoggedExport(handlers.ExportProducts, "商品"))
			auth.GET("/export/purchase-orders", handlers.LoggedExport(handlers.ExportPurchaseOrders, "采购单"))
			auth.GET("/export/sales-orders", handlers.LoggedExport(handlers.ExportSalesOrders, "销售单"))
			auth.GET("/export/inventory", handlers.LoggedExport(handlers.ExportInventory, "库存"))
		}
	}

	return r
}
