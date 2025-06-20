package models

// EventType represents the type of user event
type EventType string

// Event type constants (enum-like)
const (
	CallPressed                    EventType = "CALL_PRESSED_PROPERTY"
	WhatsAppPressed                EventType = "WHATS_APP_PRESSED_PROPERTY"
	PostRentalPropertyPressed      EventType = "POST_RENTAL_PROPERTY_PRESSED"
	ConstructionCallPressed        EventType = "CONSTRUCTION_CALL_PRESSED"
	ConstructionWhatsAppPressed    EventType = "CONSTRUCTION_WHATS_APP_PRESSED"
	ConstructionBrochureDownloaded EventType = "CONSTRUCTION_BROCHURE_DOWNLOADED"
	CustomPropertySearchRequest    EventType = "CUSTOM_PROPERTY_SEARCH_REQUEST"
)

// String returns the string representation of the EventType
func (e EventType) String() string {
	return string(e)
}

// IsValid checks if the event type is valid
func (e EventType) IsValid() bool {
	switch e {
	case CallPressed, WhatsAppPressed, PostRentalPropertyPressed,
		ConstructionCallPressed, ConstructionWhatsAppPressed,
		ConstructionBrochureDownloaded, CustomPropertySearchRequest:
		return true
	default:
		return false
	}
}

// RequiresPropertyID returns true if this event type requires a property ID
func (e EventType) RequiresPropertyID() bool {
	switch e {
	case CallPressed, WhatsAppPressed:
		return true
	default:
		return false
	}
}

// GetCategory returns the category of the event
func (e EventType) GetCategory() string {
	switch e {
	case CallPressed, WhatsAppPressed:
		return "property_interaction"
	case ConstructionCallPressed, ConstructionWhatsAppPressed, ConstructionBrochureDownloaded:
		return "construction"
	case PostRentalPropertyPressed:
		return "rental"
	case CustomPropertySearchRequest:
		return "search"
	default:
		return "unknown"
	}
}

type SellRequest struct {
	Id                int    `json:"id"`
	Notes             string `json:"notes"`
	Price             string `json:"price"`
	Address           string `json:"address"`
	UserID            string `json:"user_id" binding:"required"`
	Attended          bool   `json:"attended"`
	AssignTo          string `json:"assign_to"`
	CreatedAt         string `json:"created_at"`
	PropertyType      string `json:"property_type"`
	LastCommunicated  string `json:"last_communicated"`
	ReceivedProperty  bool   `json:"recieved_property"`
	CommunicationMode string `json:"communication mode"`
}

type LogsRequest struct {
	Id               int       `json:"id"`
	CreatedAt        string    `json:"created_at"`
	UserID           string    `json:"user_id"`
	EventType        EventType `json:"event_type"`
	PropertyID       *int      `json:"property_id"`
	Attended         *bool     `json:"attended"`
	Notes            *string   `json:"notes"`
	AssignedTo       *string   `json:"assigned_to"`
	LastCommunicated *string   `json:"last_communicated"`
}

// WebhookPayload represents the outer JSON structure for webhook requests
type WebhookPayload struct {
	Type      string      `json:"type"`
	Table     string      `json:"table"`
	Record    SellRequest `json:"record"`
	Schema    string      `json:"schema"`
	OldRecord interface{} `json:"old_record"`
}

type LogsWebhookPayload struct {
	Type      string      `json:"type"`
	Table     string      `json:"table"`
	Record    LogsRequest `json:"record"`
	Schema    string      `json:"schema"`
	OldRecord interface{} `json:"old_record"`
}
