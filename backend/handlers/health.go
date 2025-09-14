package handlers

import (
	"github.com/gin-gonic/gin"
)

// Health 健康检查
func Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
