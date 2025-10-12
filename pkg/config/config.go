package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	URL             string `mapstructure:"url"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// Config holds application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// Load reads configuration from YAML files and environment variables using Viper
func Load() *Config {
	return LoadWithStage(getStage())
}

// LoadWithStage loads configuration for a specific stage
func LoadWithStage(stage string) *Config {
	v := viper.New()

	// Set config paths (check multiple locations)
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	v.AddConfigPath("/etc/myapp/config")

	// Load base configuration first
	v.SetConfigName("base")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		// If no config file found, set defaults for backward compatibility
		log.Printf("No base config file found, using defaults: %v", err)
		setDefaults(v)
	}

	// Merge stage-specific configuration
	v.SetConfigName(stage)
	if err := v.MergeInConfig(); err != nil {
		log.Printf("No %s config file found, using base configuration: %v", stage, err)
	}

	// Enable environment variable overrides
	v.AutomaticEnv()
	// Support nested config keys via environment variables (e.g., SERVER_PORT for server.port)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// For backward compatibility, also bind old-style env vars
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("database.url", "DATABASE_URL")
	v.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	v.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	v.BindEnv("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")
	v.BindEnv("jwt.secret", "JWT_SECRET")

	// Unmarshal configuration into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal configuration: %v", err)
	}

	return &config
}

// getStage returns the current stage from environment variable or default
func getStage() string {
	stage := os.Getenv("APP_STAGE")
	if stage == "" {
		stage = "development"
	}
	return stage
}

// setDefaults sets default values for backward compatibility when no config files exist
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", "8080")
	v.SetDefault("database.url", "postgres://user:password@localhost:5432/myapp?sslmode=disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", 30)
	v.SetDefault("jwt.secret", "your-secret-key")
}
