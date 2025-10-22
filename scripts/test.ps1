# Test script for Windows
$ErrorActionPreference = "Stop"

Write-Host "Running tests..." -ForegroundColor Cyan
go test ./... -v

if ($LASTEXITCODE -eq 0) {
    Write-Host "[OK] All tests passed!" -ForegroundColor Green
} else {
    Write-Host "[FAILED] Tests failed!" -ForegroundColor Red
    exit 1
}
