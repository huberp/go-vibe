# Build script for Windows
$ErrorActionPreference = "Stop"

Write-Host "Building application..." -ForegroundColor Cyan
go build -o server.exe ./cmd/server

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Build complete! Binary: .\server.exe" -ForegroundColor Green
} else {
    Write-Host "❌ Build failed!" -ForegroundColor Red
    exit 1
}
