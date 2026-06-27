package routes

import (
	"net/http"
	"ticket-system/config"
	"ticket-system/controllers"
	"ticket-system/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter registers all routing groups and maps them to controllers.
func SetupRouter(cfg *config.Config, authCtrl *controllers.AuthController, ticketCtrl *controllers.TicketController) *gin.Engine {
	r := gin.Default()

	// Apply Gin default recovery middleware to recover from any panics gracefully
	r.Use(gin.Recovery())

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Ticket System",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Health check endpoint (Public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth endpoint group (Public)
	auth := r.Group("/auth")
	{
		auth.POST("/register", authCtrl.Register)
		auth.POST("/login", authCtrl.Login)
	}

	// Tickets endpoint group (Protected by JWT AuthMiddleware)
	tickets := r.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware(cfg))
	{
		tickets.POST("", ticketCtrl.Create)
		tickets.GET("", ticketCtrl.List)
		tickets.GET("/:id", ticketCtrl.GetByID)
		tickets.PATCH("/:id/status", ticketCtrl.UpdateStatus)
	}

	return r
}
