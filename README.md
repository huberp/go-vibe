# Production-Ready User Management Microservice

A production-ready microservice built with Go 1.25.2, Gin v1.11.0, following TDD principles and designed for Kubernetes deployment.

## Features

- ✅ RESTful API for user management (CRUD operations)
- ✅ JWT-based authentication and authorization
- ✅ Role-based access control (admin/user)
- ✅ PostgreSQL database with GORM
- ✅ **Database migrations with golang-migrate** (version-controlled schema changes)
- ✅ **Auto-generated OpenAPI/Swagger documentation** (accessible at /swagger)
- ✅ **Registrable health check system with scopes** (startup, liveness, readiness, base)
- ✅ YAML-based configuration with stage support (development, staging, production)
- ✅ Flexible configuration: YAML files + environment variable overrides
- ✅ **Configurable rate limiting** (via YAML/environment variables)
- ✅ Structured logging with Zap
- ✅ W3C trace context support for distributed tracing
- ✅ OpenTelemetry (OTEL) tracing integration
- ✅ Prometheus metrics (including user count)
- ✅ CORS middleware for cross-origin requests
- ✅ API versioning (/v1/...) for backward compatibility
- ✅ Enhanced bcrypt security (cost factor: 12)
- ✅ 100% test coverage for handlers and middleware
- ✅ **Dependency management with go mod tidy in CI**
- ✅ Multi-stage Docker build (Alpine-based)
- ✅ Helm chart for Kubernetes deployment
- ✅ GitHub Actions CI/CD pipelines
- ✅ **Shell scripts for common development tasks** (Linux/macOS and Windows PowerShell)

## Architecture

```
myapp/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── handlers/               # HTTP handlers (TDD)
│   ├── models/                 # GORM models
│   ├── routes/                 # Gin routes
│   ├── middleware/             # Auth, logging, metrics
│   └── repository/             # Database layer
├── pkg/
│   ├── config/                 # Configuration
│   ├── logger/                 # Zap logger
│   ├── migration/              # Database migrations
│   └── utils/                  # Utilities (JWT, hashing)
├── migrations/                 # SQL migration files
├── docs/                       # Generated Swagger docs
├── helm/                       # Helm chart
├── .github/workflows/          # CI/CD pipelines
├── scripts/                    # Development scripts (.sh/.ps1)
├── Dockerfile                  # Multi-stage build
├── go.mod
└── README.md
```

## Tech Stack

| Component | Version | Purpose |
|-----------|---------|---------|
| Go | 1.25.2 | Programming language |
| Gin | v1.11.0 | HTTP framework |
| GORM | v1.31.0 | ORM |
| PostgreSQL Driver | v1.6.0 | Database driver |
| **golang-migrate** | **v4.19.0** | **Database migrations** |
| JWT | v5.3.0 | Authentication |
| Viper | v1.21.0 | Configuration management |
| Zap | v1.27.0 | Structured logging |
| Testify | v1.11.1 | Testing framework |
| Prometheus | v1.23.2 | Metrics |
| OpenTelemetry | v1.37.0 | Distributed tracing |
| CORS | v1.7.0 | Cross-origin resource sharing |
| Rate Limiter | golang.org/x/time | Request rate limiting |
| **Swagger** | **v1.16.6** | **API documentation** |

## Prerequisites

- Go 1.25.2+
- PostgreSQL 13+
- Docker (for containerization)
- Kubernetes cluster (for deployment)
- Helm 3+ (for deployment)
- **Windows PowerShell users**: Run `.\scripts\set-execution-policy.ps1` before using PowerShell scripts

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/huberp/go-vibe.git
cd go-vibe
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set environment variables

Note: Please refer to ./docs/database/postgresql.md for installing postgre and setting up the User "myapp" and DB "myapp"

```bash
export DATABASE_URL="postgres://myapp:myapp@localhost:5432/myapp?sslmode=disable"
export JWT_SECRET="your-secret-key-change-in-production"
export SERVER_PORT="8080"
```

**Windows PowerShell users**: Before running PowerShell scripts, set the execution policy:
```powershell
.\scripts\set-execution-policy.ps1
```

Then set environment variables:
```powershell
$env:DATABASE_URL="postgres://myapp:myapp@localhost:5432/myapp?sslmode=disable"
$env:JWT_SECRET="your-secret-key-change-in-production"
$env:SERVER_PORT="8080"
```

### 4. Run the application

