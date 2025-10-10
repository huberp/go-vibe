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
NC='\033[0m' # No Color

# Test health endpoint
echo "1. Testing Health Endpoint..."
curl -s "$BASE_URL/health" | jq .
echo ""

# Create a regular user
echo "2. Creating a regular user..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "role": "user"
  }')
echo "$USER_RESPONSE" | jq .
USER_ID=$(echo "$USER_RESPONSE" | jq -r '.id')
echo ""

# Create an admin user
echo "3. Creating an admin user..."
ADMIN_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@example.com",
    "password": "admin123",
    "role": "admin"
  }')
echo "$ADMIN_RESPONSE" | jq .
ADMIN_ID=$(echo "$ADMIN_RESPONSE" | jq -r '.id')
echo ""

# Login as admin
echo "4. Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }')
echo "$LOGIN_RESPONSE" | jq .
ADMIN_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
echo ""

# Login as regular user
echo "5. Logging in as regular user..."
USER_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }')
echo "$USER_LOGIN_RESPONSE" | jq .
USER_TOKEN=$(echo "$USER_LOGIN_RESPONSE" | jq -r '.token')
echo ""

# Get all users (admin only)
echo "6. Getting all users (admin token)..."
curl -s -X GET "$BASE_URL/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
echo ""

# Try to get all users as regular user (should fail)
echo "7. Trying to get all users as regular user (should fail)..."
curl -s -X GET "$BASE_URL/users" \
  -H "Authorization: Bearer $USER_TOKEN" | jq .
echo ""

# Get user by ID
echo "8. Getting user by ID..."
curl -s -X GET "$BASE_URL/users/$USER_ID" \
  -H "Authorization: Bearer $USER_TOKEN" | jq .
echo ""

# Update user
echo "9. Updating user..."
curl -s -X PUT "$BASE_URL/users/$USER_ID" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "email": "john.updated@example.com"
  }' | jq .
echo ""

# Delete user (admin only)
echo "10. Deleting user as admin..."
curl -s -X DELETE "$BASE_URL/users/$USER_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -w "\nHTTP Status: %{http_code}\n"
echo ""

# Check Prometheus metrics
echo "11. Checking Prometheus metrics..."
curl -s "$BASE_URL/metrics" | grep "http_requests_total" | head -5
echo ""

echo "================================"
echo "API Tests Completed!"
echo "================================"
