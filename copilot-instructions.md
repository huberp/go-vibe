# Copilot Instructions for go-vibe

This document provides overarching guidelines and best practices for contributing to the go-vibe project. These instructions apply across the entire codebase.

## Project Overview

go-vibe is a production-ready user management microservice built with Go 1.24, following Test-Driven Development (TDD) principles and designed for cloud-native Kubernetes deployment.

## Tech Stack

- **Language**: Go 1.24+
- **Web Framework**: Gin v1.10.0
- **ORM**: GORM v1.25.10
- **Database**: PostgreSQL 13+
- **Authentication**: JWT (github.com/golang-jwt/jwt/v5 v5.2.1)
- **Logging**: Zap v1.27.0
- **Testing**: Testify v1.9.0 + gomock
- **Metrics**: Prometheus v1.19.1
- **Deployment**: Docker, Kubernetes, Helm

## Architecture Patterns

### Clean Architecture
- Use repository pattern for data access layer
- Implement dependency injection
- Design with interfaces for testability
- Maintain clear separation of concerns:
  - `cmd/`: Application entry points
  - `internal/handlers/`: HTTP request handlers
  - `internal/models/`: Data models
  - `internal/repository/`: Data access layer
  - `internal/middleware/`: HTTP middleware
  - `internal/routes/`: Route configuration
  - `pkg/`: Reusable packages (config, logger, utils)

### Request Flow
1. Request → Middleware (logging, auth) → Handler → Repository → Database
2. Use context propagation throughout the request chain
3. Store user information in Gin context after authentication

## Code Style and Standards

### General Go Guidelines
- Follow standard Go conventions and idioms
- Use `gofmt` for code formatting
- Keep functions focused and small
- Prefer explicit error handling over panic
- Use meaningful variable and function names

### Error Handling
- Return errors explicitly, don't use panic
- Use custom error types when appropriate (e.g., `ErrUserNotFound`)
- Provide appropriate HTTP status codes:
  - `200 OK`: Successful GET/PUT
  - `201 Created`: Successful POST
  - `204 No Content`: Successful DELETE
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Missing or invalid authentication
  - `403 Forbidden`: Insufficient permissions
  - `404 Not Found`: Resource not found
  - `500 Internal Server Error`: Server-side errors
- Return consistent error response format: `{"error": "message"}`
- Log errors with structured logging before returning to client

### Logging
- Use structured logging with Zap
- Include relevant context in logs (request_id, user_id, etc.)
- Log levels:
  - `Info`: Normal operations, request/response
  - `Error`: Errors that need attention
  - `Fatal`: Unrecoverable errors (startup only)
- Never log sensitive data (passwords, tokens, etc.)
- Request logging should include: method, path, status, duration, client_ip

## Testing Strategy

### Test-Driven Development (TDD)
- **Write tests before implementation**
- Follow the Red-Green-Refactor cycle
- Tests define the expected behavior

### Test Structure
- Use table-driven tests for multiple scenarios
- Organize tests with `t.Run()` for subtests
- Test naming: `TestFunctionName_Scenario` or `should describe behavior`
- Always test:
  - Happy path (successful operations)
  - Error cases (invalid input, database errors)
  - Edge cases (boundary conditions)
  - Authentication/authorization failures

### Mocking
- Use gomock for interface mocking
- Generate mocks with: `mockgen -source=file.go -destination=file_mock.go`
- Mock external dependencies (database, external services)
- Don't mock what you own (internal structs)

### Test Coverage
- Aim for high coverage on critical paths (handlers, middleware)
- Run tests with: `go test ./... -v`
- Generate coverage: `go test ./... -coverprofile=coverage.out`
- Coverage HTML report: `go tool cover -html=coverage.out`

## Security Best Practices

### Authentication & Authorization
- Use JWT tokens (HS256) for authentication
- Store JWT secret in environment variables, never hardcode
- Validate tokens in middleware before protected routes
- Implement role-based access control (RBAC)
- Support roles: `user`, `admin`

### Data Protection
- Hash passwords with bcrypt (never store plain text)
- Use GORM parameterized queries to prevent SQL injection
- Validate all user input with Gin binding tags
- Never return password hashes in API responses (use `json:"-"`)
- Never log sensitive information

### Input Validation
- Use Gin validator tags: `binding:"required,email,min=6"` etc.
- Validate role values: `binding:"oneof=user admin"`
- Return `400 Bad Request` for invalid input
- Sanitize inputs before database operations

## Database Guidelines

### GORM Best Practices
- Use GORM model tags for schema definition
- Always use context in repository methods
- Use `AutoMigrate()` for schema migrations
- Enable connection pooling (default in GORM)
- Use transactions for multi-step operations

### Model Definition
- Embed `gorm.Model` or define your own base fields
- Use JSON tags to control API responses: `json:"field"` or `json:"-"`
- Use GORM tags for constraints: `gorm:"not null;uniqueIndex"`
- Define default values: `gorm:"default:'user'"`

### Repository Pattern
- Define repository interface in `internal/repository/`
- Implement PostgreSQL version and mock version
- All methods accept `context.Context` as first parameter
- Return domain models, not database-specific types