Using scripts:
```bash
# Linux/macOS: Generate swagger docs and run the server
./scripts/swagger.sh && go run ./cmd/server

# Windows PowerShell:
.\scripts\swagger.ps1; go run ./cmd/server
```

Or directly with Go:
```bash
go run ./cmd/server
```

The server will start on `http://localhost:8080`

### 5. Access API Documentation

Open your browser and navigate to:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health (aggregated), http://localhost:8080/health/liveness, http://localhost:8080/health/readiness
- **Metrics**: http://localhost:8080/metrics

## Development Commands

The project includes shell scripts in `scripts/` for common development tasks. See [scripts/README.md](scripts/README.md) for full documentation.

### Quick Reference

**Linux/macOS:**
```bash
# Database management
./scripts/db-start.sh          # Start PostgreSQL
./scripts/db-stop.sh           # Stop PostgreSQL

# Build the application
./scripts/build.sh

# Run tests
./scripts/test.sh

# Run tests with coverage
./scripts/test-coverage.sh

# Generate Swagger documentation
./scripts/swagger.sh

# Database migrations
./scripts/migrate.sh up              # Apply all migrations
./scripts/migrate.sh down            # Rollback last migration
./scripts/migrate.sh create xxx      # Create new migration

# Run server in background
./scripts/run-background.sh

# Stop server
./scripts/stop.sh
```

**Windows PowerShell:**
```powershell
# Database management
.\scripts\db-start.ps1         # Start PostgreSQL
.\scripts\db-stop.ps1          # Stop PostgreSQL

# Build the application
.\scripts\build.ps1

# Run tests
.\scripts\test.ps1

# Run tests with coverage
.\scripts\test-coverage.ps1

# Generate Swagger documentation
.\scripts\swagger.ps1

# Database migrations
.\scripts\migrate.ps1 up              # Apply all migrations
.\scripts\migrate.ps1 down            # Rollback last migration
.\scripts\migrate.ps1 create xxx      # Create new migration

# Run server in background
.\scripts\run-background.ps1

# Stop server
.\scripts\stop.ps1
```

## Configuration

The application supports flexible YAML-based configuration with multiple deployment stages, while maintaining full backward compatibility with environment variables.

### Configuration Methods

#### 1. YAML Configuration Files (Recommended)

Configuration files are organized by stage in the `config/` directory:

```
config/
├── base.yaml              # Base/shared configuration
├── development.yaml       # Development overrides
├── staging.yaml          # Staging overrides
└── production.yaml       # Production overrides
```

**Stage Selection:**

Using command line flag:
```bash
./server --stage=production
```

Using environment variable:
```bash
export APP_STAGE=production
./server
```

Default stage is `development` if not specified.

#### 2. Environment Variables (Backward Compatible)

All configuration values can be overridden with environment variables:

```bash
export DATABASE_URL="postgres://user:pass@host:5432/db"
export JWT_SECRET="your-secret-key"
export SERVER_PORT="8080"
./server
```

Environment variables take precedence over YAML configuration.

### Configuration Structure

**Base Configuration (config/base.yaml):**
```yaml
server:
  port: "8080"

database:
  url: "postgres://user:password@localhost:5432/myapp?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: 30

jwt:
  secret: "your-secret-key"

rate_limit:
  requests_per_second: 100
  burst: 200
```

**Stage-Specific Overrides:**

Each stage file (development.yaml, staging.yaml, production.yaml) can override base values:

```yaml
# config/production.yaml
database:
  url: "${DATABASE_URL}"  # Use environment variable
  max_open_conns: 100
  max_idle_conns: 25
  conn_max_lifetime: 60

jwt:
  secret: "${JWT_SECRET}"  # Use environment variable

server:
  port: "${SERVER_PORT:8080}"  # Default to 8080 if not set

rate_limit:
  requests_per_second: 50  # More conservative for production
  burst: 100
```

### Configuration Loading Order

1. Load `config/base.yaml` (shared defaults)
2. Merge `config/{stage}.yaml` (stage-specific overrides)
3. Apply environment variable overrides (highest priority)

### Environment Variable Mapping

YAML keys map to environment variables using underscores:

