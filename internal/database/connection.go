package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnection creates a new database connection pool
func NewConnection(databaseURL string) (*pgxpool.Pool, error) {
	fmt.Println("Connecting to database...")

	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Test the connection
	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Database connected'").Scan(&greeting)
	if err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("database connection test failed: %v", err)
	}

	fmt.Println("Database connection successful:", greeting)
	return dbpool, nil
}

// Close closes the database connection pool
func Close(dbpool *pgxpool.Pool) {
	if dbpool != nil {
		dbpool.Close()
	}
}
