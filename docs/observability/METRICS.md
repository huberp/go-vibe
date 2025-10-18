# Metrics Documentation

This document describes the metrics exposed by the application at the `/metrics` endpoint.

## Overview

The application exposes metrics in **Prometheus format** at the `/metrics` endpoint. These metrics include:
1. Custom HTTP request metrics
2. Go runtime metrics (including runtime.MemStats)
3. Go garbage collection metrics

All metrics are automatically collected and exposed by the Prometheus Go client library.

## Metric Categories

### 1. HTTP Request Metrics

Custom metrics tracking HTTP request patterns:

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `http_requests_total` | Counter | Total number of HTTP requests | `method`, `path`, `status` |
| `http_request_duration_seconds` | Histogram | HTTP request duration distribution | `method`, `path` |

**Example:**
```
http_requests_total{method="GET",path="/users",status="200"} 42
http_request_duration_seconds_bucket{method="GET",path="/users",le="0.005"} 38
```

### 2. Go Runtime Memory Metrics (runtime.MemStats)

Metrics derived from Go's `runtime.MemStats` structure, providing detailed memory statistics:

#### Heap Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `go_memstats_alloc_bytes` | Gauge | Bytes of allocated heap objects currently in use |
| `go_memstats_heap_alloc_bytes` | Gauge | Same as alloc_bytes (heap bytes allocated and in use) |
| `go_memstats_heap_sys_bytes` | Gauge | Heap memory obtained from OS |
| `go_memstats_heap_idle_bytes` | Gauge | Heap bytes waiting to be used |
| `go_memstats_heap_inuse_bytes` | Gauge | Heap bytes that are in use |
| `go_memstats_heap_released_bytes` | Gauge | Heap bytes released to OS |
| `go_memstats_heap_objects` | Gauge | Number of allocated heap objects |

#### Memory Allocation Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `go_memstats_sys_bytes` | Gauge | Total bytes of memory obtained from OS |
| `go_memstats_alloc_bytes_total` | Counter | Total bytes allocated (even if freed) |
| `go_memstats_mallocs_total` | Counter | Total number of heap allocations |
| `go_memstats_frees_total` | Counter | Total number of heap frees |
| `go_memstats_lookups_total` | Counter | Total number of pointer lookups |

#### Garbage Collection Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `go_memstats_gc_sys_bytes` | Gauge | Bytes used for garbage collection metadata |
| `go_memstats_last_gc_time_seconds` | Gauge | Time of last GC (seconds since Unix epoch) |
| `go_memstats_next_gc_bytes` | Gauge | Target heap size for next GC |
| `go_memstats_gc_cpu_fraction` | Gauge | Fraction of CPU time used by GC |

#### Stack and System Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `go_memstats_stack_inuse_bytes` | Gauge | Bytes in stack spans currently in use |
| `go_memstats_stack_sys_bytes` | Gauge | Bytes of stack memory obtained from OS |
| `go_memstats_mspan_inuse_bytes` | Gauge | Bytes of mspan structures in use |
| `go_memstats_mspan_sys_bytes` | Gauge | Bytes of mspan structures obtained from OS |
| `go_memstats_mcache_inuse_bytes` | Gauge | Bytes of mcache structures in use |
| `go_memstats_mcache_sys_bytes` | Gauge | Bytes of mcache structures obtained from OS |
| `go_memstats_buck_hash_sys_bytes` | Gauge | Bytes used by profiling bucket hash table |
| `go_memstats_other_sys_bytes` | Gauge | Bytes used for other system allocations |

### 3. Go Runtime Metrics

General Go runtime metrics:

| Metric | Type | Description |
|--------|------|-------------|
| `go_goroutines` | Gauge | Number of goroutines currently running |
| `go_threads` | Gauge | Number of OS threads created |
| `go_info` | Gauge | Information about Go version |

### 4. Garbage Collection Performance

| Metric | Type | Description |
|--------|------|-------------|
| `go_gc_duration_seconds` | Summary | GC pause duration distribution (quantiles: 0, 0.25, 0.5, 0.75, 1.0) |

## Prometheus Format

All metrics follow the Prometheus exposition format:

```
# HELP metric_name Description of the metric
# TYPE metric_name metric_type
metric_name{label1="value1",label2="value2"} value timestamp

# Example:
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.234567e+06
```

### Metric Types

- **Counter**: Cumulative metric that only increases (e.g., total requests)
- **Gauge**: Metric that can go up and down (e.g., current memory usage)
- **Histogram**: Distribution of values in configurable buckets
- **Summary**: Similar to histogram with pre-calculated quantiles

## How It Works

The application uses the Prometheus Go client library which automatically:

1. **Registers default collectors** including the Go collector when using `promhttp.Handler()`
2. **Collects runtime.MemStats** by calling `runtime.ReadMemStats()` periodically
3. **Exposes metrics** in Prometheus format at the `/metrics` endpoint

### Code Implementation

The metrics endpoint is set up in `internal/routes/routes.go`:

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

// Metrics endpoint - automatically includes Go runtime metrics
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

Custom HTTP metrics are collected in `internal/middleware/metrics.go`:

```go
var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    // ...
)
```

The `promauto` package automatically registers metrics with the default Prometheus registry, which is then exposed by `promhttp.Handler()`.

## Usage Examples

### Querying Metrics

```bash
# Get all metrics
curl http://localhost:8080/metrics

# Filter for memory metrics only
curl http://localhost:8080/metrics | grep go_memstats

# Get current heap allocation
curl -s http://localhost:8080/metrics | grep "^go_memstats_heap_alloc_bytes"

# Check goroutine count
curl -s http://localhost:8080/metrics | grep "^go_goroutines"
```

### Prometheus Queries (PromQL)

Once scraped by Prometheus, you can query metrics:

```promql
# Current heap memory usage
go_memstats_heap_alloc_bytes

# Memory allocation rate (bytes/sec)
rate(go_memstats_alloc_bytes_total[5m])

# HTTP request rate
rate(http_requests_total[5m])

# 95th percentile request duration
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Number of goroutines over time
go_goroutines
```

### Grafana Dashboard Examples

Create visualizations for:
- **Memory usage over time**: `go_memstats_heap_alloc_bytes`
- **GC pause duration**: `go_gc_duration_seconds`
- **Request rate**: `rate(http_requests_total[5m])`
- **Error rate**: `rate(http_requests_total{status=~"5.."}[5m])`

## Testing Metrics

The application includes comprehensive tests to verify metrics exposure:

```bash
# Run metrics tests
go test ./internal/routes -v -run TestMetricsEndpoint

# Verify metrics script
./scripts/verify-metrics.sh
```

## Monitoring Best Practices

1. **Track memory trends**: Monitor `go_memstats_heap_alloc_bytes` for memory leaks
2. **Watch GC pressure**: Check `go_gc_duration_seconds` and `go_memstats_gc_cpu_fraction`
3. **Monitor goroutines**: Sudden spikes in `go_goroutines` may indicate goroutine leaks
4. **Alert on anomalies**: Set up alerts for unusual patterns in memory or GC metrics
5. **Correlate with requests**: Compare `http_requests_total` with memory metrics

## References

- [Prometheus Exposition Formats](https://prometheus.io/docs/instrumenting/exposition_formats/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Go runtime.MemStats](https://pkg.go.dev/runtime#MemStats)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
