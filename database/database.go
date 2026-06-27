package database

import (
	"log"
	"ticket-system/config"
	"ticket-system/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance.
var DB *gorm.DB

// InitDB initializes the SQLite connection using GORM and runs AutoMigrate on models.
func InitDB(cfg *config.Config) *gorm.DB {
	var err error

	// Connect to the SQLite database. GORM is configured with Info logs for traceability.
	DB, err = gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite database at '%s': %v", cfg.DatabaseURL, err)
	}

	log.Printf("Database connection established at '%s'\n", cfg.DatabaseURL)

	// Automatically run migrations to sync schemas.
	err = DB.AutoMigrate(&models.User{}, &models.Ticket{})
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	log.Println("Database schema auto-migrations completed successfully")

	return DB
}