| YAML Path | Environment Variable |
|-----------|---------------------|
| `server.port` | `SERVER_PORT` |
| `database.url` | `DATABASE_URL` |
| `database.max_open_conns` | `DB_MAX_OPEN_CONNS` |
| `database.max_idle_conns` | `DB_MAX_IDLE_CONNS` |
| `database.conn_max_lifetime` | `DB_CONN_MAX_LIFETIME` |
| `jwt.secret` | `JWT_SECRET` |
| `rate_limit.requests_per_second` | `RATE_LIMIT_REQUESTS_PER_SECOND` |
| `rate_limit.burst` | `RATE_LIMIT_BURST` |

### Configuration Examples

**Development (default):**
```bash
# Uses config/development.yaml
./server
```

**Staging with secret override:**
```bash
export JWT_SECRET="staging-secret-key"
./server --stage=staging
```

**Production with all secrets from environment:**
```bash
export APP_STAGE=production
export DATABASE_URL="postgres://prod-user:pass@prod-host:5432/proddb"
export JWT_SECRET="production-secret-key"
export SERVER_PORT="8080"
./server
```

**Backward compatible (environment only):**
```bash
# No YAML files needed - works as before
export DATABASE_URL="postgres://user:pass@localhost:5432/db"
export JWT_SECRET="secret"
export SERVER_PORT="8080"
./server
```

### Security Best Practices

- ✅ Store secrets (DATABASE_URL, JWT_SECRET) in environment variables, not in YAML files
- ✅ Use placeholders in YAML: `${JWT_SECRET}` to reference environment variables
- ✅ Never commit production secrets to version control
- ✅ Use different secrets for each stage
- ✅ In Kubernetes, use Secrets for sensitive values

## API Documentation

### Interactive API Documentation (Swagger)

The API documentation is automatically generated from code annotations and available at:
- **Swagger UI**: http://localhost:8080/swagger/index.html

Features:
- Interactive API explorer
- Try out endpoints directly from the browser
- Auto-generated from code (always up-to-date)
- Request/response examples

To regenerate Swagger docs after code changes:
```bash
# Linux/macOS
./scripts/swagger.sh

# Windows PowerShell
.\scripts\swagger.ps1

# Or directly with swag CLI
swag init -g cmd/server/main.go --output docs --parseDependency --parseInternal
```

### API Versioning

The API supports versioning for backward compatibility:
- **v1 endpoints**: `/v1/login`, `/v1/users`, `/v1/users/{id}` (recommended)
- **Legacy endpoints**: `/login`, `/users`, `/users/{id}` (for backward compatibility)

All new integrations should use the v1 endpoints.

### OpenAPI Specification

```yaml
openapi: 3.0.0
info:
  title: User Management API
  version: 1.0.0
  description: Production-ready user management microservice

servers:
  - url: http://localhost:8080
    description: Local development

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
        email:
          type: string
          format: email
        role:
          type: string
          enum: [user, admin]
    
    CreateUserRequest:
      type: object
      required:
        - name
        - email
        - password
      properties:
        name:
          type: string
        email:
          type: string
          format: email
        password:
          type: string
          minLength: 6
        role:
          type: string
          enum: [user, admin]
    
    LoginRequest:
      type: object
      required:
        - email
        - password
      properties:
        email:
          type: string
          format: email
        password:
          type: string
    
    LoginResponse:
      type: object
      properties:
        token:
          type: string
        user:
          $ref: '#/components/schemas/User'

paths:
  /health:
    get:
      summary: Overall health check
      description: Returns aggregated health status with all component checks
      responses:
        '200':
          description: All components are healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    enum: [UP, DOWN]
                  components:
                    type: object
                    properties:
                      database:
                        type: object
                        properties:
                          status:
                            type: string
                            enum: [UP, DOWN]
                          details:
                            type: object
        '503':
          description: One or more components are unhealthy

  /health/startup:
    get:
      summary: Kubernetes startup probe
      description: Indicates if the application has started successfully
      responses:
        '200':
          description: Application has started
        '503':
          description: Application has not started

  /health/liveness:
    get:
      summary: Kubernetes liveness probe
      description: Indicates if the application is running and should not be restarted
      responses:
        '200':
          description: Application is alive

  /health/readiness:
    get:
      summary: Kubernetes readiness probe
      description: Indicates if the application is ready to accept traffic
      responses:
        '200':
          description: Application is ready
        '503':
          description: Application is not ready

  /login:
    post:
      summary: Authenticate user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginRequest'
      responses:
        '200':
          description: Login successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/LoginResponse'
        '401':
          description: Invalid credentials

  /users:
    get:
      summary: List all users (Admin only)
      security:
        - BearerAuth: []
      responses:
        '200':
          description: List of users
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/User'
        '401':
          description: Unauthorized
        '403':
          description: Forbidden
    
    post:
      summary: Create user (Public signup)
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
      responses:
        '201':
          description: User created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '400':
          description: Invalid input

  /users/{id}:
    get:
      summary: Get user by ID (Owner or Admin)
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: User details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '404':
          description: User not found
    
    put:
      summary: Update user (Owner or Admin)
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                email:
                  type: string
                  format: email
      responses:
        '200':
          description: User updated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '404':
          description: User not found
    
    delete:
      summary: Delete user (Admin only)
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '204':
          description: User deleted
        '404':
          description: User not found

  /metrics:
    get:
      summary: Prometheus metrics
      responses:
        '200':
          description: Metrics in Prometheus format
```

