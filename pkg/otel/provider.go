package otel

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitProvider initializes the OpenTelemetry TracerProvider with OTLP gRPC exporter.
// It returns the TracerProvider and a cleanup function that should be called on shutdown.
func InitProvider(ctx context.Context, serviceName, endpoint string, logger *zap.Logger) (*sdktrace.TracerProvider, func(context.Context) error, error) {
	if endpoint == "" {
		endpoint = "localhost:4317" // Default OTLP gRPC endpoint
	}

	logger.Info("Initializing OpenTelemetry provider",
		zap.String("service_name", serviceName),
		zap.String("otlp_endpoint", endpoint))

	// Create OTLP trace exporter with gRPC
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		otlptracegrpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Create resource with service name and other attributes
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
		),
		resource.WithSchemaURL(semconv.SchemaURL),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create TracerProvider with batch span processor
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set global TracerProvider
	otel.SetTracerProvider(tp)

	logger.Info("OpenTelemetry provider initialized successfully")

	// Return cleanup function
	cleanup := func(ctx context.Context) error {
		logger.Info("Shutting down OpenTelemetry provider")
		return tp.Shutdown(ctx)
	}

	return tp, cleanup, nil
}
