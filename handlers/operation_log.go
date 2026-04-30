package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/models"
)

func GetOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	action := c.Query("action")
	target := c.Query("target")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query := database.DB.Model(&models.OperationLog{})
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if target != "" {
		query = query.Where("target = ?", target)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR detail LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	var logs []models.OperationLog
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// ExportOperationLogs exports all operation logs matching the filter
func ExportOperationLogs(c *gin.Context) {
	action := c.Query("action")
	target := c.Query("target")
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.OperationLog{}).Order("id DESC")
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if target != "" {
		query = query.Where("target = ?", target)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR detail LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var logs []models.OperationLog
	query.Find(&logs)

	// Build CSV
	csv := "\xEF\xBB\xBF" // BOM for Excel UTF-8
	csv += "ID,用户,操作,对象,对象ID,详情,IP,时间\n"
	for _, l := range logs {
		csv += strconv.Itoa(int(l.ID)) + ","
		csv += l.Username + ","
		csv += l.Action + ","
		csv += l.Target + ","
		csv += strconv.Itoa(int(l.TargetID)) + ","
		csv += l.Detail + ","
		csv += l.IP + ","
		csv += l.CreatedAt.Format("2006-01-02 15:04:05") + "\n"
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=operation_logs.csv")
	c.String(http.StatusOK, csv)
}
