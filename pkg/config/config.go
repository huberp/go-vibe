package config

import (
	"log"

	"github.com/spf13/viper"
)

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	URL             string `mapstructure:"DATABASE_URL"`
	MaxOpenConns    int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime int    `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

// Config holds application configuration
type Config struct {
	Database   DatabaseConfig `mapstructure:",squash"`
	JWTSecret  string         `mapstructure:"JWT_SECRET"`
	ServerPort string         `mapstructure:"SERVER_PORT"`
}

// Load reads configuration from environment variables using Viper
func Load() *Config {
	// Set default values
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/myapp?sslmode=disable")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 30)

	// Automatically read from environment variables
	viper.AutomaticEnv()

	// Unmarshal configuration into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal configuration: %v", err)
	}

	return &config
}
