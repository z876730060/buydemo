package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/z876730060/buydemo/config"
	"github.com/z876730060/buydemo/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) {
	// ensure data directory exists
	dataDir := filepath.Dir(cfg.DBPath)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Auto migrate
	err = DB.AutoMigrate(
		&models.User{},
		&models.Supplier{},
		&models.Product{},
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.Inventory{},
		&models.InventoryLog{},
		&models.Customer{},
		&models.SalesOrder{},
		&models.SalesOrderItem{},
		&models.AccountPayable{},
		&models.AccountReceivable{},
		&models.Expense{},
		&models.PaymentRecord{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed default admin user
	seedAdmin(cfg)
}

func seedAdmin(cfg *config.Config) {
	var count int64
	DB.Model(&models.User{}).Where("username = ?", cfg.AdminUser).Count(&count)
	if count > 0 {
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	admin := models.User{
		Username: cfg.AdminUser,
		Password: string(hashedPass),
		RealName: "系统管理员",
		Role:     "admin",
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	fmt.Println("Default admin user created: admin / admin123")
}
