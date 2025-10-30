package otel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/zap"
)

// TestOtelWorkflowWithExporter demonstrates the complete OTEL workflow
// including span creation, propagation, and export
func TestOtelWorkflowWithExporter(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()
	
	t.Run("should create and export spans with in-memory exporter", func(t *testing.T) {
		ctx := context.Background()

		// Create an in-memory span exporter for testing
		exporter := tracetest.NewInMemoryExporter()
		
		// Create TracerProvider with in-memory exporter
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSyncer(exporter),
		)
		
		// Set global TracerProvider
		otel.SetTracerProvider(tp)
		defer tp.Shutdown(ctx)

		// Create a tracer
		tracer := tp.Tracer("test-tracer")

		// Create a span
		_, span := tracer.Start(ctx, "test-operation")
		span.SetAttributes(
			attribute.String("test.key", "test.value"),
			attribute.Int("test.count", 42),
		)
		span.SetStatus(codes.Ok, "Operation completed successfully")
		span.End()

		// Force flush to ensure span is exported
		err := tp.ForceFlush(ctx)
		assert.NoError(t, err, "ForceFlush should not error")

		// Verify span was exported
		spans := exporter.GetSpans()
		assert.Len(t, spans, 1, "Should have exported 1 span")
		
		exportedSpan := spans[0]
		assert.Equal(t, "test-operation", exportedSpan.Name, "Span name should match")
		assert.Equal(t, codes.Ok, exportedSpan.Status.Code, "Span status should be OK")
		
		// Verify attributes
		attrs := exportedSpan.Attributes
		assert.Contains(t, attrs, attribute.String("test.key", "test.value"))
		assert.Contains(t, attrs, attribute.Int("test.count", 42))
	})

	t.Run("should trace HTTP requests with otelgin middleware", func(t *testing.T) {
		ctx := context.Background()
		gin.SetMode(gin.TestMode)

		// Create an in-memory span exporter for testing
		exporter := tracetest.NewInMemoryExporter()
		
		// Create TracerProvider with in-memory exporter
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSyncer(exporter),
		)
		
		// Set global TracerProvider
		otel.SetTracerProvider(tp)
		defer tp.Shutdown(ctx)

		// Create Gin router with otelgin middleware
		router := gin.New()
		router.Use(otelgin.Middleware("test-service"))

		// Add a test endpoint
		router.GET("/test", func(c *gin.Context) {
			// Get tracer from context and create a child span
			tracer := otel.Tracer("test-handler")
			_, span := tracer.Start(c.Request.Context(), "process-request")
			span.SetAttributes(attribute.String("handler", "test"))
			span.End()

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Make a test request
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		// Force flush to ensure spans are exported
		err := tp.ForceFlush(ctx)
		assert.NoError(t, err, "ForceFlush should not error")

		// Verify spans were exported
		spans := exporter.GetSpans()
		assert.GreaterOrEqual(t, len(spans), 2, "Should have at least 2 spans (middleware + handler)")

		// Find the handler span
		var handlerSpan *tracetest.SpanStub
		for i := range spans {
			if spans[i].Name == "process-request" {
				handlerSpan = &spans[i]
				break
			}
		}

		assert.NotNil(t, handlerSpan, "Handler span should exist")
		assert.Equal(t, "process-request", handlerSpan.Name)
		assert.Contains(t, handlerSpan.Attributes, attribute.String("handler", "test"))
	})

	t.Run("should demonstrate OTLP exporter initialization", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// This test verifies that the OTLP exporter can be initialized
		// Even if the collector is not running, initialization should succeed
		tp, cleanup, err := InitProvider(ctx, "test-service", "localhost:4317", logger)
		
		assert.NoError(t, err, "InitProvider should not error during initialization")
		assert.NotNil(t, tp, "TracerProvider should not be nil")
		assert.NotNil(t, cleanup, "Cleanup function should not be nil")

		// Create a span to demonstrate the workflow
		tracer := otel.Tracer("test-tracer")
		_, span := tracer.Start(ctx, "test-span-with-otlp")
		span.SetAttributes(
			attribute.String("test.workflow", "otlp-exporter"),
			attribute.Bool("test.enabled", true),
		)
		span.End()

		// Cleanup
		if cleanup != nil {
			err := cleanup(context.Background())
			assert.NoError(t, err, "Cleanup should not error")
		}
	})
}
