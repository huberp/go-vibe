# Run server in background for Windows
$ErrorActionPreference = "Stop"

$local:PID_FILE = "server.pid"

# Check if server is already running
if (Test-Path $PID_FILE) {
    $local:PID_OF_GO_PROCESS = Get-Content $PID_FILE
    $local:process = Get-Process -Id $local:PID_OF_GO_PROCESS -ErrorAction SilentlyContinue
    if ($process) {
        Write-Host "Server is already running with PID $local:PID_OF_GO_PROCESS" -ForegroundColor Yellow
        exit 1
    }
}

# Build first
Write-Host "Building application..." -ForegroundColor Cyan
go build -tags=go_json -o server.exe ./cmd/server

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

# Start server in background
Write-Host "Starting server in background..." -ForegroundColor Cyan
$job = Start-Process -FilePath ".\server.exe" -RedirectStandardOutput "server.log" -RedirectStandardError "server-error.log" -PassThru -WindowStyle Hidden

# Save PID to file
$job.Id | Out-File -FilePath $PID_FILE -NoNewline
Write-Host "Captured PID: $($job.Id)" -ForegroundColor Cyan
Write-Host "PID file content: $(Get-Content $PID_FILE)" -ForegroundColor Cyan

# Wait a moment and check if server is running
Start-Sleep -Seconds 2
$process = Get-Process -Id $job.Id -ErrorAction SilentlyContinue
Write-Host "Process: $($process) $($job.Id)" -ForegroundColor Cyan
if ($process) {
    Write-Host "Server started successfully with PID $($job.Id)" -ForegroundColor Green
    Write-Host "Logs: Get-Content server.log -Wait" -ForegroundColor Cyan
} else {
    Write-Host "Server failed to start. Check server.log and server-error.log for details." -ForegroundColor Red
    Remove-Item -Path $PID_FILE -ErrorAction SilentlyContinue
    exit 1
}
