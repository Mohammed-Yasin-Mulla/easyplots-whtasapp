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

		log.Printf("WhatsApp Debug: Attempting to send message to: %s", userData.Phone)

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

		err = whatsappService.SendMessage(c.Request.Context(), userData.Phone, services.SellRequestWhatAppMessage(userName))

		var whatsappStatus string
		if err != nil {
			log.Printf("WhatsApp Error: Failed to send message: %v", err)
			whatsappStatus = "Failed to send WhatsApp message: " + err.Error()
		} else {
			log.Printf("WhatsApp Success: Message sent successfully to %s", userData.Phone)
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

	// Return the response
	response := gin.H{
		"message":   "User logs processed successfully",
		"user_data": userData,
	}

	// Handle different event types
	switch logRequestData.EventType {
	case models.CallPressed, models.WhatsAppPressed:
		log.Printf("User %s initiated contact for property %v (event: %s)", userData.Name, *logRequestData.PropertyID, logRequestData.EventType)
		response["action"] = "Property contact initiated"
		if logRequestData.EventType.RequiresPropertyID() && logRequestData.PropertyID == nil {
			response["warning"] = "Event requires a property ID, but none was provided."
		}
		PropertyInterest(c, userData, int(*logRequestData.PropertyID))

	case models.ConstructionCallPressed, models.ConstructionWhatsAppPressed:
		log.Printf("User %s interacted with construction item (event: %s)", userData.Name, logRequestData.EventType)
		response["action"] = "Construction interaction"
		ConstructionServicesEnq(c, userData)

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

	internalWAMessage := fmt.Sprintf(` *New Rental property post Received*
	👤 *Name:* %s
	📞 *Phone:* %s
	`, user.Name, user.Phone)

	customerName := "Sir/Madam"
	if user.Name != "" {
		customerName = user.Name
	}

	userFacingMessage := fmt.Sprintf(`Hello %s,

Fantastic! We're thrilled to help you list your rental property on Easyplots and connect you with qualified tenants quickly.

To create a standout listing that gets maximum visibility, please provide the following:

📸 *Photos:* 2-3 clear images of your property
📍 *Location:* A Google Maps link for accuracy
📝 *Description:* A brief summary (e.g., 2BHK, ground floor, key amenities)
💰 *Rent:* The expected monthly rent
📞 *Contact:* Your preferred phone & WhatsApp number

Once we have these details, we'll get your property live for thousands of potential renters to see.

Best,
The Easyplots Team`, customerName)

	whatsappService.SendGroupMessage(c.Request.Context(), services.InternalGroupWhatsAppId, internalWAMessage)
	whatsappService.SendMessage(c.Request.Context(), user.Phone, userFacingMessage)
}

func CustomPropertySearch(c *gin.Context, user models.User) {
	whatsappService, exists := middleware.GetWhatsApp(c)
	if !exists {
		// Handle error: service not found
		return
	}

	internalWAMessage := fmt.Sprintf(`🤷 *Custom Property Request*
	👤 *Name:* %s
	📞 *Phone:* %s
	`, user.Name, user.Phone)

	customerName := "Sir/Madam"
	if user.Name != "" {
		customerName = user.Name
	}

	userFacingMessage := fmt.Sprintf(`Hello %s,

Searching for the perfect property? Let our experts do the heavy lifting for you!

We offer a complimentary custom search service to match you with exclusive listings that meet your exact needs.

Simply reply with your requirements (e.g., location, budget, property type, size), and we'll send you a curated list of the best options available.

We look forward to finding your ideal property.

Warm regards,
The Easyplots Team`, customerName)

	whatsappService.SendGroupMessage(c.Request.Context(), services.InternalGroupWhatsAppId, internalWAMessage)
	whatsappService.SendMessage(c.Request.Context(), user.Phone, userFacingMessage)
}

func ConstructionServicesEnq(c *gin.Context, user models.User) {
	whatsappService, exists := middleware.GetWhatsApp(c)
	if !exists {
		// Handle error: service not found
		return
	}

	internalWAMessage := fmt.Sprintf(`🏗️ *Construction Services*
	👤 *Name:* %s
	📞 *Phone:* %s
	`, user.Name, user.Phone)

	customerName := "Sir/Madam"
	if user.Name != "" {
		customerName = user.Name
	}

	userFacingMessage := fmt.Sprintf(`Hello %s,

Thank you for your interest in our construction services!

Whether you're planning to build your dream home or a new commercial project, our expert team is here to bring your vision to life with quality craftsmanship and on-time delivery.

To provide you with a tailored consultation, could you tell us a bit more about your project?

One of our specialists is ready to connect and discuss how we can help.

Best regards,
The Easyplots Team`, customerName)

	whatsappService.SendGroupMessage(c.Request.Context(), services.InternalGroupWhatsAppId, internalWAMessage)
	whatsappService.SendMessage(c.Request.Context(), user.Phone, userFacingMessage)
}

func PropertyInterest(c *gin.Context, user models.User, propertyId int) {
	whatsappService, waExists := middleware.GetWhatsApp(c)
	db, dbExists := middleware.GetDB(c)
	if !waExists || !dbExists {
		log.Println("Could not get whatsapp service or db from context")
		return
	}

	propertyData, err := utils.GetPropertyDataById(propertyId, db)
	if err != nil {
		log.Printf("Failed to get property data for id %d: %v", propertyId, err)
		// Optionally, notify the group that an error occurred
		errorMsg := fmt.Sprintf("⚠️ Error fetching property details for ID: %d\nUser: %s\nError: %v", propertyId, user.Name, err)
		whatsappService.SendGroupMessage(c.Request.Context(), services.InternalGroupWhatsAppId, errorMsg)
		return
	}

	internalWAMessage := fmt.Sprintf(`✅ *Property Interest*
👤 *Name:* %s
📞 *Phone:* %s

*Property Details:*
ID: %d
Title: %s
Size: %s
Link: https://easyplots.in/property/%d

Message is not sent to the user, please call them directly.`, user.Name, user.Phone, propertyData.ID, propertyData.Title, propertyData.Size, propertyData.ID)

	whatsappService.SendGroupMessage(c.Request.Context(), services.InternalGroupWhatsAppId, internalWAMessage)
}
