#!/bin/bash
set -e

# Database migration script
# Usage: ./migrate.sh [up|down|create|force] [args...]

COMMAND=${1:-up}
shift || true

# Check if migrate CLI is installed
if ! command -v migrate &> /dev/null; then
    echo "❌ migrate CLI not installed"
    echo "   Install from: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"
    echo ""
    echo "   Quick install:"
    echo "   curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz"
    echo "   sudo mv migrate /usr/local/bin/"
    exit 1
fi

# Check DATABASE_URL
if [ -z "$DATABASE_URL" ]; then
    echo "❌ DATABASE_URL environment variable not set"
    echo "   Example: export DATABASE_URL='postgres://user:pass@localhost:5432/myapp?sslmode=disable'"
    exit 1
fi

case "$COMMAND" in
    up)
        echo "Running database migrations..."
        migrate -path migrations -database "$DATABASE_URL" up
        echo "✅ Migrations applied successfully"
        ;;
    down)
        echo "Rolling back last migration..."
        migrate -path migrations -database "$DATABASE_URL" down 1
        echo "✅ Migration rolled back"
        ;;
    create)
        if [ -z "$1" ]; then
            echo "❌ Migration name required"
            echo "   Usage: ./migrate.sh create <migration_name>"
            exit 1
        fi
        echo "Creating migration: $1"
        migrate create -ext sql -dir migrations -seq "$1"
        echo "✅ Migration files created in ./migrations"
        ;;
    force)
        if [ -z "$1" ]; then
            echo "❌ Version required"
            echo "   Usage: ./migrate.sh force <version>"
            exit 1
        fi
        echo "Forcing migration version to: $1"
        migrate -path migrations -database "$DATABASE_URL" force "$1"
        echo "✅ Migration version forced to $1"
        ;;
    version)
        migrate -path migrations -database "$DATABASE_URL" version
        ;;
    *)
        echo "Usage: ./migrate.sh [up|down|create|force|version] [args...]"
        echo ""
        echo "Commands:"
        echo "  up              - Apply all pending migrations"
        echo "  down            - Rollback last migration"
        echo "  create <name>   - Create new migration files"
        echo "  force <version> - Force migration to specific version"
        echo "  version         - Show current migration version"
        exit 1
        ;;
esac
