# Database Migrations Guide

For db schema migration ``golang-migrate`` is used.

## Overview

Database migrations provide a version-controlled way to manage schema changes:
- Track all database changes in version control
- Apply changes consistently across environments
- Rollback changes when needed
- Maintain schema history

## Migration Files

Migrations are stored in the `migrations/` directory with the naming convention:
```
{version}_{description}.up.sql    # Applied when migrating up
{version}_{description}.down.sql  # Applied when rolling back
```

Example:
```
migrations/
├── 000001_create_users_table.up.sql
└── 000001_create_users_table.down.sql
```

## How It Works

The application automatically runs migrations on startup using the `pkg/migration` package:
1. Attempts to run migrations from `migrations/` directory
2. Falls back to GORM AutoMigrate if migrations fail (backward compatibility)
3. Logs migration status

## Creating New Migrations

### Using the migration script (Recommended)

**Linux/macOS:**
```bash
# Create a new migration
./scripts/migrate.sh create add_user_profile

# This creates:
# migrations/000002_add_user_profile.up.sql
# migrations/000002_add_user_profile.down.sql
```

**Windows PowerShell:**
```powershell
# Create a new migration
.\scripts\migrate.ps1 create add_user_profile

# This creates:
# migrations/000002_add_user_profile.up.sql
# migrations/000002_add_user_profile.down.sql
```

### Using migrate CLI directly

```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir migrations -seq add_user_profile
```

### Manual Creation

Create two files following the naming convention:
```sql
-- 000002_add_user_profile.up.sql
ALTER TABLE users ADD COLUMN profile_picture VARCHAR(255);

-- 000002_add_user_profile.down.sql
ALTER TABLE users DROP COLUMN profile_picture;
```

## Running Migrations

### Automatic (Application Startup)

Migrations run automatically when the application starts. Check logs:
```
2024-01-01T00:00:00Z INFO Migrations applied successfully
```

### Manual Using Migration Scripts

**Linux/macOS:**
```bash
# Apply all pending migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down

# Force to specific version (if stuck)
./scripts/migrate.sh force 1
```

**Windows PowerShell:**
```powershell
# Apply all pending migrations
.\scripts\migrate.ps1 up

# Rollback last migration
.\scripts\migrate.ps1 down

# Force to specific version (if stuck)
.\scripts\migrate.ps1 force 1
```

### Manual Using migrate CLI

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/myapp?sslmode=disable"

# Apply migrations
migrate -path migrations -database "${DATABASE_URL}" up

# Rollback one migration
migrate -path migrations -database "${DATABASE_URL}" down 1

# Go to specific version
migrate -path migrations -database "${DATABASE_URL}" goto 2

# Force version (use with caution)
migrate -path migrations -database "${DATABASE_URL}" force 1
```

## Best Practices

### 1. Always Create Both Up and Down Migrations
```sql
-- ✅ Good: up.sql
CREATE TABLE posts (id SERIAL PRIMARY KEY, title VARCHAR(255));

-- ✅ Good: down.sql
DROP TABLE posts;
```

### 2. Make Migrations Idempotent
```sql
-- ✅ Good: won't fail if already exists
CREATE TABLE IF NOT EXISTS posts (...);
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts(title);

-- ❌ Bad: fails if already exists
CREATE TABLE posts (...);
CREATE INDEX idx_posts_title ON posts(title);
```

### 3. Use Transactions (when supported)
```sql
-- ✅ Good
BEGIN;
ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';
UPDATE users SET status = 'active' WHERE status IS NULL;
ALTER TABLE users ALTER COLUMN status SET NOT NULL;
COMMIT;
```

### 4. Test Both Directions

**Linux/macOS:**
```bash
# Test up
./scripts/migrate.sh up

# Verify changes
psql $DATABASE_URL -c "\d users"

# Test down
./scripts/migrate.sh down

# Verify rollback worked
psql $DATABASE_URL -c "\d users"
```

**Windows PowerShell:**
```powershell
# Test up
.\scripts\migrate.ps1 up

# Verify changes
psql $env:DATABASE_URL -c "\d users"

# Test down
.\scripts\migrate.ps1 down

# Verify rollback worked
psql $env:DATABASE_URL -c "\d users"
```

### 5. Never Modify Existing Migrations
- Once a migration is committed and deployed, create a new migration to make changes
- Modifying existing migrations can cause version conflicts

## Migration Strategy

### Development
- Run migrations automatically on app startup
- Create migrations as you develop features
- Test both up and down migrations locally

### Staging/Production
- Migrations run automatically on deployment
- Always test migrations in staging first
- Have a rollback plan ready
- Monitor migration logs during deployment

## Troubleshooting

### "Dirty database version"
This happens when a migration partially fails:

**Linux/macOS:**
```bash
# Check current version
migrate -path migrations -database "${DATABASE_URL}" version

# Force to a known good version
./scripts/migrate.sh force 1

# Then re-run migrations
./scripts/migrate.sh up
```

**Windows PowerShell:**
```powershell
# Check current version
migrate -path migrations -database "$env:DATABASE_URL" version

# Force to a known good version
.\scripts\migrate.ps1 force 1

# Then re-run migrations
.\scripts\migrate.ps1 up
```

### "File does not exist"
Ensure you're running from the project root where `migrations/` directory exists.

### Migration Already Applied
The tool tracks applied migrations in the `schema_migrations` table:
```sql
-- View migration history
SELECT * FROM schema_migrations;
```

### Rollback All Migrations
```bash
# Careful! This drops all tables
migrate -path migrations -database "${DATABASE_URL}" down -all
```

## CI/CD Integration

Migrations run automatically on application startup, so no separate CI/CD step is needed. However, you can add a validation step:

```yaml
# .github/workflows/test.yml
- name: Validate migrations
  run: |
    # Start test database
    docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=test postgres:13
    sleep 5
    
    # Run migrations
    export DATABASE_URL="postgres://postgres:test@localhost:5432/postgres?sslmode=disable"
    ./scripts/migrate.sh up
    
    # Verify migrations can rollback
    ./scripts/migrate.sh down
    ./scripts/migrate.sh up
```

## Docker Considerations

The Dockerfile copies migrations into the container:
```dockerfile
# Migrations are in the working directory
COPY . .
# So migrations/ is available at runtime
```

Ensure migrations run before the app starts:
```bash
# In production, migrations run automatically in main.go
# No manual intervention needed
```

## Kubernetes Deployment

For Kubernetes, you have two options:

### Option 1: Init Container (Recommended for Production)
```yaml
# deployment.yaml
initContainers:
- name: migrate
  image: migrate/migrate
  args:
    - "-path=/migrations"
    - "-database=$(DATABASE_URL)"
    - "up"
  volumeMounts:
  - name: migrations
    mountPath: /migrations
```

### Option 2: Application Startup (Current Implementation)
The app runs migrations on startup, so no init container is needed. This is simpler but means the app must have database schema modification permissions.

## Further Reading

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Best Practices](https://www.postgresql.org/docs/current/ddl.html)
- [Database Migration Best Practices](https://www.prisma.io/dataguide/types/relational/migrations)
