package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	DatabaseURL string
	JWTSecret   string
	ServerPort  string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/myapp?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
