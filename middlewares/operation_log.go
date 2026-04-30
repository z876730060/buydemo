package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/models"
)

func LogOperation(action, target string, targetID uint, detail string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only log on success (status < 400)
		if c.Writer.Status() >= 400 {
			return
		}

		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		ip := c.ClientIP()
		// Get real IP from headers
		if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
			ip = strings.Split(xff, ",")[0]
		} else if xri := c.GetHeader("X-Real-IP"); xri != "" {
			ip = xri
		}

		uid, _ := userID.(uint)
		uname, _ := username.(string)

		log := models.OperationLog{
			UserID:   uid,
			Username: uname,
			Action:   action,
			Target:   target,
			TargetID: targetID,
			Detail:   detail,
			IP:       ip,
		}

		database.DB.Create(&log)
	}
}

// SimpleLog logs an operation without being a middleware
func SimpleLog(c *gin.Context, action, target string, targetID uint, detail string) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	ip := c.ClientIP()
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ip = strings.Split(xff, ",")[0]
	} else if xri := c.GetHeader("X-Real-IP"); xri != "" {
		ip = xri
	}

	uid, _ := userID.(uint)
	uname, _ := username.(string)

	log := models.OperationLog{
		UserID:   uid,
		Username: uname,
		Action:   action,
		Target:   target,
		TargetID: targetID,
		Detail:   detail,
		IP:       ip,
	}

	database.DB.Create(&log)
}
