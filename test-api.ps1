# API Testing Script for User Management Microservice
# This script demonstrates all API endpoints

$BaseUrl = "http://localhost:8080"

Write-Host "================================"
Write-Host "User Management API Test Script"
Write-Host "================================"
Write-Host ""

# Test health endpoint
Write-Host "1. Testing Health Endpoint..."
$response = Invoke-RestMethod -Uri "$BaseUrl/health" -Method Get
$response | ConvertTo-Json
Write-Host ""

# Create a regular user
Write-Host "2. Creating a regular user..."
$userBody = @{
    name = "John Doe"
    email = "john@example.com"
    password = "password123"
    role = "user"
} | ConvertTo-Json

$userResponse = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Post -Body $userBody -ContentType "application/json"
$userResponse | ConvertTo-Json
$userId = $userResponse.id
Write-Host ""

# Create an admin user
Write-Host "3. Creating an admin user..."
$adminBody = @{
    name = "Admin User"
    email = "admin@example.com"
    password = "admin123"
    role = "admin"
} | ConvertTo-Json

$adminResponse = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Post -Body $adminBody -ContentType "application/json"
$adminResponse | ConvertTo-Json
$adminId = $adminResponse.id
Write-Host ""

# Login as admin
Write-Host "4. Logging in as admin..."
$adminLoginBody = @{
    email = "admin@example.com"
    password = "admin123"
} | ConvertTo-Json

$adminLoginResponse = Invoke-RestMethod -Uri "$BaseUrl/v1/login" -Method Post -Body $adminLoginBody -ContentType "application/json"
$adminLoginResponse | ConvertTo-Json
$adminToken = $adminLoginResponse.token
Write-Host ""

# Login as regular user
Write-Host "5. Logging in as regular user..."
$userLoginBody = @{
    email = "john@example.com"
    password = "password123"
} | ConvertTo-Json

$userLoginResponse = Invoke-RestMethod -Uri "$BaseUrl/v1/login" -Method Post -Body $userLoginBody -ContentType "application/json"
$userLoginResponse | ConvertTo-Json
$userToken = $userLoginResponse.token
Write-Host ""

# Get all users (admin only)
Write-Host "6. Getting all users (admin token)..."
$headers = @{
    Authorization = "Bearer $adminToken"
}
$allUsers = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Get -Headers $headers
$allUsers | ConvertTo-Json
Write-Host ""

# Try to get all users as regular user (should fail)
Write-Host "7. Trying to get all users as regular user (should fail)..."
$userHeaders = @{
    Authorization = "Bearer $userToken"
}
try {
    $userAllUsers = Invoke-RestMethod -Uri "$BaseUrl/v1/users" -Method Get -Headers $userHeaders
    $userAllUsers | ConvertTo-Json
} catch {
    Write-Host "Error: $($_.Exception.Message)"
}
Write-Host ""

# Get user by ID
Write-Host "8. Getting user by ID..."
$userById = Invoke-RestMethod -Uri "$BaseUrl/v1/users/$userId" -Method Get -Headers $userHeaders
$userById | ConvertTo-Json
Write-Host ""

# Update user
Write-Host "9. Updating user..."
$updateBody = @{
    name = "John Updated"
    email = "john.updated@example.com"
} | ConvertTo-Json

$updatedUser = Invoke-RestMethod -Uri "$BaseUrl/v1/users/$userId" -Method Put -Body $updateBody -Headers $userHeaders -ContentType "application/json"
$updatedUser | ConvertTo-Json
Write-Host ""

# Delete user (admin only)
Write-Host "10. Deleting user as admin..."
try {
    Invoke-RestMethod -Uri "$BaseUrl/v1/users/$userId" -Method Delete -Headers $headers
    Write-Host "User deleted successfully"
} catch {
    Write-Host "Error: $($_.Exception.Message)"
}
Write-Host ""

# Check Prometheus metrics
Write-Host "11. Checking Prometheus metrics..."
$metrics = Invoke-WebRequest -Uri "$BaseUrl/metrics" -UseBasicParsing
$metricsLines = $metrics.Content -split "`n" | Select-String -Pattern "http_requests_total" | Select-Object -First 5
$metricsLines | ForEach-Object { Write-Host $_ }
Write-Host ""

Write-Host "================================"
Write-Host "API Tests Completed!"
Write-Host "================================"
