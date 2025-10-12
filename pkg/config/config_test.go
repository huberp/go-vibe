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

		assert.Equal(t, "postgres://testuser:testpass@testhost:5432/testdb", cfg.Database.URL)
		assert.Equal(t, "test-jwt-secret", cfg.JWT.Secret)
		assert.Equal(t, "9090", cfg.Server.Port)
	})

	t.Run("should use default values when environment variables are not set", func(t *testing.T) {
		// Ensure environment variables are not set
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("APP_STAGE")

		cfg := Load()

		// Should use defaults from base.yaml or fallback defaults
		assert.NotEmpty(t, cfg.Database.URL)
		assert.NotEmpty(t, cfg.JWT.Secret)
		assert.Equal(t, "8080", cfg.Server.Port)
	})

	t.Run("should allow partial environment variable override", func(t *testing.T) {
		// Set only some environment variables
		os.Setenv("JWT_SECRET", "custom-secret")
		defer os.Unsetenv("JWT_SECRET")

		// Ensure others are not set
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("SERVER_PORT")

		cfg := Load()

		assert.NotEmpty(t, cfg.Database.URL)
		assert.Equal(t, "custom-secret", cfg.JWT.Secret)
		assert.Equal(t, "8080", cfg.Server.Port)
	})
}

func TestLoadWithStage(t *testing.T) {
	t.Run("should load development stage", func(t *testing.T) {
		cfg := LoadWithStage("development")
		assert.NotNil(t, cfg)
		assert.Equal(t, "8080", cfg.Server.Port)
		assert.NotEmpty(t, cfg.Database.URL)
		assert.NotEmpty(t, cfg.JWT.Secret)
	})

	t.Run("should load production stage", func(t *testing.T) {
		// Set required env vars for production
		os.Setenv("DATABASE_URL", "postgres://prod:pass@prod-host:5432/prod")
		os.Setenv("JWT_SECRET", "prod-secret")
		defer func() {
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("JWT_SECRET")
		}()

		cfg := LoadWithStage("production")
		assert.NotNil(t, cfg)
		assert.Equal(t, "postgres://prod:pass@prod-host:5432/prod", cfg.Database.URL)
		assert.Equal(t, "prod-secret", cfg.JWT.Secret)
	})

	t.Run("should load staging stage", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "staging-secret")
		defer os.Unsetenv("JWT_SECRET")

		cfg := LoadWithStage("staging")
		assert.NotNil(t, cfg)
		assert.Equal(t, "staging-secret", cfg.JWT.Secret)
	})

	t.Run("should support environment variable overrides", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "9090")
		defer os.Unsetenv("SERVER_PORT")

		cfg := LoadWithStage("development")
		assert.Equal(t, "9090", cfg.Server.Port)
	})
}

func TestGetStage(t *testing.T) {
	t.Run("should default to development stage", func(t *testing.T) {
		os.Unsetenv("APP_STAGE")
		stage := getStage()
		assert.Equal(t, "development", stage)
	})

	t.Run("should read stage from environment variable", func(t *testing.T) {
		os.Setenv("APP_STAGE", "production")
		defer os.Unsetenv("APP_STAGE")

		stage := getStage()
		assert.Equal(t, "production", stage)
	})

	t.Run("should read staging from environment variable", func(t *testing.T) {
		os.Setenv("APP_STAGE", "staging")
		defer os.Unsetenv("APP_STAGE")

		stage := getStage()
		assert.Equal(t, "staging", stage)
	})
}

