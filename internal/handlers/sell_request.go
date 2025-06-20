package handlers

import (
	"log"
	"net/http"

	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/middleware"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/models"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/services"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/utils"
	"github.com/gin-gonic/gin"
)

// NewSellRequestHandler returns a handler function with WhatsApp service injected
func NewSellRequestHandler(whatsappService *services.WhatsAppService) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Fetch user data to get phone number and name for WhatsApp
		var userData, err = utils.GetUserDataById(sellRequestData.UserID, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch userData",
				"details": err.Error(),
			})
			return
		}

		// Extract name safely (handle potential nil pointer)
		userName := userData.Name
		if userName == "" {
			userName = "Valued Customer"
		}

		log.Printf("WhatsApp Debug: User data - Name: %s, Phone: %s", userName, userData.Phone)

		// For now, using hardcoded number for testing - you can switch to userData.Phone later
		phoneNumber := "919480382078" // Change this to userData.Phone when ready
		log.Printf("WhatsApp Debug: Attempting to send message to: %s", phoneNumber)

		// Check if we should send to group instead of individual
		log.Printf("WhatsApp Debug: Sending to group: %s", services.InternalGroupWhatsAppId)
		err = whatsappService.SendSellRequestToGroup(
			c.Request.Context(),
			services.InternalGroupWhatsAppId,
			userName,
			sellRequestData.PropertyType,
			sellRequestData.Address,
			sellRequestData.Price,
			userData.Phone,
		)

		if err != nil {
			//TODO: Need to add logging system
		}

		err = whatsappService.SendMessage(c.Request.Context(), phoneNumber, services.SellRequestWhatAppMessage(userName))

		var whatsappStatus string
		if err != nil {
			log.Printf("WhatsApp Error: Failed to send message: %v", err)
			whatsappStatus = "Failed to send WhatsApp message: " + err.Error()
		} else {
			log.Printf("WhatsApp Success: Message sent successfully to %s", phoneNumber)
			whatsappStatus = "WhatsApp message sent successfully"
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Sell request received successfully",
			"data":      sellRequestData,
			"user_data": userData,
			"whatsapp":  whatsappStatus,
		})
	}
}

// SellRequestHandler is the original handler (kept for backward compatibility)
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
		"message":   "Sell request received successfully",
		"data":      sellRequestData,
		"user_data": userData,
	})
}
