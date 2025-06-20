package middleware

import (
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

// GetDB retrieves database connection from context
func GetDB(c *gin.Context) (*pgxpool.Pool, bool) {
	db, exists := c.Get("db")
	if !exists {
		return nil, false
	}
	dbpool, ok := db.(*pgxpool.Pool)
	return dbpool, ok
}
