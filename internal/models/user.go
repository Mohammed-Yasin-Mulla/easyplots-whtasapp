package models

import (
	"time"
)

type User struct {
	Name                   string    `json:"name" db:"name"`
	Role                   *string   `json:"role" db:"role"` // nullable
	IsBlocked              bool      `json:"is_blocked" db:"is_blocked"`
	ID                     string    `json:"id" db:"id"`
	Phone                  string    `json:"phone" db:"phone"`
	PrefLang               *string   `json:"pref_lang" db:"pref_lang"` // nullable
	Address                string    `json:"address" db:"address"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	PushNotificationTokens []string  `json:"push_notification_tokens" db:"push_notification_tokens"`
	Notes                  *string   `json:"notes" db:"notes"` // nullable
	SendPushNotifications  bool      `json:"send_push_notifications" db:"send_push_notifications"`
}
