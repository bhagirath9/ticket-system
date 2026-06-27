package main

import (
	"log"
	"ticket-system/config"
	"ticket-system/controllers"
	"ticket-system/database"
	"ticket-system/repository"
	"ticket-system/routes"
	"ticket-system/services"
)

func main() {
	log.Println("Initializing Ticket Management System Backend...")

	// 1. Load Configurations from environment or file
	cfg := config.LoadConfig()

	// 2. Connect to database and perform migrations
	db := database.InitDB(cfg)

	// 3. Instantiate repositories (DB isolation layer)
	userRepo := repository.NewUserRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	// 4. Instantiate services (Business logic layer)
	authService := services.NewAuthService(userRepo, cfg)
	ticketService := services.NewTicketService(ticketRepo)

	// 5. Instantiate controllers (HTTP transport layer)
	authCtrl := controllers.NewAuthController(authService)
	ticketCtrl := controllers.NewTicketController(ticketService)

	// 6. Setup Gin Router and routes mapping
	r := routes.SetupRouter(cfg, authCtrl, ticketCtrl)

	// 7. Start the server
	log.Printf("HTTP server is starting on port %s...\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Fatal: Failed to start server: %v", err)
	}
}
