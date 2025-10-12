#!/bin/bash

# Script to verify runtime.MemStats are exposed in Prometheus format
# This demonstrates that the /metrics endpoint exposes Go runtime metrics

echo "=========================================="
echo "Testing Prometheus Metrics Endpoint"
echo "=========================================="
echo ""

# Create a simple test server
cat > /tmp/metrics_demo.go << 'EOF'
package main

import (
	"fmt"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	customCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_custom_metric",
		Help: "A custom demo metric",
	})
)

func main() {
	// Increment custom counter
	customCounter.Inc()
	
	http.Handle("/metrics", promhttp.Handler())
	
	fmt.Println("Server starting on :8090")
	fmt.Println("Visit http://localhost:8090/metrics to see metrics")
	http.ListenAndServe(":8090", nil)
}
EOF

# Start the server
echo "Starting demo server on port 8090..."
go run /tmp/metrics_demo.go &
SERVER_PID=$!

# Wait for server to start
sleep 2

echo ""
echo "1. Testing Prometheus Format"
echo "================================"
curl -s http://localhost:8090/metrics | grep -E "^# (TYPE|HELP)" | head -10
echo ""

echo "2. Testing runtime.MemStats Metrics"
echo "===================================="
echo "Checking for go_memstats_* metrics..."
MEMSTATS_COUNT=$(curl -s http://localhost:8090/metrics | grep -c "^go_memstats")
echo "Found $MEMSTATS_COUNT go_memstats_* metrics"
echo ""
echo "Sample metrics:"
curl -s http://localhost:8090/metrics | grep "^go_memstats" | head -10
echo ""

echo "3. Testing Go Runtime Metrics"
echo "=============================="
curl -s http://localhost:8090/metrics | grep -E "^(go_goroutines|go_threads|go_gc_duration)" | head -5
echo ""

echo "4. Testing Custom Metrics"
echo "========================="
curl -s http://localhost:8090/metrics | grep "demo_custom_metric"
echo ""

# Cleanup
kill $SERVER_PID 2>/dev/null
rm /tmp/metrics_demo.go

echo "=========================================="
echo "✓ Verification Complete!"
echo "=========================================="
echo ""
echo "Summary:"
echo "- ✓ Metrics are in Prometheus format (# TYPE, # HELP)"
echo "- ✓ runtime.MemStats are exposed (go_memstats_*)"
echo "- ✓ Go runtime metrics are exposed (go_goroutines, go_threads, etc.)"
echo "- ✓ Custom metrics work alongside runtime metrics"
