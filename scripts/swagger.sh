#!/bin/bash
set -e

echo "Generating Swagger documentation..."

# Check if swag is installed
if ! command -v swag &> /dev/null && ! [ -x "$HOME/go/bin/swag" ]; then
    echo "swag not found. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Use swag from PATH or go/bin
SWAG_CMD="swag"
if ! command -v swag &> /dev/null && [ -x "$HOME/go/bin/swag" ]; then
    SWAG_CMD="$HOME/go/bin/swag"
fi

# Generate swagger docs
$SWAG_CMD init -g cmd/server/main.go --output docs --parseDependency --parseInternal --exclude examples

echo "✅ Swagger documentation generated in ./docs"
echo "   View at: http://localhost:8080/swagger/index.html (when server is running)"
