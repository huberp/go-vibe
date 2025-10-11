package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("should load configuration from environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("DATABASE_URL", "postgres://testuser:testpass@testhost:5432/testdb")
		os.Setenv("JWT_SECRET", "test-jwt-secret")
		os.Setenv("SERVER_PORT", "9090")
		defer func() {
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("SERVER_PORT")
		}()

		cfg := Load()

		assert.Equal(t, "postgres://testuser:testpass@testhost:5432/testdb", cfg.DatabaseURL)
		assert.Equal(t, "test-jwt-secret", cfg.JWTSecret)
		assert.Equal(t, "9090", cfg.ServerPort)
	})

	t.Run("should use default values when environment variables are not set", func(t *testing.T) {
		// Ensure environment variables are not set
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("SERVER_PORT")

		cfg := Load()

		assert.Equal(t, "postgres://user:password@localhost:5432/myapp?sslmode=disable", cfg.DatabaseURL)
		assert.Equal(t, "your-secret-key", cfg.JWTSecret)
		assert.Equal(t, "8080", cfg.ServerPort)
	})

	t.Run("should allow partial environment variable override", func(t *testing.T) {
		// Set only some environment variables
		os.Setenv("JWT_SECRET", "custom-secret")
		defer os.Unsetenv("JWT_SECRET")

		// Ensure others are not set
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("SERVER_PORT")

		cfg := Load()

		assert.Equal(t, "postgres://user:password@localhost:5432/myapp?sslmode=disable", cfg.DatabaseURL)
		assert.Equal(t, "custom-secret", cfg.JWTSecret)
		assert.Equal(t, "8080", cfg.ServerPort)
	})
}
