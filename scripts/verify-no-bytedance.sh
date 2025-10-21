#!/bin/bash

# Bytedance Dependency Verification Script
# This script checks for ByteDance dependencies in the project

set -e

echo "=========================================="
echo "ByteDance Dependency Verification"
echo "=========================================="
echo ""

echo "Checking for ByteDance libraries in go.mod and go.sum..."
echo ""

BYTEDANCE_DEPS=$(go list -m all | grep -i bytedance || true)

if [ -z "$BYTEDANCE_DEPS" ]; then
    echo "✅ SUCCESS: No ByteDance dependencies found!"
    echo ""
    echo "The project is free of ByteDance/TikTok libraries."
    exit 0
else
    echo "⚠️  WARNING: ByteDance dependencies detected:"
    echo ""
    echo "$BYTEDANCE_DEPS"
    echo ""
    echo "Dependencies breakdown:"
    echo "$BYTEDANCE_DEPS" | wc -l | xargs echo "  - Total ByteDance packages:"
    echo ""
    
    # Check dependency chain
    echo "Dependency chain analysis:"
    echo "----------------------------------------"
    for dep in $(echo "$BYTEDANCE_DEPS" | awk '{print $1}'); do
        echo ""
        echo "📦 $dep"
        echo "   Used by:"
        go mod why $dep | tail -n +2 | head -n 5 | sed 's/^/   → /'
    done
    echo ""
    echo "=========================================="
    echo ""
    echo "To remove ByteDance dependencies, see:"
    echo "  - BYTEDANCE_SUBSTITUTION_SUMMARY.md"
    echo "  - BYTEDANCE_ANALYSIS.md"
    echo ""
    exit 1
fi
