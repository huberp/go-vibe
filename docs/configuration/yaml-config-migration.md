# YAML Configuration Migration Guide

This guide helps you migrate from environment-only configuration to the new YAML-based configuration system.

## Overview

The new configuration system supports:
- ✅ YAML configuration files with stage support
- ✅ File-per-stage approach (base.yaml + stage-specific overrides)
- ✅ Environment variable overrides (for secrets and runtime config)
- ✅ Full backward compatibility with environment-only configuration
- ✅ Easy stage selection via CLI flag or environment variable

## Migration Options

### Option 1: Keep Using Environment Variables (No Migration Needed)

**Nothing changes!** Your existing environment-based configuration continues to work exactly as before.

```bash
# This still works
export DATABASE_URL="postgres://user:pass@host:5432/db"
export JWT_SECRET="your-secret"
export SERVER_PORT="8080"
./server
```

### Option 2: Adopt YAML Configuration (Recommended for New Deployments)

Use YAML files for base configuration and environment variables for secrets.

#### Step 1: Choose Your Stage

The application includes three pre-configured stages:
- **development** - Local development (default)
- **staging** - Staging environment
- **production** - Production environment

#### Step 2: Review Stage Configuration

**Development (config/development.yaml):**
```yaml
database:
  url: "postgres://user:password@localhost:5432/myapp?sslmode=disable"
jwt:
  secret: "dev-secret-key"
```

**Staging (config/staging.yaml):**
```yaml
database:
  url: "${DATABASE_URL}"  # From environment variable
  max_open_conns: 50
jwt:
  secret: "${JWT_SECRET}"  # From environment variable
```

**Production (config/production.yaml):**
```yaml
database:
  url: "${DATABASE_URL}"  # From environment variable
  max_open_conns: 100
  max_idle_conns: 25
  conn_max_lifetime: 60
jwt:
  secret: "${JWT_SECRET}"  # From environment variable
server:
  port: "${SERVER_PORT:8080}"
```

#### Step 3: Set Required Environment Variables

Only set secrets and runtime-specific values:

```bash
# Staging
export APP_STAGE=staging
export DATABASE_URL="postgres://user:pass@staging-host:5432/db"
export JWT_SECRET="staging-secret-key"
./server

# Production
export APP_STAGE=production
export DATABASE_URL="postgres://user:pass@prod-host:5432/db"
export JWT_SECRET="production-secret-key"
./server
```

Or use the CLI flag:

```bash
export DATABASE_URL="postgres://user:pass@prod-host:5432/db"
export JWT_SECRET="production-secret-key"
./server --stage=production
```

### Option 3: Hybrid Approach (Environment Variables + Stage Selection)

Use stage configuration for structural settings, environment variables for overrides.

```bash
# Use staging config but override specific values
export APP_STAGE=staging
export DATABASE_URL="postgres://custom:pass@custom-host:5432/db"
export SERVER_PORT="9090"
./server
```

## Kubernetes/Helm Migration

### Before (Environment Variables Only)

```yaml
# Old deployment
env:
  - name: DATABASE_URL
    valueFrom:
      secretKeyRef:
        name: myapp-secrets
        key: database-url
  - name: JWT_SECRET
    valueFrom:
      secretKeyRef:
        name: myapp-secrets
        key: jwt-secret
  - name: SERVER_PORT
    value: "8080"
```

### After (YAML Configuration with Stage)

```bash
# Install with stage selection
helm install myapp ./helm/myapp \
  --namespace production \
  --set config.stage=production \
  --create-namespace
```

The Helm chart automatically sets `APP_STAGE` environment variable based on `config.stage`.

### Optional: Use ConfigMap for YAML Files

```bash
# Mount YAML config files from ConfigMap
helm install myapp ./helm/myapp \
  --namespace production \
  --set config.stage=production \
  --set config.useConfigMap=true \
  --create-namespace
```

This creates a ConfigMap with all stage configurations and mounts them to `/etc/myapp/config`.

## Configuration Loading Priority

The configuration system uses this priority order (highest to lowest):

1. **Environment Variables** (highest priority)
2. **Stage-specific YAML** (e.g., `config/production.yaml`)
3. **Base YAML** (e.g., `config/base.yaml`)
4. **Default values** (fallback if no config files exist)

### Example

