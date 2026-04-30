package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/config"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/middlewares"
	"github.com/z876730060/buydemo/models"
)

var backupDir = "./data/backups"

func init() {
	os.MkdirAll(backupDir, 0755)
}

// ========== System Settings ==========

func GetSystemSettings(c *gin.Context) {
	var settings []models.SystemSetting
	database.DB.Find(&settings)

	data := make(map[string]string)
	for _, s := range settings {
		data[s.Key] = s.Value
	}

	// Ensure defaults exist
	defaults := map[string]string{
		"company_name":       "Buy-Demo ERP",
		"company_address":    "",
		"company_phone":      "",
		"low_stock_threshold": "10",
		"currency":           "¥",
		"auto_backup":        "false",
		"backup_retention":   "30",
	}
	for k, v := range defaults {
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func UpdateSystemSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}

	for key, value := range req {
		var setting models.SystemSetting
		if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
			// Create
			database.DB.Create(&models.SystemSetting{
				Key:   key,
				Value: value,
				Desc:  getSettingDesc(key),
			})
		} else {
			database.DB.Model(&setting).Update("value", value)
		}
	}

	middlewares.SimpleLog(c, "update", "system_setting", 0, "更新系统配置")
	c.JSON(http.StatusOK, gin.H{"message": "保存成功"})
}

func getSettingDesc(key string) string {
	descs := map[string]string{
		"company_name":        "公司名称",
		"company_address":     "公司地址",
		"company_phone":       "联系电话",
		"low_stock_threshold": "低库存预警阈值",
		"currency":            "货币符号",
		"auto_backup":         "自动备份",
		"backup_retention":    "备份保留天数",
	}
	if d, ok := descs[key]; ok {
		return d
	}
	return key
}

// ========== Database Backup & Restore ==========

func GetBackupRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	var total int64
	database.DB.Model(&models.BackupRecord{}).Count(&total)

	var records []models.BackupRecord
	database.DB.Preload("Creator").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"data":  records,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func CreateBackup(c *gin.Context) {
	dbPath := config.Load().DBPath

	// Create backup dir
	os.MkdirAll(backupDir, 0755)

	// Build backup filename
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("buydemo_backup_%s.db", timestamp)
	filePath := filepath.Join(backupDir, fileName)

	userID := c.GetUint("user_id")

	backup := models.BackupRecord{
		FileName:  fileName,
		FilePath:  filePath,
		Status:    "success",
		CreatedBy: userID,
	}

	// Copy database file
	srcFile, err := os.Open(dbPath)
	if err != nil {
		backup.Status = "failed"
		backup.Remark = "无法打开数据库文件: " + err.Error()
		database.DB.Create(&backup)
		c.JSON(http.StatusInternalServerError, gin.H{"error": backup.Remark})
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(filePath)
	if err != nil {
		backup.Status = "failed"
		backup.Remark = "无法创建备份文件: " + err.Error()
		database.DB.Create(&backup)
		c.JSON(http.StatusInternalServerError, gin.H{"error": backup.Remark})
		return
	}
	defer dstFile.Close()

	copied, err := io.Copy(dstFile, srcFile)
	if err != nil {
		backup.Status = "failed"
		backup.Remark = "备份写入失败: " + err.Error()
		database.DB.Create(&backup)
		c.JSON(http.StatusInternalServerError, gin.H{"error": backup.Remark})
		return
	}

	backup.FileSize = copied
	database.DB.Create(&backup)

	middlewares.SimpleLog(c, "backup", "database", backup.ID, fmt.Sprintf("数据库备份: %s (%s)", fileName, formatSize(copied)))

	c.JSON(http.StatusOK, gin.H{"data": backup, "message": "备份成功"})
}

func DownloadBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var record models.BackupRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份记录不存在"})
		return
	}

	if _, err := os.Stat(record.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份文件不存在"})
		return
	}

	middlewares.SimpleLog(c, "download_backup", "database", record.ID, "下载备份: "+record.FileName)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", record.FileName))
	c.Header("Content-Type", "application/octet-stream")
	c.File(record.FilePath)
}

func DeleteBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var record models.BackupRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份记录不存在"})
		return
	}

	// Delete file
	os.Remove(record.FilePath)
	database.DB.Delete(&record)

	middlewares.SimpleLog(c, "delete_backup", "database", record.ID, "删除备份: "+record.FileName)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func RestoreBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var record models.BackupRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份记录不存在"})
		return
	}

	if _, err := os.Stat(record.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份文件不存在"})
		return
	}

	dbPath := config.Load().DBPath

	// Close current DB connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据库连接失败"})
		return
	}
	sqlDB.Close()

	// Restore: copy backup to db path
	srcFile, err := os.Open(record.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取备份文件"})
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dbPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法写入数据库文件"})
		return
	}
	defer dstFile.Close()

	io.Copy(dstFile, srcFile)

	// Re-initialize database
	database.Init(config.Load())

	middlewares.SimpleLog(c, "restore", "database", record.ID, "恢复备份: "+record.FileName)
	c.JSON(http.StatusOK, gin.H{"message": "恢复成功，系统已从备份重启"})
}

func formatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	} else {
		return fmt.Sprintf("%.1f MB", float64(size)/1024/1024)
	}
}
