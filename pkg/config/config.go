package config

import (
	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	DatabaseURL string
	JWTSecret   string
	ServerPort  string
}

// Load reads configuration from environment variables using Viper
func Load() *Config {
	// Set default values
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/myapp?sslmode=disable")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("SERVER_PORT", "8080")

	// Automatically read from environment variables
	viper.AutomaticEnv()

	return &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		ServerPort:  viper.GetString("SERVER_PORT"),
	}
}
