# Database migration script
# Usage: .\migrate.ps1 [up|down|create|force|version] [args...]
$ErrorActionPreference = "Stop"

$Command = if ($args.Count -gt 0) { $args[0] } else { "up" }
$RestArgs = if ($args.Count -gt 1) { $args[1..($args.Count-1)] } else { @() }

# Check if migrate CLI is installed
$migratePath = Get-Command migrate -ErrorAction SilentlyContinue
if (-not $migratePath) {
    Write-Host "❌ migrate CLI not installed" -ForegroundColor Red
    Write-Host "   Install from: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "   Windows install:" -ForegroundColor Yellow
    Write-Host "   Download from releases page and add to PATH" -ForegroundColor Yellow
    exit 1
}

# Check DATABASE_URL
if (-not $env:DATABASE_URL) {
    Write-Host "❌ DATABASE_URL environment variable not set" -ForegroundColor Red
    Write-Host "   Example: `$env:DATABASE_URL='postgres://user:pass@localhost:5432/myapp?sslmode=disable'" -ForegroundColor Yellow
    exit 1
}

switch ($Command) {
    "up" {
        Write-Host "Running database migrations..." -ForegroundColor Cyan
        migrate -path migrations -database $env:DATABASE_URL up
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Migrations applied successfully" -ForegroundColor Green
        } else {
            Write-Host "❌ Migration failed!" -ForegroundColor Red
            exit 1
        }
    }
    "down" {
        Write-Host "Rolling back last migration..." -ForegroundColor Cyan
        migrate -path migrations -database $env:DATABASE_URL down 1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Migration rolled back" -ForegroundColor Green
        } else {
            Write-Host "❌ Rollback failed!" -ForegroundColor Red
            exit 1
        }
    }
    "create" {
        if ($RestArgs.Count -eq 0) {
            Write-Host "❌ Migration name required" -ForegroundColor Red
            Write-Host "   Usage: .\migrate.ps1 create <migration_name>" -ForegroundColor Yellow
            exit 1
        }
        $migrationName = $RestArgs[0]
        Write-Host "Creating migration: $migrationName" -ForegroundColor Cyan
        migrate create -ext sql -dir migrations -seq $migrationName
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Migration files created in ./migrations" -ForegroundColor Green
        } else {
            Write-Host "❌ Failed to create migration!" -ForegroundColor Red
            exit 1
        }
    }
    "force" {
        if ($RestArgs.Count -eq 0) {
            Write-Host "❌ Version required" -ForegroundColor Red
            Write-Host "   Usage: .\migrate.ps1 force <version>" -ForegroundColor Yellow
            exit 1
        }
        $version = $RestArgs[0]
        Write-Host "Forcing migration version to: $version" -ForegroundColor Cyan
        migrate -path migrations -database $env:DATABASE_URL force $version
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Migration version forced to $version" -ForegroundColor Green
        } else {
            Write-Host "❌ Force failed!" -ForegroundColor Red
            exit 1
        }
    }
    "version" {
        migrate -path migrations -database $env:DATABASE_URL version
    }
    default {
        Write-Host "Usage: .\migrate.ps1 [up|down|create|force|version] [args...]" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "Commands:" -ForegroundColor Cyan
        Write-Host "  up              - Apply all pending migrations"
        Write-Host "  down            - Rollback last migration"
        Write-Host "  create <name>   - Create new migration files"
        Write-Host "  force <version> - Force migration to specific version"
        Write-Host "  version         - Show current migration version"
        exit 1
    }
}
