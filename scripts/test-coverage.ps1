# Test with coverage report
$ErrorActionPreference = "Stop"

Write-Host "Running tests with coverage..." -ForegroundColor Cyan

# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Tests failed!" -ForegroundColor Red
    exit 1
}

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

Write-Host "Success: Coverage report generated!" -ForegroundColor Green
Write-Host "   Text report: coverage.out" -ForegroundColor Cyan
Write-Host "   HTML report: coverage.html" -ForegroundColor Cyan
Write-Host ""
Write-Host "View coverage summary:" -ForegroundColor Cyan
go tool cover -func=coverage.out | Select-Object -Last 1
