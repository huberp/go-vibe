#!/bin/bash

PID_FILE="server.pid"

if [ ! -f "$PID_FILE" ]; then
    echo "⚠️  No PID file found. Server may not be running."
    exit 1
fi

PID=$(cat "$PID_FILE")
echo "📋 PID file content: $PID"

if ps -p $PID > /dev/null 2>&1; then
    echo "Stopping server (PID $PID)..."
    kill $PID
    
    # Wait for process to stop
    for i in {1..10}; do
        if ! ps -p $PID > /dev/null 2>&1; then
            break
        fi
        sleep 1
    done
    
    # Force kill if still running
    if ps -p $PID > /dev/null 2>&1; then
        echo "Force stopping server..."
        kill -9 $PID
    fi
    
    rm -f "$PID_FILE"
    echo "✅ Server stopped"
else
    echo "⚠️  Server with PID $PID is not running"
    rm -f "$PID_FILE"
fi
