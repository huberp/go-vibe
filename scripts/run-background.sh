#!/bin/bash
set -e

PID_FILE="server.pid"

# Check if server is already running
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null 2>&1; then
        echo "⚠️  Server is already running with PID $PID"
        exit 1
    fi
fi

# Build first
echo "Building application..."
go build -o server ./cmd/server

# Start server in background
echo "Starting server in background..."
nohup ./server > server.log 2>&1 &
SERVER_PID=$!

# Save PID to file
echo $SERVER_PID > "$PID_FILE"
echo "📋 Captured PID: $SERVER_PID"
echo "📋 PID file content: $(cat "$PID_FILE")"

# Wait a moment and check if server is running
sleep 2
if ps -p $SERVER_PID > /dev/null 2>&1; then
    echo "✅ Server started successfully with PID $SERVER_PID"
    echo "📝 Logs: tail -f server.log"
else
    echo "❌ Server failed to start. Check server.log for details."
    rm -f "$PID_FILE"
    exit 1
fi
