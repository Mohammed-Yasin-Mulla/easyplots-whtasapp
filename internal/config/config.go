package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL         string
	Port                string
	GinMode             string
	WhatsAppPairingMode string // "phone" or "qr"
	WhatsAppPhoneNumber string // Phone number for pairing
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or error loading .env file")
	}

	config := &Config{
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		Port:                getEnv("PORT", "8080"),
		GinMode:             getEnv("GIN_MODE", "debug"),
		WhatsAppPairingMode: getEnv("WHATSAPP_PAIRING_MODE", "phone"), // Default to phone pairing
		WhatsAppPhoneNumber: getEnv("WHATSAPP_PHONE_NUMBER", ""),
	}

	return config, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Validate checks if required configuration values are present
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	return nil
}
