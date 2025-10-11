package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
}

// Load reads configuration from environment variables using Viper
func Load() *Config {
	// Set default values
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/myapp?sslmode=disable")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("SERVER_PORT", "8080")

	// Automatically read from environment variables
	viper.AutomaticEnv()

	// Unmarshal configuration into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal configuration: %v", err)
	}

	return &config
}
