#!/bin/bash
set -e

echo "Building application..."
go build -tags=go_json -o server ./cmd/server
echo "✅ Build complete! Binary: ./server"