## API Examples

### Using v1 API (Recommended)

### 1. Create a user (Public signup)

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "role": "user"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

### 3. Get all users (Admin only) with W3C trace context

```bash
TOKEN="your-jwt-token"
curl -X GET http://localhost:8080/v1/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
```

### 4. Get user by ID

```bash
TOKEN="your-jwt-token"
curl -X GET http://localhost:8080/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

### 5. Update user

```bash
TOKEN="your-jwt-token"
curl -X PUT http://localhost:8080/v1/users/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "email": "john.updated@example.com"
  }'
```

### 6. Delete user (Admin only)

```bash
TOKEN="your-jwt-token"
curl -X DELETE http://localhost:8080/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

### 7. Health checks

The application provides a **registrable health check system** with comprehensive endpoints for Kubernetes probes. Components can register custom health checks with specific scopes (base, startup, ready, live).

```bash
# Overall health check (aggregates all checks)
curl http://localhost:8080/health

# Startup probe (for Kubernetes startup probe)
curl http://localhost:8080/health/startup

# Liveness probe (for Kubernetes liveness probe)
curl http://localhost:8080/health/liveness

# Readiness probe (for Kubernetes readiness probe)
curl http://localhost:8080/health/readiness
```

Example response from `/health` endpoint:
```json
{
  "status": "UP",
  "components": {
    "database": {
      "status": "UP",
      "details": {
        "idle": 2,
        "in_use": 0,
        "max_open_connections": 25,
        "open_connections": 2
      }
    }
  }
}
```

**For detailed information on creating custom health checks**, see [docs/HEALTH_CHECKS.md](docs/HEALTH_CHECKS.md).

### 8. Prometheus metrics (includes users_total)

```bash
curl http://localhost:8080/metrics | grep -E "(http_requests_total|users_total)"
```

Example output:
```
http_requests_total{method="GET",path="/v1/users",status="200"} 42
users_total 156
```

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema version control.

### How Migrations Work

1. **Automatic**: Migrations run automatically when the application starts
2. **Version-controlled**: All schema changes are tracked in `migrations/` directory
3. **Rollback support**: Each migration has up and down SQL files
4. **Idempotent**: Can safely re-run migrations

### Migration Files

```
migrations/
├── 000001_create_users_table.up.sql    # Create users table
└── 000001_create_users_table.down.sql  # Drop users table
```

### Creating New Migrations

**Linux/macOS:**
```bash
./scripts/migrate.sh create add_user_profile
```

**Windows PowerShell:**
```powershell
.\scripts\migrate.ps1 create add_user_profile
```

**Or using migrate CLI directly:**
```bash
migrate create -ext sql -dir migrations -seq add_user_profile
```

This creates:
- `migrations/000002_add_user_profile.up.sql` - Apply changes
- `migrations/000002_add_user_profile.down.sql` - Revert changes

### Manual Migration Commands

**Linux/macOS:**
```bash
# Apply all pending migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down

# Create new migration
./scripts/migrate.sh create migration_name

# Force to specific version (recovery)
./scripts/migrate.sh force 1
```

**Windows PowerShell:**
```powershell
# Apply all pending migrations
.\scripts\migrate.ps1 up

# Rollback last migration
.\scripts\migrate.ps1 down

# Create new migration
.\scripts\migrate.ps1 create migration_name

# Force to specific version (recovery)
.\scripts\migrate.ps1 force 1
```

### Migration Best Practices

1. **Always test both directions**: up and down migrations
2. **Make migrations idempotent**: Use `IF NOT EXISTS`, `IF EXISTS`
3. **Never modify existing migrations**: Create new ones instead
4. **Keep migrations small**: One logical change per migration
5. **Review generated SQL**: Ensure it's safe for production

