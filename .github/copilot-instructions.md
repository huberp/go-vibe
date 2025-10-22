# Copilot Instructions for go-vibe

> **📘 Note**: This file provides comprehensive development guidelines for GitHub Copilot and other AI coding assistants, as well as human contributors. It ensures consistent code quality and adherence to project standards.

This document provides overarching guidelines and best practices for contributing to the go-vibe project. These instructions apply across the entire codebase.

## Table of Contents

- [Project Overview](#project-overview)
- [Copilot Communication Style](#copilot-communication-style)
- [Tech Stack](#tech-stack)
- [Architecture Patterns](#architecture-patterns)
- [Code Style and Standards](#code-style-and-standards)
- [Testing Strategy](#testing-strategy)
- [Security Best Practices](#security-best-practices)
- [Database Guidelines](#database-guidelines)
- [HTTP Handlers](#http-handlers)
- [Environment Configuration](#environment-configuration)
- [Docker & Kubernetes](#docker--kubernetes)
- [CI/CD](#cicd)
- [Observability](#observability)
- [API Design](#api-design)
- [Development Workflow](#development-workflow)
- [Code Review Checklist](#code-review-checklist)
- [Common Patterns](#common-patterns)
- [Performance Considerations](#performance-considerations)
- [Communication and Contribution](#communication-and-contribution)
- [Do's and Don'ts](#dos-and-donts)
- [Dependency Management](#dependency-management)
- [External Resources](#external-resources)
- [Maintenance](#maintenance)

## Project Overview

go-vibe is a production-ready user management microservice built with Go 1.25.2, following Test-Driven Development (TDD) principles and designed for cloud-native Kubernetes deployment.

**Target Audience**: Backend developers working on Go microservices with PostgreSQL databases, deployed on Kubernetes.

**Key Design Decisions**:
- TDD-first approach for all new features
- Repository pattern for data abstraction
- JWT-based stateless authentication
- Horizontal scalability via Kubernetes HPA
- Observability-first design with structured logging and metrics

## Copilot Communication Style

These guidelines define how Copilot communicates in code, comments, commit messages, and documentation to ensure consistency, clarity, and professionalism across the project.

### Be Concise
- Keep responses short, clear, and direct
- Avoid unnecessary filler words or overly detailed explanations
- Focus on delivering essential information

**Example**:
- ❌ "This function is responsible for validating the user's input to ensure it meets the required criteria."
- ✅ "Copilot has created a function to validate user input."

### Use Active Voice
- Prefer active voice to make communication more direct and engaging
- Use phrases like "Copilot has generated," "Copilot recommends," or "Copilot suggests"

**Examples**:
- ❌ "A configuration file was generated."
- ✅ "Copilot has generated a configuration file."

- ❌ "This code snippet is provided to handle authentication."
- ✅ "Copilot has provided this code snippet to handle authentication."

### Neutral Tone with Explicit Self-Reference
- Always refer to itself as "Copilot"
- Avoid referring to itself as "I," "we," or "an assistant"
- Use neutral, professional language

**Examples**:
- ❌ "I suggest using this function to handle user authentication."
- ✅ "Copilot suggests using this function to handle user authentication."

- ❌ "This is the code I generated for you."
- ✅ "This is the code Copilot has generated."

### Clarity Over Length
- Prioritize clarity and precision over lengthy explanations
- Use bulleted lists, tables, or headings to organize information when needed

**Examples**:
- ❌ "This API endpoint is designed to handle user authentication and authorization, ensuring that only valid users can access the application."
- ✅ "Copilot has designed this API endpoint to handle user authentication and authorization."

### Avoid Redundancy
- Eliminate repetitive phrasing and redundant information
- Focus on the most critical details

**Example**:
- ❌ "Copilot has created this method to handle API requests, and this method is responsible for processing data from the client API requests."
- ✅ "Copilot has created this method to process API requests."

### Code Comments Style
- Keep inline comments short and focused on the intent of the code
- Avoid long, detailed comments unless necessary for complex logic
- Use Copilot self-reference in comments when appropriate

**Example**:
```go
// Copilot has added this handler to log errors.
func logError(err error) {
    fmt.Println(err)
}
```

### Commit Messages with Copilot Reference
- Use active voice and refer explicitly to Copilot in commit messages
- Follow [Conventional Commits](https://www.conventionalcommits.org/) format
- Use imperative mood with Copilot attribution

**Examples**:
- ✅ `feat(auth): Copilot has added Google login support`
- ✅ `fix(routes): Copilot has resolved a routing error`
- ✅ `refactor(handlers): Copilot has simplified error handling logic`

### Tools for Enforcing Guidelines
- Use tools like **Vale** or **Grammarly** to validate clarity and conciseness
- Regularly review generated content for adherence to these guidelines
- Apply linters and formatters consistently

## Tech Stack

- **Language**: Go 1.25.2+
- **Web Framework**: Gin v1.11.0
- **ORM**: GORM v1.31.0
- **Database**: PostgreSQL 13+
- **Authentication**: JWT (github.com/golang-jwt/jwt/v5 v5.3.0)
- **Logging**: Zap v1.27.0
- **Testing**: Testify v1.11.1 + gomock
- **Metrics**: Prometheus v1.23.2
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

### Naming Conventions

**Variables**:
- Use camelCase for local variables: `userID`, `requestCount`
- Use descriptive names, avoid single letters except for: `i`, `j` (loops), `c` (Gin context), `t` (tests)
- Boolean variables start with `is`, `has`, `should`: `isValid`, `hasPermission`

**Functions**:
- Use PascalCase for exported functions: `CreateUser`, `ValidateToken`
- Use camelCase for private functions: `hashPassword`, `parseRequest`
- Function names should be verbs or verb phrases: `GetUser`, `UpdateRecord`, `CalculateTotal`

**Files**:
- Use snake_case for Go files: `user_handler.go`, `auth_middleware.go`
- Test files: `*_test.go` (e.g., `user_handler_test.go`)
- Mock files: `*_mock.go` (e.g., `user_repository_mock.go`)

**Packages**:
- Use short, lowercase, single-word names: `handlers`, `models`, `utils`
- Avoid underscores or mixed caps
- Package name should match directory name

**Constants**:
- Use PascalCase for exported: `DefaultTimeout`, `MaxRetries`
- Use camelCase or UPPER_CASE for private based on context

### Code Comments

**When to Comment**:
- ✅ Package documentation (package-level comment)
- ✅ Exported functions, types, and constants (godoc format)
- ✅ Complex business logic or algorithms
- ✅ TODOs with ticket references: `// TODO(#123): Implement retry logic`
- ❌ Don't comment obvious code
- ❌ Don't leave commented-out code

**Godoc Format**:
```go
// Package handlers provides HTTP request handlers for the user management API.
package handlers

// CreateUser handles user creation requests.
// It validates input, hashes the password, and stores the user in the database.
// Returns 201 on success, 400 for invalid input, or 500 for server errors.
func (h *UserHandler) CreateUser(c *gin.Context) {
    // implementation
}
```

**Documentation Style**:
- Write in complete sentences
- Start with the name of the thing being described
- Be concise but complete
- Include error conditions and return values for functions

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
- Use Prometheus client library v1.23.2

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

## Communication and Contribution

### Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or modifications
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

**Examples**:
```
feat(auth): Copilot has added JWT token refresh endpoint

Copilot implements token refresh mechanism for expired tokens.
Closes #42

fix(handlers): Copilot has validated user ID before database query

Copilot prevents panic when invalid ID format is provided.
```

**Rules**:
- Keep subject line under 72 characters
- Use imperative mood with Copilot attribution: "Copilot has added" or "Copilot has fixed"
- No period at the end of subject line
- Reference issue numbers in footer
- Always use active voice with explicit Copilot self-reference

### Pull Request Guidelines

**Branch Naming**:
- Feature: `feature/<issue-number>-short-description` (e.g., `feature/42-add-user-export`)
- Bug fix: `fix/<issue-number>-short-description` (e.g., `fix/43-auth-token-validation`)
- Hotfix: `hotfix/<issue-number>-short-description`

**PR Title Format**:
```
<type>: <description>
```

Examples:
- `feat: Copilot has added user export functionality`
- `fix: Copilot has resolved JWT token validation issue`

**PR Description Template**:
```markdown
## Overview
Brief description of changes

## Changes Made
- List of changes
- Another change

## Testing
How changes were tested

## Related Issues
Closes #<issue-number>
```

**Review Requirements**:
- ✅ All tests passing
- ✅ Code coverage maintained or improved
- ✅ No security vulnerabilities introduced
- ✅ Documentation updated if needed
- ✅ At least one approval from maintainer

### Code Review Checklist

When reviewing PRs, verify:

- [ ] Tests written and passing (TDD approach)
- [ ] Error handling implemented correctly
- [ ] Logging added with appropriate context
- [ ] Security best practices followed
- [ ] Input validation included
- [ ] Documentation updated (if needed)
- [ ] No sensitive data in logs or responses
- [ ] Code follows Go conventions
- [ ] Dependencies are justified and minimal
- [ ] Commit messages follow conventions
- [ ] No commented-out code

## Do's and Don'ts

### ✅ DO:

- **DO** write tests before implementation (TDD)
- **DO** use dependency injection for testability
- **DO** validate all user inputs
- **DO** hash passwords with bcrypt
- **DO** use structured logging with context
- **DO** handle errors explicitly at every level
- **DO** use GORM parameterized queries
- **DO** return appropriate HTTP status codes
- **DO** document exported functions with godoc
- **DO** use context.Context for cancellation and timeouts
- **DO** run `go mod tidy` after dependency changes
- **DO** use interfaces for repository layer
- **DO** store secrets in environment variables
- **DO** close database connections and resources properly

### ❌ DON'T:

- **DON'T** use panic for error handling
- **DON'T** commit secrets or credentials
- **DON'T** log sensitive data (passwords, tokens, PII)
- **DON'T** return password hashes in API responses
- **DON'T** use global variables for state
- **DON'T** mock what you own (internal structs)
- **DON'T** skip input validation
- **DON'T** use string concatenation for SQL queries
- **DON'T** hardcode configuration values
- **DON'T** ignore errors (use `if err != nil`)
- **DON'T** use `panic()` or `os.Exit()` in library code
- **DON'T** modify request/response after calling `c.Next()`
- **DON'T** use time.Sleep in production code (use proper timeouts)

### Common Pitfalls to Avoid

1. **N+1 Query Problem**: Use eager loading with GORM's `Preload()`
2. **Goroutine Leaks**: Always ensure goroutines can exit
3. **Missing Context**: Always propagate context through call chain
4. **Race Conditions**: Run tests with `-race` flag
5. **Improper Error Wrapping**: Use `fmt.Errorf("context: %w", err)`
6. **Missing Middleware Order**: Auth before business logic
7. **Unbounded Slices**: Set capacity when size is known
8. **Pointer to Loop Variable**: Use local variable in loop closures

## Dependency Management

### Adding New Dependencies

1. **Before Adding**:
   - Check if functionality exists in standard library
   - Verify the package is actively maintained
   - Review security advisories
   - Consider the dependency tree size

2. **How to Add**:
```bash
# Add dependency
go get github.com/package/name@version

# Update go.mod and go.sum
go mod tidy

# Verify
go mod verify
```

3. **Approval Process**:
   - Discuss in issue or PR why dependency is needed
   - Get maintainer approval for new dependencies
   - Document the purpose in PR description

### Preferred Libraries

**For Common Tasks**:
- **HTTP Router**: Gin (already in use)
- **ORM**: GORM (already in use)
- **Testing**: testify + gomock (already in use)
- **Logging**: zap (already in use)
- **UUID**: github.com/google/uuid (already in use)
- **JWT**: github.com/golang-jwt/jwt/v5 (already in use)
- **Validation**: Use Gin's built-in validator
- **Environment**: Viper for configuration management (preferred) or `os.Getenv()`

**Avoid**:
- Unnecessary web frameworks (stick with Gin)
- Multiple logging libraries (use zap)
- Alternative ORMs (use GORM)

## External Resources

### Official Documentation

- **Go Language**: https://go.dev/doc/
- **Gin Framework**: https://gin-gonic.com/docs/
- **GORM**: https://gorm.io/docs/
- **PostgreSQL**: https://www.postgresql.org/docs/
- **JWT**: https://jwt.io/introduction
- **Zap Logging**: https://pkg.go.dev/go.uber.org/zap
- **Testify**: https://pkg.go.dev/github.com/stretchr/testify
- **Prometheus**: https://prometheus.io/docs/

### API References

- **Gin Context**: https://pkg.go.dev/github.com/gin-gonic/gin#Context
- **GORM Models**: https://gorm.io/docs/models.html
- **GORM Associations**: https://gorm.io/docs/associations.html
- **bcrypt**: https://pkg.go.dev/golang.org/x/crypto/bcrypt

### Internal Documentation

- **README.md**: Project setup and API documentation
- **IMPLEMENTATION_SUMMARY.md**: Architecture and implementation details
- **OpenAPI Spec**: In README.md, defines all API endpoints
- **Helm Chart**: `helm/myapp/` for Kubernetes deployment

### Related Standards

- **Effective Go**: https://go.dev/doc/effective_go
- **Go Code Review Comments**: https://github.com/golang/go/wiki/CodeReviewComments
- **Uber Go Style Guide**: https://github.com/uber-go/guide/blob/master/style.md
- **12-Factor App**: https://12factor.net/ (for cloud-native principles)

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

**Note**: This file should be located at `.github/copilot-instructions.md` for optimal GitHub Copilot workspace agent integration.

Remember: Write tests first, handle errors explicitly, log with context, and never commit secrets!
