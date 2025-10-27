#!/bin/bash

# Run smoke test
# Usage: ./run-smoke-test.sh [BASE_URL]

set -e

BASE_URL=${1:-http://localhost:8080}

echo "Running smoke test against: $BASE_URL"
echo "================================================"

if ! command -v k6 &> /dev/null; then
    echo "Error: k6 is not installed."
    echo "Please install k6 from: https://k6.io/docs/getting-started/installation/"
    exit 1
fi

k6 run --env BASE_URL="$BASE_URL" scripts/smoke-test.js

echo ""
echo "Smoke test completed!"
