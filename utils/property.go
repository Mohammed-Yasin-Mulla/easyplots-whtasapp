package utils

import (
	"context"
	"fmt"

	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPropertyDataById(propertyId int, db *pgxpool.Pool) (models.Property, error) {
	var propertyData models.Property

	// Select all columns from property table
	query := `SELECT id, category_id, title, size, developer_id, owner_id, address_id, 
			  recommended, banner_id, custom_phone_no, status, 
			  facing, estimated_price, negotiable, map_centerpoint, custom_zoom, 
			  featured, created_at, show_as_new, visibility_score, rental, 
			  rent_amount, reveal_location 
			  FROM property WHERE id = $1`

	err := db.QueryRow(context.Background(), query, propertyId).Scan(
		&propertyData.ID,
		&propertyData.CategoryID,
		&propertyData.Title,
		&propertyData.Size,
		&propertyData.DeveloperID,
		&propertyData.OwnerID,
		&propertyData.AddressID,
		&propertyData.Recommended,
		&propertyData.BannerID,
		&propertyData.CustomPhoneNo,
		&propertyData.Status,
		&propertyData.Facing,
		&propertyData.EstimatedPrice,
		&propertyData.Negotiable,
		&propertyData.MapCenterpoint,
		&propertyData.CustomZoom,
		&propertyData.Featured,
		&propertyData.CreatedAt,
		&propertyData.ShowAsNew,
		&propertyData.VisibilityScore,
		&propertyData.Rental,
		&propertyData.RentAmount,
		&propertyData.RevealLocation,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return propertyData, fmt.Errorf("property with ID %s not found", propertyId)
		}
		return propertyData, fmt.Errorf("failed to fetch property: %v", err)
	}

	return propertyData, nil
}
