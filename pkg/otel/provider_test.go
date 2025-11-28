package otel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func TestInitProvider(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()

	t.Run("should initialize provider with default endpoint", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Initialize provider (this will attempt to connect but may fail if collector is not running)
		tp, cleanup, err := InitProvider(ctx, "test-service", "", logger)
		
		// We expect the provider to be created even if the collector is not available
		// The error only occurs during actual trace export, not during initialization
		assert.NoError(t, err, "InitProvider should not error during initialization")
		assert.NotNil(t, tp, "TracerProvider should not be nil")
		assert.NotNil(t, cleanup, "Cleanup function should not be nil")

		// Verify the global tracer provider is set
		globalTP := otel.GetTracerProvider()
		assert.NotNil(t, globalTP, "Global TracerProvider should be set")

		// Cleanup
		if cleanup != nil {
			err := cleanup(context.Background())
			assert.NoError(t, err, "Cleanup should not error")
		}
	})

	t.Run("should initialize provider with custom endpoint", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tp, cleanup, err := InitProvider(ctx, "test-service", "otel-collector:4317", logger)
		
		assert.NoError(t, err, "InitProvider should not error during initialization")
		assert.NotNil(t, tp, "TracerProvider should not be nil")
		assert.NotNil(t, cleanup, "Cleanup function should not be nil")

		// Cleanup
		if cleanup != nil {
			err := cleanup(context.Background())
			assert.NoError(t, err, "Cleanup should not error")
		}
	})

	t.Run("should set global tracer provider", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, cleanup, err := InitProvider(ctx, "test-service", "localhost:4317", logger)
		defer func() {
			if cleanup != nil {
				cleanup(context.Background())
			}
		}()

		assert.NoError(t, err, "InitProvider should not error")
		
		// Get the global tracer provider
		globalTP := otel.GetTracerProvider()
		assert.NotNil(t, globalTP, "Global TracerProvider should be set")

		// Create a tracer to verify it works
		tracer := globalTP.Tracer("test-tracer")
		assert.NotNil(t, tracer, "Tracer should not be nil")
	})
}
