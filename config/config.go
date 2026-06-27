package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the environment configuration for the application.
type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
}

// LoadConfig loads environment configurations from .env file or system environment variables.
func LoadConfig() *Config {
	// Load .env file if it exists; in production environments, environment variables
	// are typically injected directly, so we do not fail if .env is missing.
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_fallback_jwt_secret_key"
		log.Println("WARNING: JWT_SECRET environment variable is not set. Using insecure fallback key!")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "tickets.db"
	}

	return &Config{
		Port:        port,
		JWTSecret:   jwtSecret,
		DatabaseURL: dbURL,
	}
}
