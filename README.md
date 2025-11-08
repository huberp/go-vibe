# go-vibe

[![Build](https://github.com/huberp/go-vibe/workflows/Build/badge.svg)](https://github.com/huberp/go-vibe/actions/workflows/build.yml)
[![Test](https://github.com/huberp/go-vibe/workflows/Test/badge.svg)](https://github.com/huberp/go-vibe/actions/workflows/test.yml)
[![Deploy](https://github.com/huberp/go-vibe/workflows/Deploy/badge.svg)](https://github.com/huberp/go-vibe/actions/workflows/deploy.yml)
[![codecov](https://codecov.io/gh/huberp/go-vibe/branch/main/graph/badge.svg)](https://codecov.io/gh/huberp/go-vibe)
[![Go Version](https://img.shields.io/badge/Go-1.25.2-blue.svg)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/huberp/go-vibe)](https://goreportcard.com/report/github.com/huberp/go-vibe)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**A production-ready Go microservice template with observability, testing, and deployment best practices built-in.**

This template provides a solid foundation for building cloud-native microservices with Go, following Test-Driven Development (TDD) principles and industry best practices.

## ✨ What's Included

### Core Infrastructure
- ✅ **RESTful API Framework** - Gin web framework with clean architecture
- ✅ **PostgreSQL Integration** - GORM ORM with connection pooling
- ✅ **Database Migrations** - Version-controlled schema management with golang-migrate
- ✅ **OpenAPI/Swagger** - Auto-generated API documentation
- ✅ **JWT Authentication** - Secure token-based authentication middleware
- ✅ **Role-Based Access Control** - RBAC middleware for authorization

### Observability
- ✅ **Prometheus Metrics** - HTTP request metrics and custom collectors
- ✅ **Structured Logging** - Zap logger with request correlation
- ✅ **OpenTelemetry Tracing** - Distributed tracing support with W3C trace context
- ✅ **Health Checks** - Kubernetes-ready startup, liveness, and readiness probes
- ✅ **Info Endpoint** - Runtime information and custom stats

### Reliability & Security
- ✅ **Rate Limiting** - Configurable per-environment rate limiting
- ✅ **CORS Middleware** - Cross-Origin Resource Sharing support
- ✅ **Input Validation** - Request validation with Gin binding
- ✅ **Error Handling** - Consistent error responses
- ✅ **Recovery Middleware** - Panic recovery

### Development & Deployment
- ✅ **Docker Containerization** - Multi-stage build for minimal image size
- ✅ **Kubernetes/Helm Charts** - Production-ready K8s deployment
- ✅ **CI/CD Pipelines** - GitHub Actions for build, test, and deploy
- ✅ **YAML Configuration** - Stage-specific configuration (dev, staging, production)
- ✅ **Development Scripts** - Shell and PowerShell scripts for common tasks
- ✅ **Test-Driven Development** - Comprehensive test suite with >85% coverage

## 🚀 Quick Start

### Prerequisites

- Go 1.25.2 or higher
- PostgreSQL 13+ (or use Docker Compose)
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
   - Info: http://localhost:8080/info

### Using Docker Compose

```bash
docker-compose up -d
```

This will start both the application and PostgreSQL database.

## 📚 Core API Endpoints

### Monitoring & Documentation

| Method | Endpoint             | Description                        |
|--------|----------------------|------------------------------------|
| GET    | `/health`            | Overall health check               |
| GET    | `/health/startup`    | Kubernetes startup probe           |
| GET    | `/health/liveness`   | Kubernetes liveness probe          |
| GET    | `/health/readiness`  | Kubernetes readiness probe         |
| GET    | `/info`              | Build info and runtime statistics  |
| GET    | `/metrics`           | Prometheus metrics                 |
| GET    | `/swagger/*`         | OpenAPI/Swagger documentation      |

### Adding Your Own API

The template provides infrastructure - add your domain-specific API endpoints in `internal/routes/routes.go`. See the [examples/user-management](examples/user-management/) for a complete working example.

## 📖 Examples

### User Management API

A complete example demonstrating how to build a user management microservice with authentication, authorization, and CRUD operations.

**Location:** [examples/user-management/](examples/user-management/)

**Features:**
- User registration and authentication
- JWT token generation and validation
- Role-based access control (admin/user)
- CRUD operations with owner-based permissions
- Database migrations
- Comprehensive tests
- API testing scripts

See [examples/user-management/README.md](examples/user-management/README.md) for detailed documentation.

## 🏗️ Project Structure

```
.
├── .github/workflows/       # CI/CD pipelines
├── cmd/server/             # Application entry point
├── config/                 # YAML configuration files
├── docs/                   # Documentation and Swagger specs
├── examples/               # Example implementations
│   └── user-management/    # Complete user management example
├── helm/myapp/            # Kubernetes Helm chart
├── internal/              # Private application code
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # HTTP middleware (auth, logging, metrics)
│   ├── models/            # Data models (add your own here)
│   ├── repository/        # Data access layer (add your own here)
│   └── routes/            # Route configuration
├── migrations/            # Database migrations (add your own here)
├── pkg/                   # Public libraries
│   ├── config/            # Configuration loader
│   ├── health/            # Health check system
│   ├── info/              # Info endpoint system
│   ├── logger/            # Logging setup
│   ├── migration/         # Migration runner
│   └── utils/             # Utilities (JWT, password hashing)
├── scripts/               # Build and deployment scripts
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yml     # Local development setup
└── go.mod                 # Go module dependencies
```

## ⚙️ Configuration

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

## 🧪 Testing

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

## 📊 Observability

### Prometheus Metrics

The application exposes the following metrics at `/metrics`:

- `http_requests_total` - Total HTTP requests by method, path, and status
- `http_request_duration_seconds` - HTTP request duration
- Go runtime metrics (memory, GC, goroutines, etc.)

Add your own custom metrics using the Prometheus client library.

### Structured Logging

All logs are structured using Zap with the following fields:
- `method` - HTTP method
- `path` - Request path
- `status` - Response status code
- `duration` - Request duration
- `client_ip` - Client IP address
- `trace_id` - W3C trace ID (if available)
- `span_id` - OpenTelemetry span ID (if available)

### Health Checks

The template includes a flexible health check system with scope-based providers:

- `/health` - Overall health (all providers)
- `/health/startup` - Startup readiness (database, critical dependencies)
- `/health/liveness` - Application is running (no database checks)
- `/health/readiness` - Ready to accept traffic (database, dependencies)

See [docs/HEALTH_CHECKS.md](docs/HEALTH_CHECKS.md) for details.

### Tracing

The application supports W3C Trace Context and OpenTelemetry tracing:
- Automatic trace ID extraction from `traceparent` header
- Span context propagation
- Integration with OpenTelemetry backends

## 🗄️ Database Migrations

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

## 🚢 Deployment

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

## 🛠️ Development Scripts

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

## 🔒 Security

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

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes (TDD approach)
4. Ensure all tests pass (`./scripts/test.sh`)
5. Commit your changes with conventional commits format
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

See [.github/copilot-instructions.md](.github/copilot-instructions.md) for detailed contribution guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📚 Documentation

For comprehensive documentation, see the [docs/](docs/) directory:

- **API Documentation**: [docs/api/](docs/api/) - OpenAPI/Swagger specifications
- **Configuration Guide**: [docs/configuration/](docs/configuration/) - YAML config and migration
- **Database Guide**: [docs/database/](docs/database/) - Migrations and PostgreSQL setup
- **Deployment Guide**: [docs/deployment/](docs/deployment/) - Kubernetes and local setup
- **Development Guide**: [docs/development/](docs/development/) - Code review and workflows
- **Observability**: [docs/observability/](docs/observability/) - Metrics and monitoring
- **Health Checks**: [docs/HEALTH_CHECKS.md](docs/HEALTH_CHECKS.md) - Health check system

## 🏗️ Architecture

This microservice template follows Clean Architecture principles:

- **Repository Pattern** for data access abstraction
- **Dependency Injection** for testability
- **Interface-based design** for flexibility
- **Separation of concerns** across layers
- **Test-Driven Development** (TDD) for all features

## 🛠️ Tech Stack

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

## 🚀 Getting Started with Your Own API

1. **Define your domain models** in `internal/models/`
2. **Create repository interfaces and implementations** in `internal/repository/`
3. **Build your handlers** in `internal/handlers/`
4. **Add routes** in `internal/routes/routes.go`
5. **Create database migrations** in `migrations/`
6. **Write tests** for all components (TDD approach)
7. **Update Swagger annotations** for API documentation
8. **Run tests and build** to verify everything works

For a complete example, see [examples/user-management/](examples/user-management/).

## 🙏 Acknowledgments

Built with Test-Driven Development (TDD) principles and following Go best practices.

---

For more information, see the [Implementation Summary](IMPLEMENTATION_SUMMARY.md).
