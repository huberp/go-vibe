# Observability

go-vibe is built observability-first. Metrics, structured logs, and health checks are all wired up out of the box.

## Prometheus Metrics

The `/metrics` endpoint exposes Prometheus-compatible metrics. Scrape it with any Prometheus server — no extra configuration needed.

### Available Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `http_requests_total` | Counter | `method`, `path`, `status` | Total HTTP requests by method, path, and response status |
| `http_request_duration_seconds` | Histogram | `method`, `path` | Latency distribution — p50, p95, p99 |
| `users_total` | Gauge | — | Current count of registered users in the database |
| `go_*` | Various | — | Standard Go runtime metrics (GC, goroutines, memory) |
| `process_*` | Various | — | OS process metrics (CPU, file descriptors) |

### Example Metrics Output

```
# HELP http_requests_total Total number of HTTP requests processed
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/v1/users",status="200"} 1847
http_requests_total{method="POST",path="/v1/login",status="200"} 532
http_requests_total{method="POST",path="/v1/login",status="401"} 14
http_requests_total{method="GET",path="/health",status="200"} 8921

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/v1/users",le="0.005"} 1204
http_request_duration_seconds_bucket{method="GET",path="/v1/users",le="0.01"} 1687
http_request_duration_seconds_bucket{method="GET",path="/v1/users",le="0.025"} 1843
http_request_duration_seconds_bucket{method="GET",path="/v1/users",le="+Inf"} 1847
http_request_duration_seconds_sum{method="GET",path="/v1/users"} 3.21
http_request_duration_seconds_count{method="GET",path="/v1/users"} 1847

# HELP users_total Total number of registered users
# TYPE users_total gauge
users_total 127
```

### Prometheus Scrape Config

Add this to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: go-vibe
    static_configs:
      - targets: ["localhost:8080"]
    metrics_path: /metrics
    scrape_interval: 15s
```

For Kubernetes, annotate the Pod or use a `ServiceMonitor` (Prometheus Operator):

```yaml
# Pod annotations approach
annotations:
  prometheus.io/scrape: "true"
  prometheus.io/path: "/metrics"
  prometheus.io/port: "8080"
```

```yaml
# ServiceMonitor (Prometheus Operator)
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: go-vibe
spec:
  selector:
    matchLabels:
      app: go-vibe
  endpoints:
    - port: http
      path: /metrics
      interval: 15s
```

## Structured Logging

go-vibe uses [Uber Zap](https://pkg.go.dev/go.uber.org/zap) for zero-allocation structured logging. All logs are emitted as JSON, making them easy to ingest with Loki, Datadog, CloudWatch, or Splunk.

### Log Fields

Every HTTP request produces a structured log entry with these fields:

| Field | Type | Example | Description |
|-------|------|---------|-------------|
| `level` | string | `"info"` | Log level |
| `ts` | float64 | `1719878400.123` | Unix timestamp |
| `msg` | string | `"HTTP request"` | Log message |
| `method` | string | `"GET"` | HTTP method |
| `path` | string | `"/v1/users"` | Request path |
| `status` | int | `200` | Response status code |
| `duration` | string | `"1.23ms"` | Request duration |
| `client_ip` | string | `"10.0.0.1"` | Client IP address |
| `request_id` | string | `"abc123"` | Unique request identifier |

### Example Log Output

```json
{
  "level": "info",
  "ts": 1719878400.123,
  "msg": "HTTP request",
  "method": "GET",
  "path": "/v1/users",
  "status": 200,
  "duration": "1.23ms",
  "client_ip": "10.0.0.1",
  "request_id": "f4a2b1c3"
}
```

Error logs include additional context:

```json
{
  "level": "error",
  "ts": 1719878401.456,
  "msg": "database query failed",
  "error": "connection refused",
  "user_id": 42,
  "operation": "FindByID"
}
```

### Logger Initialization

```go
// pkg/logger/logger.go
func NewLogger(env string) (*zap.Logger, error) {
    if env == "production" {
        return zap.NewProduction()
    }
    return zap.NewDevelopment()
}
```

In production, set `APP_ENV=production` to get JSON logs. In development, you get human-readable colored output.

::: tip Log Levels
- **`debug`** — verbose tracing, not enabled in production
- **`info`** — normal operations, every HTTP request
- **`warn`** — unexpected conditions that don't require immediate action
- **`error`** — failures that need attention
- **`fatal`** — unrecoverable startup failures (calls `os.Exit(1)`)
:::

::: danger Never Log Sensitive Data
Passwords, JWT tokens, and PII must never appear in logs. Zap's structured fields make it easy to audit — review every `zap.String` and `zap.Any` call before merging.
:::

## Health Check

The `/health` endpoint provides a lightweight liveness signal for load balancers and Kubernetes probes:

```bash
curl -s http://localhost:8080/health
```

```json
{
  "status": "healthy"
}
```

**HTTP status**: `200 OK` when healthy. Returns `503 Service Unavailable` if the database connection is lost (planned enhancement).

## OpenTelemetry Tracing

::: info Planned Feature
Distributed tracing with OpenTelemetry is planned. The architecture is designed to support it — context propagation is already in place throughout the request chain. Contributions welcome!
:::

When implemented, traces will be exported to:
- **Jaeger** (local dev) — `http://localhost:16686`
- **OTLP** endpoint — for production (Grafana Tempo, Honeycomb, Datadog)

## Grafana Dashboard

A Grafana dashboard can be built using the Prometheus metrics exposed at `/metrics`. The following panels are recommended:

- **Request rate** — requests per second by endpoint
- **Error rate** — 4xx and 5xx percentage
- **Latency** — p50, p95, p99 from the histogram
- **User growth** — `users_total` over time

### Useful PromQL Queries

```promql
# Request rate (last 5 minutes)
rate(http_requests_total[5m])

# Error rate
rate(http_requests_total{status=~"5.."}[5m])
  / rate(http_requests_total[5m])

# p99 latency
histogram_quantile(0.99,
  rate(http_request_duration_seconds_bucket[5m])
)

# Requests per endpoint
sum by (path) (rate(http_requests_total[5m]))
```
