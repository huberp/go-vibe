# Bytedance Dependency Verification Script (PowerShell)
# This script checks for ByteDance dependencies in the project

$ErrorActionPreference = "Stop"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "ByteDance Dependency Verification" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Checking for ByteDance libraries in go.mod and go.sum..."
Write-Host ""

# Get all module dependencies
$allDeps = go list -m all 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Failed to list Go modules" -ForegroundColor Red
    exit 1
}

# Filter for ByteDance dependencies
$bytedanceDeps = $allDeps | Select-String -Pattern "bytedance" -SimpleMatch

if ($null -eq $bytedanceDeps -or $bytedanceDeps.Count -eq 0) {
    Write-Host "Success: No ByteDance dependencies found!" -ForegroundColor Green
    Write-Host ""
    Write-Host "The project is free of ByteDance/TikTok libraries."
    exit 0
} else {
    Write-Host "WARNING: ByteDance dependencies detected:" -ForegroundColor Yellow
    Write-Host ""
    $bytedanceDeps | ForEach-Object { Write-Host $_ -ForegroundColor Yellow }
    Write-Host ""
    Write-Host "Dependencies breakdown:"
    Write-Host "  - Total ByteDance packages: $($bytedanceDeps.Count)"
    Write-Host ""
    
    # Check dependency chain
    Write-Host "Dependency chain analysis:"
    Write-Host "----------------------------------------"
    
    $depNames = $bytedanceDeps | ForEach-Object { $_.ToString().Split()[0] }
    
    foreach ($dep in $depNames) {
        Write-Host ""
        Write-Host "Package: $dep" -ForegroundColor Cyan
        Write-Host "   Used by:"
        $whyOutput = go mod why $dep 2>&1 | Select-Object -Skip 1 | Select-Object -First 5
        $whyOutput | ForEach-Object { Write-Host "   -> $_" }
    }
    
    Write-Host ""
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "To remove ByteDance dependencies, see:"
    Write-Host "  - BYTEDANCE_SUBSTITUTION_SUMMARY.md"
    Write-Host "  - BYTEDANCE_ANALYSIS.md"
    Write-Host ""
    exit 1
}
