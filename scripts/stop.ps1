# Stop server for Windows
$ErrorActionPreference = "Stop"

$PID_FILE = "server.pid"

if (-not (Test-Path $PID_FILE)) {
    Write-Host "⚠️  No PID file found. Server may not be running." -ForegroundColor Yellow
    exit 1
}

$local:PID_TO_STOP = Get-Content $PID_FILE
Write-Host "📋 PID file content: $PID_TO_STOP" -ForegroundColor Cyan
$process = Get-Process -Id $PID_TO_STOP -ErrorAction SilentlyContinue

if ($process) {
    Write-Host "Stopping server (PID $PID_TO_STOP)..." -ForegroundColor Cyan
    Stop-Process -Id $PID_TO_STOP -Force
    
    # Wait for process to stop
    for ($i = 1; $i -le 10; $i++) {
        $process = Get-Process -Id $PID_TO_STOP -ErrorAction SilentlyContinue
        if (-not $process) {
            break
        }
        Start-Sleep -Seconds 1
    }
    
    Remove-Item -Path $PID_FILE -ErrorAction SilentlyContinue
    Write-Host "✅ Server stopped" -ForegroundColor Green
} else {
    Write-Host "⚠️  Server with PID $PID_TO_STOP is not running" -ForegroundColor Yellow
    Remove-Item -Path $PID_FILE -ErrorAction SilentlyContinue
}
