# API Testing Script for User Management Microservice
# This script demonstrates all API endpoints

$BaseUrl = "http://localhost:8080"

Write-Host "================================"
Write-Host "User Management API Test Script"
Write-Host "================================"
Write-Host ""

# Error counter
$ErrorCount = 0

# Function to check response and handle errors
function Test-Response {
    param(
        [Parameter(Mandatory=$true)]
        $Response,
        [Parameter(Mandatory=$true)]
        [string]$Description
    )
    
    if ($null -eq $Response) {
        Write-Host "ERROR: Empty response for $Description" -ForegroundColor Red
        $script:ErrorCount++
        return $false
    }
    
    # Check if response has an error property
    if ($Response.PSObject.Properties.Name -contains 'error') {
        Write-Host "ERROR in response for $Description" -ForegroundColor Red
        $Response | ConvertTo-Json
        $script:ErrorCount++
        return $false
    }
    
    return $true
}

# Test health endpoint
Write-Host "1. Testing Health Endpoint..."
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/health" -Method Get
    if (Test-Response -Response $response -Description "health check") {
        $response | ConvertTo-Json
    }
} catch {
    Write-Host "Warning: Health check failed - $($_.Exception.Message)" -ForegroundColor Yellow
}
Write-Host ""

# Test health subresources
Write-Host "2. Testing Health Startup Probe..."
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/health/startup" -Method Get
    if (Test-Response -Response $response -Description "startup probe") {
        $response | ConvertTo-Json
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

Write-Host "3. Testing Health Liveness Probe..."
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/health/liveness" -Method Get
    if (Test-Response -Response $response -Description "liveness probe") {
        $response | ConvertTo-Json
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

Write-Host "4. Testing Health Readiness Probe..."
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/health/readiness" -Method Get
    if (Test-Response -Response $response -Description "readiness probe") {
        $response | ConvertTo-Json
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Test info endpoint
Write-Host "5. Testing Info Endpoint..."
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/info" -Method Get
    if (Test-Response -Response $response -Description "info endpoint") {
        $response | ConvertTo-Json
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Create a regular user
Write-Host "6. Creating a regular user..."
$userBody = @{
    name = "John Doe"
    email = "john@example.com"
    password = "password123"
    role = "user"
} | ConvertTo-Json

try {
    $userResponse = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Post -Body $userBody -ContentType "application/json"
    if (Test-Response -Response $userResponse -Description "create user") {
        $userResponse | ConvertTo-Json
        $userId = $userResponse.id
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Create an admin user
Write-Host "7. Creating an admin user..."
$adminBody = @{
    name = "Admin User"
    email = "admin@example.com"
    password = "admin123"
    role = "admin"
} | ConvertTo-Json

try {
    $adminResponse = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Post -Body $adminBody -ContentType "application/json"
    if (Test-Response -Response $adminResponse -Description "create admin") {
        $adminResponse | ConvertTo-Json
        $adminId = $adminResponse.id
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Login as admin
Write-Host "8. Logging in as admin..."
$adminLoginBody = @{
    email = "admin@example.com"
    password = "admin123"
} | ConvertTo-Json

try {
    $adminLoginResponse = Invoke-RestMethod -Uri "$BaseUrl/v1/login" -Method Post -Body $adminLoginBody -ContentType "application/json"
    if (Test-Response -Response $adminLoginResponse -Description "admin login") {
        $adminLoginResponse | ConvertTo-Json
        $adminToken = $adminLoginResponse.token
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Login as regular user
Write-Host "9. Logging in as regular user..."
$userLoginBody = @{
    email = "john@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $userLoginResponse = Invoke-RestMethod -Uri "$BaseUrl/v1/login" -Method Post -Body $userLoginBody -ContentType "application/json"
    if (Test-Response -Response $userLoginResponse -Description "user login") {
        $userLoginResponse | ConvertTo-Json
        $userToken = $userLoginResponse.token
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Get all users (admin only)
Write-Host "10. Getting all users (admin token)..."
$headers = @{
    Authorization = "Bearer $adminToken"
}
try {
    $allUsers = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Get -Headers $headers
    if (Test-Response -Response $allUsers -Description "get all users") {
        $allUsers | ConvertTo-Json
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Try to get all users as regular user (should fail)
Write-Host "11. Trying to get all users as regular user (should fail)..."
$userHeaders = @{
    Authorization = "Bearer $userToken"
}
try {
    $userAllUsers = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Get -Headers $userHeaders
    $userAllUsers | ConvertTo-Json
} catch {
    # This is expected to fail
    Write-Host "Expected error: $($_.Exception.Message)" -ForegroundColor Yellow
}
Write-Host ""

# Get user by ID
Write-Host "12. Getting user by ID..."
try {
    $userById = Invoke-RestMethod -Uri "$BaseUrl/v1/users/$userId" -Method Get -Headers $userHeaders
    if (Test-Response -Response $userById -Description "get user by ID") {
        $userById | ConvertTo-Json
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Update user
Write-Host "13. Updating user..."
$updateBody = @{
    name = "John Updated"
    email = "john.updated@example.com"
} | ConvertTo-Json

try {
    $updatedUser = Invoke-RestMethod -Uri "$BaseUrl/v1/users/$userId" -Method Put -Body $updateBody -Headers $userHeaders -ContentType "application/json"
    if (Test-Response -Response $updatedUser -Description "update user") {
        $updatedUser | ConvertTo-Json
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Delete user (admin only)
Write-Host "14. Deleting user as admin..."
try {
    Invoke-RestMethod -Uri "$BaseUrl/v1/users/$userId" -Method Delete -Headers $headers
    Write-Host "User deleted successfully" -ForegroundColor Green
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

# Check Prometheus metrics (fetch once and filter multiple times)
Write-Host "15. Checking Prometheus metrics (http_requests_total)..."
try {
    $script:metrics = Invoke-WebRequest -Uri "$BaseUrl/metrics" -UseBasicParsing
    $metricsLines = $script:metrics.Content -split "`n" | Select-String -Pattern "http_requests_total" | Select-Object -First 5
    $metricsLines | ForEach-Object { Write-Host $_ }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

Write-Host "16. Checking Prometheus metrics (http_request_duration_seconds)..."
try {
    if ($script:metrics) {
        $metricsLines = $script:metrics.Content -split "`n" | Select-String -Pattern "http_request_duration_seconds" | Select-Object -First 5
        $metricsLines | ForEach-Object { Write-Host $_ }
    } else {
        Write-Host "Error: Metrics not available" -ForegroundColor Red
        $script:ErrorCount++
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

Write-Host "17. Checking Prometheus metrics (users_total)..."
try {
    if ($script:metrics) {
        $metricsLines = $script:metrics.Content -split "`n" | Select-String -Pattern "users_total"
        $metricsLines | ForEach-Object { Write-Host $_ }
    } else {
        Write-Host "Error: Metrics not available" -ForegroundColor Red
        $script:ErrorCount++
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $script:ErrorCount++
}
Write-Host ""

Write-Host "================================"
Write-Host "API Tests Completed!"
if ($ErrorCount -eq 0) {
    Write-Host "All tests passed successfully!" -ForegroundColor Green
} else {
    Write-Host "Tests completed with $ErrorCount errors" -ForegroundColor Red
}
Write-Host "================================"
