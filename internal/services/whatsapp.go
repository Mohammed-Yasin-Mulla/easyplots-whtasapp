package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	// Import PostgreSQL driver for database/sql
	_ "github.com/lib/pq"
)

type WhatsAppService struct {
	client    *whatsmeow.Client
	container *sqlstore.Container
	logger    waLog.Logger
	config    *config.Config
}

// NewWhatsAppService creates a new WhatsApp service
func NewWhatsAppService(ctx context.Context, cfg *config.Config) (*WhatsAppService, error) {
	// Set up logging
	logger := waLog.Stdout("WhatsApp", "INFO", true)

	// Create database container for WhatsApp data
	container, err := sqlstore.New(ctx, "postgres", cfg.DatabaseURL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp container: %v", err)
	}

	// Get the first device from the container
	// If no device exists, it will create a new one
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	// Create WhatsApp client
	client := whatsmeow.NewClient(deviceStore, logger)

	service := &WhatsAppService{
		client:    client,
		container: container,
		logger:    logger,
		config:    cfg,
	}

	// Add event handlers
	client.AddEventHandler(service.eventHandler)

	return service, nil
}

// Connect connects to WhatsApp using configured pairing method
func (w *WhatsAppService) Connect(ctx context.Context) error {
	// Check if already logged in
	if w.client.Store.ID == nil {
		// Not logged in, need to pair device
		log.Printf("WhatsApp: Device not paired. Using %s pairing mode.", w.config.WhatsAppPairingMode)

		// Connect to WhatsApp first
		err := w.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}

		// Choose pairing method based on configuration
		switch w.config.WhatsAppPairingMode {
		case "phone":
			return w.pairWithPhoneNumber(ctx)
		case "qr":
			return w.pairWithQR(ctx)
		default:
			log.Printf("Unknown pairing mode '%s', defaulting to phone pairing", w.config.WhatsAppPairingMode)
			return w.pairWithPhoneNumber(ctx)
		}
	} else {
		// Already logged in, just connect
		log.Println("WhatsApp: Device already paired, connecting...")
		err := w.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
	}

	log.Println("WhatsApp: Connected successfully")
	return nil
}

