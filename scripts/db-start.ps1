# Database start script for Windows
$ErrorActionPreference = "Stop"

Write-Host "Starting PostgreSQL database..." -ForegroundColor Cyan

# Check if data directory exists
if (-not (Test-Path ".\data")) {
    Write-Host "Error: Database data directory '.\data' not found." -ForegroundColor Red
    Write-Host "Please initialize the database first. See docs/database/postgresql.md for instructions." -ForegroundColor Yellow
    exit 1
}

# Start PostgreSQL
pg_ctl -D ".\data" start

if ($LASTEXITCODE -eq 0) {
    Write-Host "PostgreSQL database started successfully!" -ForegroundColor Green
} else {
    Write-Host "Failed to start PostgreSQL database!" -ForegroundColor Red
    exit 1
}
