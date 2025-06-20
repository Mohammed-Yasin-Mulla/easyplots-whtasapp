package routes

import (
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/handlers"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/middleware"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, dbpool *pgxpool.Pool, whatsappService *services.WhatsAppService) {
	router.GET("/ping", handlers.PingHandler)

	protectedRoute := router.Group("/")
	// middle-ware
	protectedRoute.Use(middleware.DatabaseMiddleware(dbpool))
	// Webhook endpoints
	protectedRoute.POST("/sell-request", handlers.NewSellRequestHandler(whatsappService))
}
