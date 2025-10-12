# Generate Swagger documentation
$ErrorActionPreference = "Stop"

Write-Host "Generating Swagger documentation..." -ForegroundColor Cyan

# Check if swag is installed
$swagPath = Get-Command swag -ErrorAction SilentlyContinue
if (-not $swagPath) {
    Write-Host "swag not found. Installing..." -ForegroundColor Yellow
    go install github.com/swaggo/swag/cmd/swag@latest
}

# Generate swagger docs
swag init -g cmd/server/main.go --output docs --parseDependency --parseInternal

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Swagger documentation generated in ./docs" -ForegroundColor Green
    Write-Host "   View at: http://localhost:8080/swagger/index.html (when server is running)" -ForegroundColor Cyan
} else {
    Write-Host "❌ Failed to generate Swagger documentation!" -ForegroundColor Red
    exit 1
}
