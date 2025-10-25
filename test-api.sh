#!/bin/bash

# API Testing Script for User Management Microservice
# This script demonstrates all API endpoints

BASE_URL="http://localhost:8080"

echo "================================"
echo "User Management API Test Script"
echo "================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Error counter
ERROR_COUNT=0

# Function to check HTTP response and handle errors
check_response() {
    local response="$1"
    local description="$2"
    
    if [ -z "$response" ]; then
        echo -e "${RED}ERROR: Empty response for $description${NC}"
        ((ERROR_COUNT++))
        return 1
    fi
    
    # Check if response contains error field
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        echo -e "${RED}ERROR in response for $description:${NC}"
        echo "$response" | jq .
        ((ERROR_COUNT++))
        return 1
    fi
    
    return 0
}

# Test health endpoint
echo "1. Testing Health Endpoint..."
HEALTH_RESPONSE=$(curl -s "$BASE_URL/health")
if check_response "$HEALTH_RESPONSE" "health check"; then
    echo "$HEALTH_RESPONSE" | jq .
else
    echo -e "${YELLOW}Warning: Health check failed, but continuing...${NC}"
fi
echo ""

# Test health subresources
echo "2. Testing Health Startup Probe..."
STARTUP_RESPONSE=$(curl -s "$BASE_URL/health/startup")
if check_response "$STARTUP_RESPONSE" "startup probe"; then
    echo "$STARTUP_RESPONSE" | jq .
fi
echo ""

echo "3. Testing Health Liveness Probe..."
LIVENESS_RESPONSE=$(curl -s "$BASE_URL/health/liveness")
if check_response "$LIVENESS_RESPONSE" "liveness probe"; then
    echo "$LIVENESS_RESPONSE" | jq .
fi
echo ""

echo "4. Testing Health Readiness Probe..."
READINESS_RESPONSE=$(curl -s "$BASE_URL/health/readiness")
if check_response "$READINESS_RESPONSE" "readiness probe"; then
    echo "$READINESS_RESPONSE" | jq .
fi
echo ""

# Test info endpoint
echo "5. Testing Info Endpoint..."
INFO_RESPONSE=$(curl -s "$BASE_URL/info")
if check_response "$INFO_RESPONSE" "info endpoint"; then
    echo "$INFO_RESPONSE" | jq .
fi
echo ""

# Create a regular user
echo "6. Creating a regular user..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "role": "user"
  }')
if check_response "$USER_RESPONSE" "create user"; then
    echo "$USER_RESPONSE" | jq .
    USER_ID=$(echo "$USER_RESPONSE" | jq -r '.id')
fi
echo ""

# Create an admin user
echo "7. Creating an admin user..."
ADMIN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@example.com",
    "password": "admin123",
    "role": "admin"
  }')
if check_response "$ADMIN_RESPONSE" "create admin"; then
    echo "$ADMIN_RESPONSE" | jq .
    ADMIN_ID=$(echo "$ADMIN_RESPONSE" | jq -r '.id')
fi
echo ""

# Login as admin
echo "8. Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }')
if check_response "$LOGIN_RESPONSE" "admin login"; then
    echo "$LOGIN_RESPONSE" | jq .
    ADMIN_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
fi
echo ""

# Login as regular user
echo "9. Logging in as regular user..."
USER_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }')
if check_response "$USER_LOGIN_RESPONSE" "user login"; then
    echo "$USER_LOGIN_RESPONSE" | jq .
    USER_TOKEN=$(echo "$USER_LOGIN_RESPONSE" | jq -r '.token')
fi
echo ""

# Get all users (admin only)
echo "10. Getting all users (admin token)..."
ALL_USERS_RESPONSE=$(curl -s -X GET "$BASE_URL/v1/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
if check_response "$ALL_USERS_RESPONSE" "get all users"; then
    echo "$ALL_USERS_RESPONSE" | jq .
fi
echo ""

# Try to get all users as regular user (should fail)
echo "11. Trying to get all users as regular user (should fail)..."
USER_ALL_RESPONSE=$(curl -s -X GET "$BASE_URL/v1/users" \
  -H "Authorization: Bearer $USER_TOKEN")
# This should fail, so we don't check it
echo "$USER_ALL_RESPONSE" | jq .
echo ""

# Get user by ID
echo "12. Getting user by ID..."
USER_BY_ID_RESPONSE=$(curl -s -X GET "$BASE_URL/v1/users/$USER_ID" \
  -H "Authorization: Bearer $USER_TOKEN")
if check_response "$USER_BY_ID_RESPONSE" "get user by ID"; then
    echo "$USER_BY_ID_RESPONSE" | jq .
fi
echo ""

# Update user
echo "13. Updating user..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/v1/users/$USER_ID" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "email": "john.updated@example.com"
  }')
if check_response "$UPDATE_RESPONSE" "update user"; then
    echo "$UPDATE_RESPONSE" | jq .
fi
echo ""

# Delete user (admin only)
echo "14. Deleting user as admin..."
DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/v1/users/$USER_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -w "\nHTTP Status: %{http_code}\n")
echo "$DELETE_RESPONSE"
echo ""

# Check Prometheus metrics
echo "15. Checking Prometheus metrics (http_requests_total)..."
curl -s "$BASE_URL/metrics" | grep "http_requests_total" | head -5
echo ""

echo "16. Checking Prometheus metrics (http_request_duration_seconds)..."
curl -s "$BASE_URL/metrics" | grep "http_request_duration_seconds" | head -5
echo ""

echo "17. Checking Prometheus metrics (users_total)..."
curl -s "$BASE_URL/metrics" | grep "users_total"
echo ""

echo "================================"
echo "API Tests Completed!"
if [ $ERROR_COUNT -eq 0 ]; then
    echo -e "${GREEN}All tests passed successfully!${NC}"
else
    echo -e "${RED}Tests completed with $ERROR_COUNT errors${NC}"
fi
echo "================================"
