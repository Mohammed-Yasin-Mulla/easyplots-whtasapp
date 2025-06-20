package models

import "time"

type Property struct {
	ID              int64     `json:"id" db:"id"`
	CategoryID      *int64    `json:"category_id" db:"category_id"`
	Title           string    `json:"title" db:"title"`
	Size            string    `json:"size" db:"size"`
	DeveloperID     *int64    `json:"developer_id" db:"developer_id"`
	OwnerID         *string   `json:"owner_id" db:"owner_id"`
	AddressID       int64     `json:"address_id" db:"address_id"`
	Recommended     bool      `json:"recommended" db:"recommended"`
	BannerID        *int64    `json:"banner_id" db:"banner_id"`
	CustomPhoneNo   *string   `json:"custom_phone_no" db:"custom_phone_no"`
	Status          string    `json:"status" db:"status"`
	Facing          *string   `json:"facing" db:"facing"`
	EstimatedPrice  *int      `json:"estimated_price" db:"estimated_price"`
	Negotiable      bool      `json:"negotiable" db:"negotiable"`
	MapCenterpoint  *string   `json:"map_centerpoint" db:"map_centerpoint"`
	CustomZoom      *int16    `json:"custom_zoom" db:"custom_zoom"`
	Featured        bool      `json:"featured" db:"featured"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	ShowAsNew       bool      `json:"show_as_new" db:"show_as_new"`
	VisibilityScore string    `json:"visibility_score" db:"visibility_score"`
	Rental          bool      `json:"rental" db:"rental"`
	RentAmount      *float64  `json:"rent_amount" db:"rent_amount"`
	RevealLocation  bool      `json:"reveal_location" db:"reveal_location"`
}