For detailed migration documentation, see [docs/database/migrations.md](docs/database/migrations.md)

## Testing

### Run all tests

```bash
go test ./... -v
```

### Run tests with coverage

```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Run specific test

```bash
go test ./internal/handlers -v -run TestGetUsers
```

### Test coverage report

```bash
go tool cover -func=coverage.out
```

## Docker

### Build image

```bash
docker build -t myapp:latest .
```

### Run container

```bash
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:password@host:5432/myapp" \
  -e JWT_SECRET="your-secret-key" \
  myapp:latest
```

### Docker Compose (with PostgreSQL)

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  db:
    image: postgres:13-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: myapp
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: "postgres://user:password@db:5432/myapp?sslmode=disable"
      JWT_SECRET: "your-secret-key"
    depends_on:
      - db

volumes:
  postgres_data:
```

Run:
```bash
docker-compose up
```

## Kubernetes Deployment

### Using Helm

1. **Create secrets:**

```bash
kubectl create secret generic myapp-secrets \
  --from-literal=database-url="postgres://user:password@postgres:5432/myapp" \
  --from-literal=jwt-secret="your-secret-key" \
  -n production
```

2. **Install with Helm:**

**Basic installation with production stage:**
```bash
helm install myapp ./helm/myapp \
  --namespace production \
  --create-namespace
```

**Install with specific stage:**
```bash
# Development
helm install myapp ./helm/myapp \
  --namespace development \
  --set config.stage=development \
  --create-namespace

# Staging
helm install myapp ./helm/myapp \
  --namespace staging \
  --set config.stage=staging \
  --create-namespace

# Production (default)
helm install myapp ./helm/myapp \
  --namespace production \
  --set config.stage=production \
  --create-namespace
```

**Install with ConfigMap-based configuration (optional):**
```bash
helm install myapp ./helm/myapp \
  --namespace production \
  --set config.stage=production \
  --set config.useConfigMap=true \
  --create-namespace
```

3. **Update deployment:**

```bash
# Update image version
helm upgrade myapp ./helm/myapp \
  --set image.tag=v1.0.1 \
  -n production

# Change configuration stage
helm upgrade myapp ./helm/myapp \
  --set config.stage=staging \
  -n production
```

4. **Uninstall:**

```bash
helm uninstall myapp -n production
```

### Helm Configuration Options

The Helm chart supports the following configuration values:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `config.stage` | Configuration stage (development, staging, production) | `production` |
| `config.useConfigMap` | Mount config files via ConfigMap instead of using env vars | `false` |
| `replicaCount` | Number of replicas | `2` |
| `image.repository` | Docker image repository | `myapp` |
| `image.tag` | Docker image tag | `latest` |
| `autoscaling.enabled` | Enable horizontal pod autoscaling | `true` |
| `autoscaling.minReplicas` | Minimum replicas | `2` |
| `autoscaling.maxReplicas` | Maximum replicas | `10` |

**Note:** When `config.stage` is set, the `APP_STAGE` environment variable is automatically configured in the pods.

### Local Testing with Kind

For local development and testing, you can deploy to a Kind cluster:

1. **Setup local Kubernetes cluster:**

```bash
# Linux/macOS
./scripts/local-k8s-setup.sh

# Windows PowerShell
.\scripts\local-k8s-setup.ps1
```

2. **Deploy to local cluster:**

```bash
# Linux/macOS
./scripts/local-k8s-deploy.sh

# Windows PowerShell
.\scripts\local-k8s-deploy.ps1
```

3. **Access the application:**

```bash
kubectl port-forward -n production svc/myapp 8080:8080
```

4. **Test the API:**

```bash
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

5. **Cleanup:**

```bash
# Linux/macOS
./scripts/local-k8s-cleanup.sh

# Windows PowerShell
.\scripts\local-k8s-cleanup.ps1
```

**Troubleshooting**: See [Local Kubernetes Troubleshooting Guide](docs/deployment/LOCAL_K8S_TROUBLESHOOTING.md) for common issues and solutions.

### Custom values

Create `my-values.yaml`:

```yaml
replicaCount: 3

image:
  repository: ghcr.io/huberp/go-vibe
  tag: "v1.0.0"

resources:
  requests:
    cpu: 500m
    memory: 512Mi
  limits:
    cpu: 1000m
    memory: 1Gi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
  targetCPUUtilizationPercentage: 75

