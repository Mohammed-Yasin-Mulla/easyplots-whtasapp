package utils

import (
	"context"
	"fmt"

	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserDataById(userId string, db *pgxpool.Pool) (models.User, error) {
	var usersData models.User

	// Select specific columns instead of * to be more explicit
	query := `SELECT name, role, is_blocked, id, phone, pref_lang, address, 
			  created_at, push_notification_tokens, notes, send_push_notifications 
			  FROM users WHERE id = $1`

	err := db.QueryRow(context.Background(), query, userId).Scan(
		&usersData.Name,
		&usersData.Role,
		&usersData.IsBlocked,
		&usersData.ID,
		&usersData.Phone,
		&usersData.PrefLang,
		&usersData.Address,
		&usersData.CreatedAt,
		&usersData.PushNotificationTokens,
		&usersData.Notes,
		&usersData.SendPushNotifications,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return usersData, fmt.Errorf("user with ID %s not found", userId)
		}
		return usersData, fmt.Errorf("failed to fetch user: %v", err)
	}

	return usersData, nil
}
