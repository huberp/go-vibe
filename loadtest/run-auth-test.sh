#!/bin/bash

# Run authentication load test
# Usage: ./run-auth-test.sh [BASE_URL]

set -e

BASE_URL=${1:-http://localhost:8080}

echo "Running authentication load test against: $BASE_URL"
echo "================================================"

if ! command -v k6 &> /dev/null; then
    echo "Error: k6 is not installed."
    echo "Please install k6 from: https://k6.io/docs/getting-started/installation/"
    exit 1
fi

k6 run --env BASE_URL="$BASE_URL" scripts/auth-load-test.js

echo ""
echo "Authentication load test completed!"
