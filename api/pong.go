package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Pong 响应测试
func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