```yaml
# config/base.yaml
server:
  port: "8080"
database:
  max_open_conns: 25

# config/production.yaml
database:
  max_open_conns: 100
```

```bash
# Environment variable
export DB_MAX_OPEN_CONNS=150
export APP_STAGE=production
./server
```

**Result:**
- `server.port` = "8080" (from base.yaml)
- `database.max_open_conns` = 150 (from environment variable - highest priority)

## Customizing Configuration

### Add Your Own Stage

Create a new stage file (e.g., `config/qa.yaml`):

```yaml
# config/qa.yaml
database:
  url: "${DATABASE_URL}"
  max_open_conns: 50

jwt:
  secret: "${JWT_SECRET}"
```

Use it:

```bash
export APP_STAGE=qa
export DATABASE_URL="postgres://user:pass@qa-host:5432/db"
export JWT_SECRET="qa-secret"
./server
```

### Modify Existing Stage

Edit the stage file in `config/` directory:

```yaml
# config/production.yaml
database:
  url: "${DATABASE_URL}"
  max_open_conns: 200  # Increased from 100
  max_idle_conns: 50   # Increased from 25
  conn_max_lifetime: 120  # Increased from 60
```

No code changes needed!

## Environment Variable Mapping

YAML configuration keys map to environment variables using underscores:

| YAML Path | Environment Variable |
|-----------|---------------------|
| `server.port` | `SERVER_PORT` |
| `database.url` | `DATABASE_URL` |
| `database.max_open_conns` | `DB_MAX_OPEN_CONNS` |
| `database.max_idle_conns` | `DB_MAX_IDLE_CONNS` |
| `database.conn_max_lifetime` | `DB_CONN_MAX_LIFETIME` |
| `jwt.secret` | `JWT_SECRET` |

## Best Practices

### ✅ DO:
- Store secrets in environment variables, not YAML files
- Use placeholders `${ENV_VAR}` in YAML for sensitive values
- Use different secrets for each stage
- Commit base.yaml and stage files to version control
- Use `--stage` flag or `APP_STAGE` env var for stage selection
- Override specific values with environment variables when needed

### ❌ DON'T:
- Don't commit production secrets to YAML files
- Don't hardcode passwords or tokens in YAML
- Don't use same JWT secret across stages
- Don't modify config files in production (use env vars for overrides)

## Troubleshooting

### Configuration not loading

**Check file paths:**
```bash
# Config files should be in one of these locations:
./config/base.yaml
../config/base.yaml
../../config/base.yaml
/etc/myapp/config/base.yaml
```

**Verify stage selection:**
```bash
# Check which stage is active
export APP_STAGE=production
./server
# Logs will show: "Starting server" stage="production"
```

### Environment variables not overriding

**Environment variables must match the mapping:**
```bash
# Correct
export SERVER_PORT="9090"

# Won't work (wrong name)
export PORT="9090"
```

### Stage not found

**Check stage file exists:**
```bash
ls -la config/
# Should show: base.yaml, development.yaml, staging.yaml, production.yaml
```

**Use correct stage name:**
```bash
# Correct
export APP_STAGE=production

# Won't work (typo)
export APP_STAGE=prod
```

## Rollback Plan

If you encounter issues, you can instantly roll back to environment-only configuration:

```bash
# Remove or rename config directory
mv config config.backup

# Set all values via environment variables
export DATABASE_URL="postgres://user:pass@host:5432/db"
export JWT_SECRET="your-secret"
export SERVER_PORT="8080"
export DB_MAX_OPEN_CONNS=25
export DB_MAX_IDLE_CONNS=10
export DB_CONN_MAX_LIFETIME=30

# Start server (uses defaults + env vars)
./server
```

The application will log: "No base config file found, using defaults"

## Support

- **Issue Tracker**: Report problems at GitHub Issues
- **Documentation**: See README.md Configuration section
- **Config Options**: Review `docs/configuration/yaml-config-options.md`

## Summary

The new YAML configuration system is:
- ✅ **Backward Compatible** - Existing deployments work unchanged
- ✅ **Flexible** - Use YAML files, environment variables, or both
- ✅ **Secure** - Secrets stay in environment variables
- ✅ **Maintainable** - Stage-specific files, easy to understand
- ✅ **Production-Ready** - Tested and validated

Choose the approach that works best for your deployment!
