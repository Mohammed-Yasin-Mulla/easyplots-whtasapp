package handlers

import (
	"net/http"

	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/middleware"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/models"
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

	// Insert the sell request into database
	query := `INSERT INTO sell_requests (user_id, notes, price, address, property_type) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var insertedID int
	err := db.QueryRow(c, query,
		sellRequestData.UserID,
		sellRequestData.Notes,
		sellRequestData.Price,
		sellRequestData.Address,
		sellRequestData.PropertyType).Scan(&insertedID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to insert sell request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sell request saved successfully",
		"data":    sellRequestData,
		"id":      insertedID,
	})
}
