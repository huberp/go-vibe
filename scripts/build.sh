#!/bin/bash
set -e

echo "Building application..."
go build -o server ./cmd/server
echo "✅ Build complete! Binary: ./server"
