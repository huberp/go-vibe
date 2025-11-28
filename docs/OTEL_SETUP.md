# OpenTelemetry Setup Documentation

This document describes the OpenTelemetry (OTEL) instrumentation setup in the go-vibe application.

## Overview

The application uses OpenTelemetry to provide distributed tracing capabilities. Traces are exported to an OpenTelemetry Collector via OTLP (OpenTelemetry Protocol) over gRPC.

## Architecture

```
Application (Gin + otelgin) → OTLP Exporter → OTEL Collector → Backend (Jaeger/Console)
```

### Components

1. **Application Instrumentation**
   - Uses `otelgin` middleware to automatically trace HTTP requests
   - Custom spans can be created using the OTEL SDK
   - Traces include W3C trace context headers

2. **OTLP Exporter**
   - Sends traces to OTEL Collector via gRPC (port 4317)
   - Configured in `pkg/otel/provider.go`
   - Uses batch span processor for efficiency

3. **OTEL Collector**
   - Receives traces via OTLP protocol
   - Processes traces (batching, filtering, etc.)
   - Exports to various backends (console, Jaeger, etc.)
   - Configuration in `otel-collector-config.yaml`

## Configuration

### Environment Variables

- `OBSERVABILITY_OTEL`: Enable/disable OTEL tracing (boolean, default: false)
- `OTEL_EXPORTER_OTLP_ENDPOINT`: OTLP endpoint (default: "localhost:4317")

### YAML Configuration

**config/base.yaml** or **config/development.yaml**:
```yaml
observability:
  otel: true
  otel_endpoint: "localhost:4317"
```

### OTEL Collector Configuration

The collector is configured in `otel-collector-config.yaml`:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  batch:
    timeout: 10s

exporters:
  logging:
    loglevel: debug

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging]
```

## Usage

### Running with Docker Compose

Start all services including OTEL collector:

```bash
docker-compose up
```

This starts:
- PostgreSQL database
- OTEL Collector (ports 4317/4318)
- Application (port 8080)

### Running Locally

1. Start OTEL Collector:
```bash
docker run -p 4317:4317 -p 4318:4318 \
  -v $(pwd)/otel-collector-config.yaml:/etc/otelcol/config.yaml \
  otel/opentelemetry-collector:latest
```

2. Set environment variables:
```bash
export OBSERVABILITY_OTEL=true
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
export DATABASE_URL="postgres://user:password@localhost:5432/myapp?sslmode=disable"
export JWT_SECRET="your-secret-key"
```

3. Run application:
```bash
go run ./cmd/server
```

### Testing OTEL Setup

Run the end-to-end test script:

```bash
./test-otel.sh
```

This script:
1. Starts all services via docker-compose
2. Creates a test user
3. Authenticates and obtains JWT
4. Makes traced HTTP requests
5. Shows OTEL collector logs with trace data

### Viewing Traces

**Console Output (Default)**

Traces are logged to the OTEL collector console:
```bash
docker-compose logs -f otel-collector
```

**Jaeger (Optional)**

To use Jaeger as a backend:

1. Update `otel-collector-config.yaml`:
```yaml
exporters:
  jaeger:
    endpoint: "jaeger:14250"
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, jaeger]
```

2. Add Jaeger to `docker-compose.yml`:
```yaml
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"  # Jaeger UI
      - "14250:14250"  # Jaeger gRPC
```

3. Access Jaeger UI at: http://localhost:16686

## Traced Endpoints

The following endpoints are automatically traced:

- `POST /login`
- `POST /users` (user creation)
- `GET /users/:id`
- `PUT /users/:id`
- `DELETE /users/:id`
- `GET /users` (list users)

**Excluded Endpoints** (not traced to reduce noise):
- `/health*` (all health check endpoints)
- `/metrics` (Prometheus metrics)
- `/info` (application info)

## Code Examples

### Creating Custom Spans

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

func myHandler(c *gin.Context) {
    // Get tracer
    tracer := otel.Tracer("my-service")
    
    // Create span
    ctx, span := tracer.Start(c.Request.Context(), "my-operation")
    defer span.End()
    
    // Add attributes
    span.SetAttributes(
        attribute.String("user.id", "123"),
        attribute.Int("items.count", 42),
    )
    
    // Use ctx in downstream calls
    result := doSomething(ctx)
    
    c.JSON(200, result)
}
```

### Propagating Context

Always use the request context to propagate trace context:

```go
func handler(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Pass context to repository/service calls
    user, err := repo.GetUser(ctx, id)
    
    // Context propagation ensures child spans are linked
}
```

## Testing

### Unit Tests

Run OTEL provider tests:
```bash
go test ./pkg/otel/... -v
```

### Integration Tests

The integration test (`pkg/otel/integration_test.go`) demonstrates:
- In-memory span export
- HTTP request tracing with otelgin
- OTLP exporter initialization

Run with:
```bash
go test ./pkg/otel/... -v -run TestOtelWorkflowWithExporter
```

## Troubleshooting

### Traces not appearing

1. Check OTEL is enabled:
   ```bash
   echo $OBSERVABILITY_OTEL  # should be "true"
   ```

2. Verify collector is running:
   ```bash
   docker-compose ps otel-collector
   ```

3. Check collector logs:
   ```bash
   docker-compose logs otel-collector
   ```

### Connection refused errors

- Ensure OTEL_EXPORTER_OTLP_ENDPOINT points to correct address
- In Docker: use service name (e.g., `otel-collector:4317`)
- Locally: use `localhost:4317`

### No spans in collector

- Verify requests are being made to traced endpoints
- Check that OBSERVABILITY_OTEL is set to true
- Ensure endpoints are not in the skip list (health, metrics, info)

## Performance Considerations

- **Batch Processing**: Spans are batched before export (10s timeout)
- **Sampling**: All traces are sampled by default; configure sampling in production
- **Resource Usage**: OTEL adds minimal overhead (~1-2ms per request)
- **Skip Paths**: Health checks and metrics excluded to reduce noise

## Security

- OTLP exporter uses insecure gRPC for local development
- In production, enable TLS for OTLP endpoint
- Don't expose sensitive data in span attributes
- Use network policies to restrict collector access

## References

- [OpenTelemetry Go SDK](https://pkg.go.dev/go.opentelemetry.io/otel)
- [otelgin Middleware](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin)
- [OTEL Collector](https://opentelemetry.io/docs/collector/)
- [W3C Trace Context](https://www.w3.org/TR/trace-context/)
