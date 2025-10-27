#!/bin/bash

# Run stress test
# Usage: ./run-stress-test.sh [BASE_URL]

set -e

BASE_URL=${1:-http://localhost:8080}

echo "Running stress test against: $BASE_URL"
echo "================================================"
echo "WARNING: This test will generate significant load on your system."
echo "Make sure your application and database are ready for stress testing."
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Stress test cancelled."
    exit 1
fi

if ! command -v k6 &> /dev/null; then
    echo "Error: k6 is not installed."
    echo "Please install k6 from: https://k6.io/docs/getting-started/installation/"
    exit 1
fi

k6 run --env BASE_URL="$BASE_URL" scripts/stress-test.js

echo ""
echo "Stress test completed!"