// pairWithPhoneNumber pairs the device using phone number and pairing code
func (w *WhatsAppService) pairWithPhoneNumber(ctx context.Context) error {
	// Use configured phone number or fallback
	phoneNumber := w.config.WhatsAppPhoneNumber
	if phoneNumber == "" {
		phoneNumber = "919035577330" // Fallback to hardcoded number
		log.Printf("No WHATSAPP_PHONE_NUMBER configured, using fallback: %s", phoneNumber)
	}

	log.Printf("Requesting pairing code for phone number: %s", phoneNumber)

	// Request pairing code
	code, err := w.client.PairPhone(ctx, phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		return fmt.Errorf("failed to request pairing code: %v", err)
	}

	log.Printf("✅ Pairing code: %s", code)
	log.Println("📱 Enter this code in your WhatsApp mobile app:")
	log.Println("   WhatsApp → Settings → Linked Devices → Link a Device → Link with Phone Number")
	log.Println("   Then enter the code above")

	// Wait for pairing to complete
	log.Println("⏳ Waiting for pairing to complete...")

	// Wait for the client to be logged in
	for i := 0; i < 60; i++ { // Wait up to 60 seconds
		if w.client.Store.ID != nil {
			log.Println("🎉 Successfully paired with phone number!")
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("pairing timeout - please try again")
}

// pairWithQR pairs the device using QR code method
func (w *WhatsAppService) pairWithQR(ctx context.Context) error {
	log.Println("WhatsApp: Generating QR code for pairing...")

	// Generate QR code for pairing
	qrChan, err := w.client.GetQRChannel(ctx)
	if err != nil {
		return fmt.Errorf("failed to get QR channel: %v", err)
	}

	// Wait for QR code and print it
	for evt := range qrChan {
		if evt.Event == "code" {
			fmt.Println("📱 QR Code:")
			fmt.Println(evt.Code)
			fmt.Println("Scan this QR code with your WhatsApp mobile app")
		} else {
			fmt.Printf("QR Event: %s\n", evt.Event)
			if evt.Event == "success" {
				fmt.Println("🎉 Successfully paired!")
				break
			}
		}
	}

	return nil
}

// SendMessage sends a WhatsApp message to the specified phone number
func (w *WhatsAppService) SendMessage(ctx context.Context, phoneNumber, message string) error {
	log.Printf("WhatsApp Debug: SendMessage called with phoneNumber: %s", phoneNumber)

	// Ensure client is connected
	if !w.client.IsConnected() {
		log.Printf("WhatsApp Error: Client is not connected")
		return fmt.Errorf("WhatsApp client is not connected")
	}

	log.Printf("WhatsApp Debug: Client is connected, proceeding with message send")

	// Parse phone number to JID
	// Phone number should be in format: country_code + number (e.g., "919999999999")
	jid, err := types.ParseJID(phoneNumber + "@s.whatsapp.net")
	if err != nil {
		log.Printf("WhatsApp Error: Invalid phone number format: %v", err)
		return fmt.Errorf("invalid phone number format: %v", err)
	}

	log.Printf("WhatsApp Debug: JID parsed successfully: %s", jid)

	// Create message
	msg := &waE2E.Message{
		Conversation: &message,
	}

	log.Printf("WhatsApp Debug: Message created, sending to %s", jid)

	// Send message
	response, err := w.client.SendMessage(ctx, jid, msg)
	if err != nil {
		log.Printf("WhatsApp Error: Failed to send message: %v", err)
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("WhatsApp: Message sent successfully. ID: %s, Timestamp: %s",
		response.ID, response.Timestamp)

	return nil
}

// SendGroupMessage sends a WhatsApp message to a group
func (w *WhatsAppService) SendGroupMessage(ctx context.Context, groupJID, message string) error {
	log.Printf("WhatsApp Debug: SendGroupMessage called with groupJID: %s", groupJID)

	// Ensure client is connected
	if !w.client.IsConnected() {
		log.Printf("WhatsApp Error: Client is not connected")
		return fmt.Errorf("WhatsApp client is not connected")
	}

	log.Printf("WhatsApp Debug: Client is connected, proceeding with group message send")

	// Parse group JID
	// Group JID should be in format: groupId@g.us
	jid, err := types.ParseJID(groupJID)
	if err != nil {
		log.Printf("WhatsApp Error: Invalid group JID format: %v", err)
		return fmt.Errorf("invalid group JID format: %v", err)
	}

	log.Printf("WhatsApp Debug: Group JID parsed successfully: %s", jid)

	// Create message
	msg := &waE2E.Message{
		Conversation: &message,
	}

	log.Printf("WhatsApp Debug: Group message created, sending to %s", jid)

	// Send message
	response, err := w.client.SendMessage(ctx, jid, msg)
	if err != nil {
		log.Printf("WhatsApp Error: Failed to send group message: %v", err)
		return fmt.Errorf("failed to send group message: %v", err)
	}

	log.Printf("WhatsApp: Group message sent successfully. ID: %s, Timestamp: %s",
		response.ID, response.Timestamp)

	return nil
}

const InternalGroupWhatsAppId = "120363420697230363@g.us"

// SendSellRequestToGroup sends a sell request message to a specific group
func (w *WhatsAppService) SendSellRequestToGroup(ctx context.Context, groupJID, userName, propertyType, address, price, userPhone string) error {
	log.Printf("WhatsApp Debug: SendSellRequestToGroup called - GroupJID: %s, User: %s", groupJID, userName)

	message := fmt.Sprintf(`🏠 *New Sell Request Received*
👤 *Name:* %s
🏘️ *Property Type:* %s  
📍 *Address:* %s
💰 *Price:* %s
📞 *Phone:* %s
`, userName, propertyType, address, price, userPhone)

	log.Printf("WhatsApp Debug: Group message template created, length: %d characters", len(message))

	return w.SendGroupMessage(ctx, groupJID, message)
}

// Event handler for WhatsApp events
func (w *WhatsAppService) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Handle incoming messages (optional)
		w.logger.Infof("Received message from %s: %s", v.Info.Sender, v.Message.GetConversation())
	case *events.Connected:
		w.logger.Infof("WhatsApp connected")
	case *events.Disconnected:
		w.logger.Infof("WhatsApp disconnected")
	case *events.LoggedOut:
		w.logger.Infof("WhatsApp logged out")
	}
}

// Disconnect closes the WhatsApp connection
func (w *WhatsAppService) Disconnect() {
	if w.client != nil {
		w.client.Disconnect()
	}
}

// Close closes the WhatsApp service and database connections
func (w *WhatsAppService) Close() error {
	w.Disconnect()
	if w.container != nil {
		return w.container.Close()
	}
	return nil
}

func SellRequestWhatAppMessage(userName string) string {
	greeting := "Hello Sir/Madam,"
	if userName != "" {
		greeting = fmt.Sprintf("Hello %s,", userName)
	}

	return fmt.Sprintf(`%s
We have received your property posting request 
Please share the following details to allow us to list your property 
- Property size
- Type of property (NA, Gunta etc...)
- Facing of the property
- Google map location 
- 2-3 images of your property
- Size of the property 
- Utraa Copy (for internal records keeping)

If you have any quires please feel free to call us.

Best regards,
Easyplots Team`, greeting)
}
