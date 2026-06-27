package middleware

import (
	"net/http"
	"strings"
	"ticket-system/config"
	"ticket-system/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware intercepts requests to verify the Bearer token and attach userID to context.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check prefix and format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired authentication token"})
			c.Abort()
			return
		}

		// Store user ID in Gin context for controllers
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
