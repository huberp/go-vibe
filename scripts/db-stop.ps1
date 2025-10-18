# Database stop script for Windows
$ErrorActionPreference = "Stop"

Write-Host "Stopping PostgreSQL database..." -ForegroundColor Cyan

# Check if data directory exists
if (-not (Test-Path ".\data")) {
    Write-Host "Error: Database data directory '.\data' not found." -ForegroundColor Red
    exit 1
}

# Stop PostgreSQL
pg_ctl -D ".\data" stop

if ($LASTEXITCODE -eq 0) {
    Write-Host "PostgreSQL database stopped successfully!" -ForegroundColor Green
} else {
    Write-Host "Failed to stop PostgreSQL database!" -ForegroundColor Red
    exit 1
}
