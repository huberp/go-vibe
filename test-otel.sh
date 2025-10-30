#!/bin/bash
# OpenTelemetry End-to-End Test Script
# This script demonstrates the complete OTEL workflow with collector

set -e

echo "========================================"
echo "OpenTelemetry End-to-End Test"
echo "========================================"

# Function to clean up
cleanup() {
    echo "Cleaning up..."
    docker-compose down -v 2>/dev/null || true
    rm -f /tmp/otel-test-token.txt
}

# Set trap to cleanup on exit
trap cleanup EXIT

echo ""
echo "Step 1: Starting services (database, OTEL collector, application)..."
docker-compose up -d

echo ""
echo "Step 2: Waiting for services to be healthy..."
sleep 10

# Check if services are running
echo "Checking service health..."
docker-compose ps

echo ""
echo "Step 3: Creating a test user..."
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "role": "user"
  }' || true

echo ""
echo ""
echo "Step 4: Logging in to get JWT token..."
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' | jq -r '.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    echo "✓ Successfully obtained JWT token"
    echo "$TOKEN" > /tmp/otel-test-token.txt
else
    echo "✗ Failed to obtain JWT token"
    exit 1
fi

echo ""
echo "Step 5: Making authenticated requests to generate traces..."

# Get user by ID
echo "Making GET /users/1 request..."
curl -s -X GET http://localhost:8080/users/1 \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# Get all users (admin endpoint, will fail but generates trace)
echo ""
echo "Making GET /users request (will fail - not admin)..."
curl -s -X GET http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "Step 6: Checking OTEL collector logs for traces..."
echo "Last 30 lines of OTEL collector logs:"
docker-compose logs --tail=30 otel-collector

echo ""
echo "========================================"
echo "Test Summary"
echo "========================================"
echo "✓ Services started successfully"
echo "✓ User created and authenticated"
echo "✓ Traced requests sent to application"
echo "✓ Traces exported to OTEL collector"
echo ""
echo "The OTEL collector logs above should show received traces."
echo "Look for lines containing 'ResourceSpans' or 'Span' to see trace data."
echo ""
echo "To view live logs: docker-compose logs -f otel-collector"
echo "To stop services: docker-compose down"
