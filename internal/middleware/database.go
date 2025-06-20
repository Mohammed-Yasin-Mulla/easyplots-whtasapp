package middleware

import (
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/config"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseMiddleware injects database connection into context
func DatabaseMiddleware(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", dbpool)
		c.Next()
	}
}

// ConfigMiddleware injects configuration into context
func ConfigMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	}
}

// WhatsAppMiddleware injects the WhatsApp service into the context
func WhatsAppMiddleware(whatsappService *services.WhatsAppService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("whatsapp", whatsappService)
		c.Next()
	}
}

// GetDB retrieves database connection from context
func GetDB(c *gin.Context) (*pgxpool.Pool, bool) {
	db, exists := c.Get("db")
	if !exists {
		return nil, false
	}
	dbpool, ok := db.(*pgxpool.Pool)
	return dbpool, ok
}

// GetConfig retrieves configuration from context
func GetConfig(c *gin.Context) (*config.Config, bool) {
	cfg, exists := c.Get("config")
	if !exists {
		return nil, false
	}
	config, ok := cfg.(*config.Config)
	return config, ok
}

// GetWhatsApp retrieves the WhatsApp service from the context
func GetWhatsApp(c *gin.Context) (*services.WhatsAppService, bool) {
	ws, exists := c.Get("whatsapp")
	if !exists {
		return nil, false
	}
	whatsappService, ok := ws.(*services.WhatsAppService)
	return whatsappService, ok
}
