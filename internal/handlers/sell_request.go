package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/middleware"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/models"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/services"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/utils"
	"github.com/gin-gonic/gin"
)

// NewSellRequestHandler returns a handler function that uses the WhatsApp service from the context
func NewSellRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get services from context
		db, exists := middleware.GetDB(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
			return
		}
		whatsappService, exists := middleware.GetWhatsApp(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "WhatsApp service not available"})
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

func HandleUserLogs(c *gin.Context) {
	db, exists := middleware.GetDB(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var LogsWebhookPayload models.LogsWebhookPayload
	if err := c.ShouldBindJSON(&LogsWebhookPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "Invalid Json format",
			"errorMessage": err,
		})
		return
	}

	logRequestData := LogsWebhookPayload.Record

	if logRequestData.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id is empty",
		})
		return
	}

	// Fetch user data
	userData, err := utils.GetUserDataById(logRequestData.UserID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch userData",
			"details": err.Error(),
		})
		return
	}

	// Fetch property data if PropertyID is provided
	var propertyData models.Property
	var propertyErr error

	if logRequestData.PropertyID != nil {
		propertyData, propertyErr = utils.GetPropertyDataById(*logRequestData.PropertyID, db)
		if propertyErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch propertyData",
				"details": propertyErr.Error(),
			})
			return
		}
	}

	// Return the response
	response := gin.H{
		"message":   "User logs processed successfully",
		"user_data": userData,
	}

	// Only include property data if it was fetched successfully
	if logRequestData.PropertyID != nil && propertyErr == nil {
		response["property_data"] = propertyData
	}

	// Handle different event types
	switch logRequestData.EventType {
	case models.CallPressed, models.WhatsAppPressed:
		log.Printf("User %s initiated contact for property %v (event: %s)", userData.Name, *logRequestData.PropertyID, logRequestData.EventType)
		response["action"] = "Property contact initiated"
		if logRequestData.EventType.RequiresPropertyID() && logRequestData.PropertyID == nil {
			response["warning"] = "Event requires a property ID, but none was provided."
		}

	case models.ConstructionCallPressed, models.ConstructionWhatsAppPressed:
		log.Printf("User %s interacted with construction item (event: %s)", userData.Name, logRequestData.EventType)
		response["action"] = "Construction interaction"

	case models.PostRentalPropertyPressed:
		log.Printf("User %s pressed post rental property button", userData.Name)
		response["action"] = "Post rental property button pressed"
		RentalPropertyPost(c, userData)

	case models.CustomPropertySearchRequest:
		log.Printf("User %s made a custom property search request", userData.Name)
		response["action"] = "Custom property search request made"
		CustomPropertySearch(c, userData)

	default:
		log.Printf("Unknown event type: %s for user %s", logRequestData.EventType, userData.Name)
		response["action"] = "Unknown event type"
		if !logRequestData.EventType.IsValid() {
			response["warning"] = "Event type not recognized"
		}
	}

	c.JSON(http.StatusOK, response)
}

func RentalPropertyPost(c *gin.Context, user models.User) {
	whatsappService, exists := middleware.GetWhatsApp(c)
	if !exists {
		// Handle error: service not found
		return
	}

	phoneNumber := "919480382078"

	internalWAMessage := fmt.Sprintf(` *New Rental property post Received*
	👤 *Name:* %s
	📞 *Phone:* %s
	`, user.Name, user.Phone)

	customerName := "Sir/Madam"
	if user.Name != "" {
		customerName = user.Name
	}

	userFacingMessage := fmt.Sprintf(`Hello %s,

We're excited to hear that you're looking to post your property for rent on our platform. To make this happen.
Please share the following details with us.

- Photos (2 to 3) 
- Google Map Location 
- A short description (e.g., 2BHK ground floor) 
- Monthly Rent 
- Contact & WhatsApp phone number 

Thank you! We're looking forward to helping you with this!
Best regards,
Easyplots Team
	`, customerName)

	whatsappService.SendGroupMessage(c.Request.Context(), services.InternalGroupWhatsAppId, internalWAMessage)
	whatsappService.SendMessage(c.Request.Context(), phoneNumber, userFacingMessage)
}

func CustomPropertySearch(c *gin.Context, user models.User) {
	whatsappService, exists := middleware.GetWhatsApp(c)
	if !exists {
		// Handle error: service not found
		return
	}

	phoneNumber := "919480382078"

	internalWAMessage := fmt.Sprintf(`🤷 *Custom Property Request*
	👤 *Name:* %s
	📞 *Phone:* %s
	`, user.Name, user.Phone)

	customerName := "Sir/Madam"
	if user.Name != "" {
		customerName = user.Name
	}

	userFacingMessage := fmt.Sprintf(`Hello %s,
We noticed you had some trouble finding properties on our platform that suit your needs. 

We’d love to help! If you could share your requirements with us, our team will be more than happy to search for properties that match what you're looking for.
Warm regards,
The Easyplots Team
	`, customerName)

	whatsappService.SendGroupMessage(c.Request.Context(), services.InternalGroupWhatsAppId, internalWAMessage)
	whatsappService.SendMessage(c.Request.Context(), phoneNumber, userFacingMessage)
}
