package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler handles ping requests
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