## HTTP Handlers

### Handler Structure
- Create handler struct with dependencies (repository, logger)
- Use constructor functions: `NewUserHandler(repo UserRepository)`
- Define request/response structs with binding tags
- Keep handlers thin - delegate business logic to services/repositories

### Request/Response
- Parse JSON with Gin binding: `c.ShouldBindJSON(&req)`
- Extract path params: `c.Param("id")`
- Get authenticated user from context: `c.Get("user_id")`
- Return JSON responses: `c.JSON(status, data)`

### Middleware
- Place middleware in `internal/middleware/`
- Middleware should:
  - Call `c.Next()` to continue the chain
  - Call `c.AbortWithStatusJSON()` to stop processing
- Order matters: logging → auth → business logic

## Environment Configuration

### Required Environment Variables
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT signing
- `SERVER_PORT`: Server port (default: 8080)

### Configuration Loading
- Load from environment variables in `pkg/config/`
- Provide sensible defaults where appropriate
- Fail fast if critical config is missing
- Never commit secrets to version control

## Docker & Kubernetes

### Docker Best Practices
- Use multi-stage builds (builder + runtime)
- Run as non-root user for security
- Use Alpine Linux for minimal image size
- Include health checks
- Set proper resource limits

### Kubernetes/Helm
- Define resources in `helm/myapp/templates/`
- Use values.yaml for configuration
- Store secrets in Kubernetes Secrets
- Configure resource requests/limits
- Enable horizontal pod autoscaling (HPA)
- Include liveness and readiness probes

## CI/CD

### GitHub Actions Workflows
- **Build**: Compile application, upload artifact
- **Test**: Run tests with coverage, upload to Codecov
- **Deploy**: Build Docker image, deploy to Kubernetes

### Workflow Triggers
- Push to `main` or `develop` branches
- Pull requests to `main`
- Version tags (for releases)

## Observability

### Prometheus Metrics
- Expose metrics at `/metrics` endpoint
- Track HTTP requests: `http_requests_total{method,path,status}`
- Track duration: `http_request_duration_seconds{method,path}`
- Use Prometheus client library v1.19.1

### Health Checks
- Implement `/health` endpoint
- Return `200 OK` with `{"status": "healthy"}`
- Include in Docker HEALTHCHECK and K8s probes

## API Design

### RESTful Conventions
- Use proper HTTP methods: GET (read), POST (create), PUT (update), DELETE (delete)
- Use plural nouns for collections: `/users`
- Use ID in path for specific resources: `/users/{id}`
- Return created resource with `201 Created` status

### OpenAPI/Documentation
- Document all endpoints in README.md
- Include request/response examples
- Specify authentication requirements
- Document error responses

## Development Workflow

### Getting Started
```bash
# Install dependencies
go mod download

# Run tests
go test ./... -v

# Build application
go build ./cmd/server

# Run locally
export DATABASE_URL="postgres://..."
export JWT_SECRET="secret"
go run ./cmd/server
```

### Common Commands
- `go test ./... -v`: Run all tests
- `go test ./... -race`: Run with race detection
- `go test ./... -coverprofile=coverage.out`: Generate coverage
- `go mod tidy`: Clean up dependencies
- `docker build -t myapp .`: Build Docker image
- `helm install myapp ./helm/myapp`: Deploy to K8s

## Code Review Checklist

- [ ] Tests written and passing (TDD approach)
- [ ] Error handling implemented correctly
- [ ] Logging added with appropriate context
- [ ] Security best practices followed
- [ ] Input validation included
- [ ] Documentation updated (if needed)
- [ ] No sensitive data in logs or responses
- [ ] Code follows Go conventions
- [ ] Dependencies are justified and minimal

## Common Patterns

### Creating a New Handler
1. Define handler struct with dependencies
2. Write tests first (TDD)
3. Create request/response structs with validation tags
4. Implement handler methods
5. Register routes in `internal/routes/`

### Adding Middleware
1. Create middleware function returning `gin.HandlerFunc`
2. Write tests for middleware behavior
3. Use `c.Next()` to continue or `c.AbortWithStatusJSON()` to stop
4. Apply in route setup

### Adding a New Model
1. Define struct with GORM and JSON tags
2. Add to `AutoMigrate()` in main.go
3. Create repository interface
4. Implement repository methods
5. Write tests for repository

## Performance Considerations

- Use database connection pooling
- Minimize database queries (avoid N+1)
- Use indexes on frequently queried fields
- Set appropriate timeouts for external calls
- Configure resource limits in Kubernetes
- Enable HPA for automatic scaling

## Maintenance

### Dependency Updates
- Review and update dependencies regularly
- Test thoroughly after updates
- Pin exact versions for reproducibility
- Use `go mod tidy` to clean up

### Monitoring in Production
- Monitor Prometheus metrics
- Check logs for errors
- Set up alerts for critical issues
- Monitor resource usage (CPU, memory)

---

Remember: Write tests first, handle errors explicitly, log with context, and never commit secrets!
