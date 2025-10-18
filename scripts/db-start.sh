#!/bin/bash
set -e

echo "Starting PostgreSQL database..."

# Check if data directory exists
if [ ! -d "./data" ]; then
    echo "Error: Database data directory './data' not found."
    echo "Please initialize the database first. See docs/database/postgresql.md for instructions."
    exit 1
fi

# Start PostgreSQL
pg_ctl -D "./data" start

echo "PostgreSQL database started successfully!"
