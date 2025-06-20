package models

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

// WebhookPayload represents the outer JSON structure for webhook requests
type WebhookPayload struct {
	Type      string      `json:"type"`
	Table     string      `json:"table"`
	Record    SellRequest `json:"record"`
	Schema    string      `json:"schema"`
	OldRecord interface{} `json:"old_record"`
}