database:
  url: "postgres://user:password@postgres.production.svc.cluster.local:5432/myapp"

jwt:
  secret: "production-secret-key"
```

Install with custom values:
```bash
helm install myapp ./helm/myapp -f my-values.yaml -n production
```

## CI/CD

The project includes GitHub Actions workflows for building, testing, deployment, and dependency management:

### Core Workflows

#### 1. Build (`build.yml`)
- Triggers on push to main/develop or PR to main
- Builds the application
- Uploads binary artifact

#### 2. Test (`test.yml`)
- Triggers on push to main/develop or PR to main
- Runs all tests with race detection
- Generates coverage report
- Uploads to Codecov

#### 3. Deploy (`deploy.yml`)
- Triggers on push to main or version tags
- Builds and pushes Docker image to GHCR
- Deploys to Kubernetes using Helm

### Dependency Management Workflows

#### 4. Combine Dependency PRs (`combine-dependency-prs.yml`)
- **Purpose**: Combines multiple chore(deps) PRs into a single PR
- **Trigger**: Manual (workflow_dispatch)
- **Filters**: Only PRs with labels `dependencies` AND `go` that pass Build and Test
- **Benefits**: Reduces PR clutter, easier review process

**Setup Required (One-Time):**
Before using this workflow, enable PR creation in repository settings:
1. Go to **Settings** → **Actions** → **General**
2. Under **Workflow permissions**, select **"Read and write permissions"**
3. Check **"Allow GitHub Actions to create and approve pull requests"**
4. Click **Save**

**How to use:**
1. Go to Actions tab → "Combine Dependency PRs"
2. Click "Run workflow"
3. Review and merge the combined PR
4. Source PRs are automatically closed after merge

See [Combine Dependency PRs Documentation](docs/development/COMBINE_DEPENDENCY_PRS.md) for detailed usage.

#### 5. Cleanup Combined PRs (`cleanup-combined-prs.yml`)
- **Purpose**: Automatically closes source PRs after combined PR is merged
- **Trigger**: Automatic when combined PR is merged
- **Actions**: Closes source PRs, deletes branches, adds comments

### Setup GitHub Secrets

Add these secrets to your GitHub repository:

- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: JWT signing secret
- `KUBECONFIG`: Kubernetes configuration (base64 encoded)

## Security

- ✅ Passwords hashed with bcrypt
- ✅ JWT tokens for authentication (HS256)
- ✅ Role-based access control
- ✅ Input validation
- ✅ SQL injection prevention (GORM parameterized queries)
- ✅ Non-root Docker user
- ✅ Secrets stored in Kubernetes secrets

## Monitoring

### Health Checks

The application provides comprehensive health check endpoints following Kubernetes best practices:

#### Available Endpoints

1. **Overall Health Check** (`/health`)
   - Aggregates all component health statuses
   - Returns `200 OK` when all components are healthy
   - Returns `503 Service Unavailable` when any component is unhealthy
   - Includes detailed component information (database connection stats, etc.)

2. **Startup Probe** (`/health/startup`)
   - Indicates if the application has successfully started
   - Used by Kubernetes to know when the container is ready to start accepting traffic
   - Checks database connectivity
   - Configured with longer timeout to allow for slow startup

3. **Liveness Probe** (`/health/liveness`)
   - Indicates if the application is running and should not be restarted
   - Simple check that doesn't depend on external services
   - Used by Kubernetes to determine when to restart the container
   - Always returns `200 OK` if the application can respond

4. **Readiness Probe** (`/health/readiness`)
   - Indicates if the application is ready to accept traffic
   - Checks database connectivity and other critical dependencies
   - Used by Kubernetes to determine when to route traffic to the pod
   - Returns `503 Service Unavailable` if dependencies are not ready

#### Health Check Response Format

```json
{
  "status": "UP",
  "components": {
    "database": {
      "status": "UP",
      "details": {
        "max_open_connections": 25,
        "open_connections": 2,
        "in_use": 0,
        "idle": 2
      }
    }
  }
}
```

Status values: `UP` (healthy) or `DOWN` (unhealthy)

### Prometheus Metrics

The application exposes the following metrics at `/metrics`:

#### HTTP Metrics
- `http_requests_total`: Total HTTP requests (labeled by method, path, status)
- `http_request_duration_seconds`: HTTP request duration histogram
- `users_total`: Total number of users in the database (gauge)

#### Go Runtime Metrics (runtime.MemStats)
- `go_memstats_alloc_bytes`: Bytes of allocated heap objects
- `go_memstats_sys_bytes`: Total bytes of memory obtained from OS
- `go_memstats_heap_alloc_bytes`: Heap bytes allocated and still in use
- `go_memstats_heap_sys_bytes`: Heap memory obtained from OS
- `go_memstats_heap_idle_bytes`: Heap bytes waiting to be used
- `go_memstats_heap_inuse_bytes`: Heap bytes that are in use
- `go_memstats_heap_released_bytes`: Heap bytes released to OS
- `go_memstats_heap_objects`: Number of allocated heap objects
- `go_memstats_mallocs_total`: Total number of heap allocations
- `go_memstats_frees_total`: Total number of heap frees
- `go_memstats_gc_sys_bytes`: Bytes used for garbage collection metadata
- `go_goroutines`: Number of goroutines
- `go_threads`: Number of OS threads
- `go_gc_duration_seconds`: GC duration distribution
- `go_info`: Go version information

### Structured Logging

All requests are logged with:
- `request_id`: Unique request identifier (UUID or W3C trace ID)
- `method`: HTTP method
- `path`: Request path
- `status`: Response status code
- `duration`: Request duration
- `client_ip`: Client IP address
- `trace_id`: OpenTelemetry trace ID (if available)
- `span_id`: OpenTelemetry span ID (if available)

### W3C Trace Context Support

The logging middleware supports W3C trace context propagation:
- Accepts `traceparent` header for distributed tracing
- Extracts trace ID from W3C traceparent format
- Integrates with OpenTelemetry for end-to-end tracing

### Rate Limiting

Rate limiting is enforced per IP address:
- **Default limit**: 100 requests per second
- **Burst capacity**: 200 requests
- **Response**: HTTP 429 (Too Many Requests) when limit exceeded

### CORS Configuration

CORS middleware is configured to:
- Allow all origins (configure for production)
- Support methods: GET, POST, PUT, DELETE, OPTIONS
- Accept headers: Origin, Content-Type, Authorization, traceparent, tracestate
- Support credentials for authenticated requests

**Production Note**: Update `AllowOrigins` in routes.go to restrict to your frontend domain(s).

## Security

- **Password Hashing**: bcrypt with cost factor 12 (enhanced from default 10)
- **JWT Authentication**: HS256 algorithm with configurable secret
- **Input Validation**: Gin validator tags on all request models
- **SQL Injection Prevention**: GORM parameterized queries
- **Rate Limiting**: Per-IP request limiting to prevent abuse
- **CORS**: Configurable cross-origin policy

## Performance

- Database connection pooling (GORM)
- Horizontal pod autoscaling (Kubernetes)
- Resource limits and requests
- Comprehensive health checks (startup, liveness, readiness probes)
- Multi-stage Docker build for minimal image size

## Contributing

We welcome contributions to go-vibe! To ensure consistency and quality across the codebase, please follow our comprehensive development guidelines.

### Development Guidelines

All coding standards, best practices, and architectural patterns are documented in our **[Copilot Instructions](.github/copilot-instructions.md)**. This file provides:

- 📋 Project overview and tech stack
- 🏗️ Architecture patterns and design principles
- 📝 Code style and naming conventions
- ✅ Testing strategy (TDD approach)
- 🔒 Security best practices
- 🗄️ Database and GORM guidelines
- 🚀 Development workflow and common commands
- 📚 External resources and documentation

**Quick Start for Contributors:**

1. Read the [Copilot Instructions](.github/copilot-instructions.md)
2. Follow the TDD approach: write tests first
3. Ensure all tests pass: `go test ./... -v`
4. Run code coverage: `go test ./... -coverprofile=coverage.out`
5. Follow commit message conventions (Conventional Commits)
6. Submit PR following the guidelines in the instructions

### Code Review Checklist

Before submitting a PR, ensure:
- [ ] Tests written and passing (TDD approach)
- [ ] Error handling implemented correctly
- [ ] Logging added with appropriate context
- [ ] Security best practices followed
- [ ] Input validation included
- [ ] Documentation updated (if needed)
- [ ] No sensitive data in logs or responses
- [ ] Code follows Go conventions
- [ ] Dependencies are justified and minimal

## License

MIT

## Support

For issues and questions, please open an issue on GitHub.
