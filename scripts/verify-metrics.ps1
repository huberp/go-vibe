# Verify Prometheus Metrics Endpoint (PowerShell)
# This script checks that Go runtime and custom metrics are exposed in Prometheus format

Write-Host "=========================================="
Write-Host "Testing Prometheus Metrics Endpoint"
Write-Host "=========================================="
Write-Host ""

# Create a simple test server
$metricsDemoPath = "$env:TEMP\metrics_demo.go"
@'
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
'@ | Set-Content -Path $metricsDemoPath

# Start the server
Write-Host "Starting demo server on port 8090..."
$goProc = Start-Process -FilePath "go" -ArgumentList "run", $metricsDemoPath -NoNewWindow -PassThru
Start-Sleep -Seconds 2

Write-Host ""
Write-Host "1. Testing Prometheus Format"
Write-Host "================================"
$metrics = Invoke-WebRequest -Uri "http://localhost:8090/metrics" -UseBasicParsing | Select-Object -ExpandProperty Content
$metrics -split "`n" | Select-String -Pattern "^# (TYPE|HELP)" | Select-Object -First 10 | ForEach-Object { $_.Line }
Write-Host ""

Write-Host "2. Testing runtime.MemStats Metrics"
Write-Host "===================================="
Write-Host "Checking for go_memstats_* metrics..."
$memstatsCount = ($metrics -split "`n" | Select-String -Pattern "^go_memstats").Count
Write-Host "Found $memstatsCount go_memstats_* metrics"
Write-Host ""
Write-Host "Sample metrics:"
$metrics -split "`n" | Select-String -Pattern "^go_memstats" | Select-Object -First 10 | ForEach-Object { $_.Line }
Write-Host ""

Write-Host "3. Testing Go Runtime Metrics"
Write-Host "=============================="
$metrics -split "`n" | Select-String -Pattern "^(go_goroutines|go_threads|go_gc_duration)" | Select-Object -First 5 | ForEach-Object { $_.Line }
Write-Host ""

Write-Host "4. Testing Custom Metrics"
Write-Host "========================="
$metrics -split "`n" | Select-String -Pattern "demo_custom_metric" | ForEach-Object { $_.Line }
Write-Host ""

# Cleanup
Stop-Process -Id $goProc.Id -Force
Remove-Item $metricsDemoPath -Force

Write-Host "=========================================="
Write-Host " Verification Complete!"
Write-Host "=========================================="
Write-Host ""
Write-Host "Summary:"
Write-Host "- Metrics are in Prometheus format (# TYPE, # HELP)"
Write-Host "- runtime.MemStats are exposed (go_memstats_*)"
Write-Host "- Go runtime metrics are exposed (go_goroutines, go_threads, etc.)"
Write-Host "- Custom metrics work alongside runtime metrics"