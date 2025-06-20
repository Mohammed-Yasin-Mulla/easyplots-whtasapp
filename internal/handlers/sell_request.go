package handlers

import (
	"net/http"

	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/middleware"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/models"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/utils"
	"github.com/gin-gonic/gin"
)

// SellRequestHandler handles sell request webhooks
func SellRequestHandler(c *gin.Context) {
	// Get database connection from context
	db, exists := middleware.GetDB(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var webhookPayload models.WebhookPayload
	if err := c.ShouldBindJSON(&webhookPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "Invalid JSON format",
			"errorMessage": err,
		})
		return
	}

	sellRequestData := webhookPayload.Record
	if sellRequestData.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id is empty",
		})
		return
	}

	var userData, err = utils.GetUserDataById(sellRequestData.UserID, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch userData",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Sell request saved successfully",
		"data":      sellRequestData,
		"user_data": userData,
	})
}
