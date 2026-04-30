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
			auth.POST("/suppliers", handlers.CreateSupplier)
			auth.PUT("/suppliers/:id", handlers.UpdateSupplier)
			auth.DELETE("/suppliers/:id", handlers.DeleteSupplier)

			// Products
			auth.GET("/products", handlers.GetProducts)
			auth.GET("/products/all", handlers.GetAllProducts)
			auth.GET("/products/categories", handlers.GetCategories)
			auth.GET("/products/:id", handlers.GetProduct)
			auth.POST("/products", handlers.CreateProduct)
			auth.PUT("/products/:id", handlers.UpdateProduct)
			auth.DELETE("/products/:id", handlers.DeleteProduct)

			// Purchase Orders
			auth.GET("/purchase-orders", handlers.GetPurchaseOrders)
			auth.GET("/purchase-orders/:id", handlers.GetPurchaseOrder)
			auth.POST("/purchase-orders", handlers.CreatePurchaseOrder)
			auth.PUT("/purchase-orders/:id", handlers.UpdatePurchaseOrder)
			auth.POST("/purchase-orders/:id/approve", handlers.ApproveOrder)
			auth.POST("/purchase-orders/:id/receive", handlers.ReceiveOrder)
			auth.POST("/purchase-orders/:id/cancel", handlers.CancelOrder)

			// Inventory
			auth.GET("/inventories", handlers.GetInventories)
			auth.GET("/inventories/logs", handlers.GetInventoryLogs)
			auth.GET("/inventories/low-stock", handlers.GetLowStock)

			// Users (admin only)
			auth.GET("/users", handlers.GetUsers)
			auth.POST("/users", handlers.CreateUser)
			auth.PUT("/users/:id", handlers.UpdateUser)
			auth.DELETE("/users/:id", handlers.DeleteUser)
			auth.POST("/users/:id/reset-password", handlers.ResetPassword)

			// Dashboard
			auth.GET("/dashboard/stats", handlers.GetDashboardStats)
		}
	}

	return r
}
