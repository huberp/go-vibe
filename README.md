# go-vibe

[![Build](https://github.com/huberp/go-vibe/workflows/Build/badge.svg)](https://github.com/huberp/go-vibe/actions/workflows/build.yml)
[![Test](https://github.com/huberp/go-vibe/workflows/Test/badge.svg)](https://github.com/huberp/go-vibe/actions/workflows/test.yml)
[![Deploy](https://github.com/huberp/go-vibe/workflows/Deploy/badge.svg)](https://github.com/huberp/go-vibe/actions/workflows/deploy.yml)
[![codecov](https://codecov.io/gh/huberp/go-vibe/branch/main/graph/badge.svg)](https://codecov.io/gh/huberp/go-vibe)
[![Go Version](https://img.shields.io/badge/Go-1.25.2-blue.svg)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/huberp/go-vibe)](https://goreportcard.com/report/github.com/huberp/go-vibe)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**A production-ready user management microservice built with Go, following Test-Driven Development (TDD) principles and designed for cloud-native Kubernetes deployment.**

## Features

- ✅ **RESTful API** with JWT authentication
- ✅ **Role-based Access Control** (RBAC) - admin and user roles
- ✅ **PostgreSQL** database with GORM ORM
- ✅ **Database migrations** with golang-migrate
- ✅ **OpenAPI/Swagger** documentation (auto-generated)
- ✅ **Prometheus metrics** for monitoring
- ✅ **Structured logging** with Zap
- ✅ **OpenTelemetry** tracing support
- ✅ **Rate limiting** (configurable per environment)
- ✅ **CORS** middleware
- ✅ **Docker** containerization
- ✅ **Kubernetes/Helm** deployment
- ✅ **CI/CD** with GitHub Actions
- ✅ **100% test coverage** for critical paths

## Quick Start

### Prerequisites

- Go 1.25.2 or higher
- PostgreSQL 13+
- Docker (optional, for containerized deployment)
- Kubernetes cluster (optional, for production deployment)

### Local Development

1. **Clone the repository:**
   ```bash
   git clone https://github.com/huberp/go-vibe.git
   cd go-vibe
   ```

2. **Set up environment variables:**
   ```bash
   export DATABASE_URL="postgres://user:password@localhost:5432/myapp?sslmode=disable"
   export JWT_SECRET="your-secret-key-here"
   export SERVER_PORT="8080"
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Generate Swagger documentation:**
   ```bash
   # Linux/macOS
   ./scripts/swagger.sh
   
   # Windows PowerShell
   .\scripts\swagger.ps1
   ```

5. **Run the application:**
   ```bash
   go run ./cmd/server
   ```

6. **Access the API:**
   - Health check: http://localhost:8080/health
   - Swagger UI: http://localhost:8080/swagger/index.html
   - Metrics: http://localhost:8080/metrics

### Using Docker Compose

```bash
docker-compose up -d
```

This will start both the application and PostgreSQL database.

## API Endpoints

### v1 API (Recommended)

| Method | Endpoint         | Auth            | Description          |
|--------|------------------|-----------------|----------------------|
| POST   | `/v1/users`      | None (signup)   | Create a user        |
| POST   | `/v1/login`      | None            | Authenticate user    |
| GET    | `/v1/users`      | JWT (admin)     | List all users       |
| GET    | `/v1/users/{id}` | JWT (owner/admin) | Get user by ID     |
| PUT    | `/v1/users/{id}` | JWT (owner/admin) | Update user by ID  |
| DELETE | `/v1/users/{id}` | JWT (admin)     | Delete user by ID    |

### Monitoring & Documentation

| Method | Endpoint       | Auth | Description                    |
|--------|----------------|------|--------------------------------|
| GET    | `/health`      | None | Health check                   |
| GET    | `/metrics`     | None | Prometheus metrics             |
| GET    | `/swagger/*`   | None | OpenAPI/Swagger documentation  |

For detailed API documentation, see the [Swagger UI](http://localhost:8080/swagger/index.html) when running locally.

## Configuration

The application supports multiple configuration methods:

1. **Environment Variables** (highest priority)
2. **YAML Configuration Files** (stage-specific)
3. **Default Values** (fallback)

### Environment Variables

```bash
DATABASE_URL=postgres://user:password@localhost:5432/myapp?sslmode=disable
JWT_SECRET=your-secret-key
SERVER_PORT=8080
APP_STAGE=production  # Options: development, staging, production
```

### YAML Configuration

Configuration files are located in `config/`:
- `base.yaml` - Shared defaults
- `development.yaml` - Development overrides
- `staging.yaml` - Staging overrides
- `production.yaml` - Production overrides

See [docs/configuration/yaml-config-migration.md](docs/configuration/yaml-config-migration.md) for details.

## Testing

### Run all tests:

```bash
# Linux/macOS
./scripts/test.sh

# Windows PowerShell
.\scripts\test.ps1

# Manual
go test ./... -v
```

### Run tests with coverage:

```bash
# Linux/macOS
./scripts/test-coverage.sh

# Windows PowerShell
.\scripts\test-coverage.ps1

# Manual
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test API endpoints:

```bash
# Linux/macOS
./test-api.sh

# Windows PowerShell
.\test-api.ps1
```

## Database Migrations

### Create a new migration:

```bash
# Linux/macOS
./scripts/migrate.sh create <migration_name>

# Windows PowerShell
.\scripts\migrate.ps1 create <migration_name>
```

### Apply migrations:

```bash
# Linux/macOS
./scripts/migrate.sh up

# Windows PowerShell
.\scripts\migrate.ps1 up
```

See [docs/database/migrations.md](docs/database/migrations.md) for detailed migration guide.

## Deployment

### Build Docker Image

```bash
docker build -t go-vibe:latest .
```

### Deploy to Kubernetes with Helm

```bash
helm install myapp ./helm/myapp \
  --set database.url="$DATABASE_URL" \
  --set jwt.secret="$JWT_SECRET" \
  --namespace production \
  --create-namespace
```

For local Kubernetes setup, see [docs/deployment/LOCAL_K8S_SETUP_SUMMARY.md](docs/deployment/LOCAL_K8S_SETUP_SUMMARY.md).

## Observability

### Prometheus Metrics

The application exposes the following metrics at `/metrics`:

- `http_requests_total` - Total HTTP requests by method, path, and status
- `http_request_duration_seconds` - HTTP request duration
- `users_total` - Total number of users (gauge)
- Go runtime metrics (memory, GC, goroutines, etc.)

### Structured Logging

All logs are structured using Zap with the following fields:
- `method` - HTTP method
- `path` - Request path
- `status` - Response status code
- `duration` - Request duration
- `client_ip` - Client IP address
- `trace_id` - W3C trace ID (if available)
- `span_id` - OpenTelemetry span ID (if available)

### Tracing

The application supports W3C Trace Context and OpenTelemetry tracing:
- Automatic trace ID extraction from `traceparent` header
- Span context propagation
- Integration with OpenTelemetry backends

## Project Structure

```
.
├── .github/workflows/       # CI/CD pipelines
├── cmd/server/             # Application entry point
├── config/                 # YAML configuration files
├── docs/                   # Documentation and Swagger specs
├── helm/myapp/            # Kubernetes Helm chart
├── internal/              # Private application code
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # HTTP middleware (auth, logging, metrics)
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   └── routes/            # Route configuration
├── migrations/            # Database migrations
├── pkg/                   # Public libraries
│   ├── config/            # Configuration loader
│   ├── logger/            # Logging setup
│   ├── migration/         # Migration runner
│   └── utils/             # Utilities (JWT, password hashing)
├── scripts/               # Build and deployment scripts
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yml     # Local development setup
└── go.mod                 # Go module dependencies
```

## Development Scripts

All common development tasks have helper scripts:

```bash
# Linux/macOS
./scripts/build.sh              # Build the application
./scripts/test.sh               # Run tests
./scripts/test-coverage.sh      # Run tests with coverage
./scripts/swagger.sh            # Generate Swagger docs
./scripts/migrate.sh            # Database migrations
./scripts/run-background.sh     # Start server in background
./scripts/stop.sh               # Stop background server

# Windows PowerShell
.\scripts\build.ps1
.\scripts\test.ps1
.\scripts\test-coverage.ps1
.\scripts\swagger.ps1
.\scripts\migrate.ps1
.\scripts\run-background.ps1
.\scripts\stop.ps1
```

See [scripts/README.md](scripts/README.md) for detailed script documentation.

## Security

- ✅ JWT (HS256) token-based authentication
- ✅ Password hashing with bcrypt (cost factor: 12)
- ✅ Input validation on all endpoints
- ✅ SQL injection prevention (GORM parameterized queries)
- ✅ Role-based access control (RBAC)
- ✅ Configurable rate limiting
- ✅ CORS middleware

**Security Best Practices:**
- Never commit secrets to version control
- Store JWT_SECRET and DATABASE_URL in environment variables or Kubernetes secrets
- Use HTTPS in production (configure in Kubernetes ingress)
- Regularly update dependencies

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes (TDD approach)
4. Ensure all tests pass (`./scripts/test.sh`)
5. Commit your changes with conventional commits format
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

See [.github/copilot-instructions.md](.github/copilot-instructions.md) for detailed contribution guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Documentation

For comprehensive documentation, see the [docs/](docs/) directory:

- **API Documentation**: [docs/api/](docs/api/) - OpenAPI/Swagger specifications
- **Configuration Guide**: [docs/configuration/](docs/configuration/) - YAML config and migration
- **Database Guide**: [docs/database/](docs/database/) - Migrations and PostgreSQL setup
- **Deployment Guide**: [docs/deployment/](docs/deployment/) - Kubernetes and local setup
- **Development Guide**: [docs/development/](docs/development/) - Code review and workflows
- **Observability**: [docs/observability/](docs/observability/) - Metrics and monitoring

## Architecture

This microservice follows Clean Architecture principles:

- **Repository Pattern** for data access abstraction
- **Dependency Injection** for testability
- **Interface-based design** for flexibility
- **Separation of concerns** across layers
- **Test-Driven Development** (TDD) for all features

## Tech Stack

- **Language**: Go 1.25.2
- **Web Framework**: Gin v1.11.0
- **ORM**: GORM v1.31.0
- **Database**: PostgreSQL 13+
- **Authentication**: JWT (github.com/golang-jwt/jwt/v5)
- **Logging**: Zap v1.27.0
- **Testing**: Testify v1.11.1 + gomock
- **Metrics**: Prometheus v1.23.2
- **Migrations**: golang-migrate v4.19.0
- **Documentation**: Swagger/OpenAPI (swaggo)
- **Deployment**: Docker, Kubernetes, Helm

## Acknowledgments

Built with Test-Driven Development (TDD) principles and following Go best practices.

---

For more information, see the [Implementation Summary](IMPLEMENTATION_SUMMARY.md).
