# Set PowerShell execution policy for the current process
# This allows running PowerShell scripts in this session without permanently changing system settings

Write-Host "Setting PowerShell execution policy for current process..." -ForegroundColor Cyan

Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass

if ($LASTEXITCODE -eq 0 -or $?) {
    Write-Host "Success: Execution policy set to Bypass for current process." -ForegroundColor Green
    Write-Host "You can now run PowerShell scripts in this session." -ForegroundColor White
} else {
    Write-Host "Warning: Failed to set execution policy." -ForegroundColor Yellow
}
