#!/bin/bash
set -e

echo "Running tests with coverage..."

# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

echo "✅ Coverage report generated!"
echo "   Text report: coverage.out"
echo "   HTML report: coverage.html"
echo ""
echo "View coverage summary:"
go tool cover -func=coverage.out | tail -n 1
