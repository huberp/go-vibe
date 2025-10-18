#!/bin/bash
set -e

echo "Stopping PostgreSQL database..."

# Check if data directory exists
if [ ! -d "./data" ]; then
    echo "Error: Database data directory './data' not found."
    exit 1
fi

# Stop PostgreSQL
pg_ctl -D "./data" stop

echo "PostgreSQL database stopped successfully!"
