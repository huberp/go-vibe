# Getting Started

This guide walks you through setting up go-vibe locally, exploring its API, and getting productive fast.

## Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Go | 1.21+ | [Download](https://go.dev/dl/) |
| PostgreSQL | 13+ | Or use Docker Compose |
| Docker | 24+ | For containerised workflow |
| `jq` | any | Optional — pretty-prints JSON |

## Clone and Setup

```bash
git clone https://github.com/huberp/go-vibe.git
cd go-vibe

# Download Go module dependencies
go mod download
go mod verify
```

## Configuration

go-vibe uses a **layered configuration** approach: a `config/base.yaml` file defines shared defaults, a stage-specific file (e.g. `config/development.yaml`) overrides them, and environment variables override everything.

```
config/base.yaml          ← shared defaults
config/development.yaml   ← overrides for local dev  (APP_STAGE=development)
config/staging.yaml       ← overrides for staging    (APP_STAGE=staging)
config/production.yaml    ← overrides for production (APP_STAGE=production)
```

Set the active stage with the `APP_STAGE` environment variable (defaults to `development`):

```bash
export APP_STAGE="production"
```

### Key Environment Variables

Environment variables always take precedence over YAML values and are the recommended way to inject secrets:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | ✅ | — | Full PostgreSQL DSN |
| `JWT_SECRET` | ✅ | — | Secret key for signing JWTs — use a long random string in production |
| `SERVER_PORT` | ❌ | `8080` | Port the HTTP server listens on |
| `APP_STAGE` | ❌ | `development` | Active config stage (`development`, `staging`, `production`) |

```bash
export APP_STAGE="development"
export DATABASE_URL="postgres://govibe:govibe@localhost:5432/govibe?sslmode=disable"
export JWT_SECRET="change-me-to-a-long-random-string-in-production"
```

::: tip Configuration Guide
See the [Configuration](/guide/configuration) page for the full list of YAML keys and environment variable bindings.
:::

::: warning Security
Never commit `JWT_SECRET` or `DATABASE_URL` to source control. Use Kubernetes Secrets or a secrets manager in production.
:::

## Running Locally

### Option A — Native Go

```bash
# Start PostgreSQL (skip if already running)
docker run -d \
  --name govibe-pg \
  -e POSTGRES_USER=govibe \
  -e POSTGRES_PASSWORD=govibe \
  -e POSTGRES_DB=govibe \
  -p 5432:5432 \
  postgres:15-alpine

# Set environment variables (see above), then:
go run ./cmd/server
```

The server starts on `http://localhost:8080` and auto-migrates the database schema on first run.

### Option B — Docker Compose

The included `docker-compose.yml` spins up both the application and PostgreSQL with a single command:

```bash
docker compose up --build
```

```yaml
# docker-compose.yml (excerpt)
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://govibe:govibe@db:5432/govibe?sslmode=disable
      JWT_SECRET: dev-secret-change-in-production
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: govibe
      POSTGRES_PASSWORD: govibe
      POSTGRES_DB: govibe
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U govibe"]
      interval: 5s
      timeout: 5s
      retries: 5
```

## Running the Tests

```bash
# Run all tests with verbose output
go test ./... -v

# Run with race detector (recommended before committing)
go test ./... -race

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

## Available Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/v1/users` | None | Register a new user |
| `POST` | `/v1/login` | None | Login, returns JWT token |
| `GET` | `/v1/users` | JWT (admin) | List all users |
| `GET` | `/v1/users/:id` | JWT (owner/admin) | Get user by ID |
| `PUT` | `/v1/users/:id` | JWT (owner/admin) | Update user |
| `DELETE` | `/v1/users/:id` | JWT (admin) | Delete user |
| `GET` | `/health` | None | Health check |
| `GET` | `/metrics` | None | Prometheus metrics |
| `GET` | `/swagger/*` | None | Swagger UI |

## Quick API Test

```bash
BASE="http://localhost:8080"

# Register a new admin user
curl -s -X POST $BASE/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123","role":"admin"}' | jq .

# Login
TOKEN=$(curl -s -X POST $BASE/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}' | jq -r .token)

echo "Token: $TOKEN"

# List users (admin JWT required)
curl -s $BASE/v1/users \
  -H "Authorization: Bearer $TOKEN" | jq .

# Health check
curl -s $BASE/health | jq .
```

## Next Steps

- Dive into the [Architecture overview](/guide/architecture) to understand how the layers fit together
- Explore [Features in depth](/guide/features) — JWT auth, RBAC, metrics, rate limiting
- Check the [API Reference](/guide/api) for full endpoint documentation
- Learn how to [deploy to Kubernetes](/guide/deployment)
