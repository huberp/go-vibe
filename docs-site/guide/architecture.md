# Architecture

go-vibe follows Clean Architecture principles, keeping concerns separated and dependencies pointing inward. The result is a codebase that is easy to test, extend, and reason about.

## Design Principles

- **Dependency Inversion** — high-level modules (handlers) depend on interfaces, not concrete implementations
- **Single Responsibility** — each package has one clear job
- **Dependency Injection** — all dependencies are injected at construction time, never grabbed globally
- **Test-Driven** — interfaces make every layer independently testable with mocks

## Request Flow

```mermaid
graph LR
    Client -->|HTTP Request| Middleware
    Middleware -->|Auth / Logging / Rate Limit| Handler
    Handler -->|Call interface method| Repository
    Repository -->|Parameterised SQL via GORM| PostgreSQL

    style Client fill:#4A90D9,color:#fff
    style Middleware fill:#F5A623,color:#fff
    style Handler fill:#7ED321,color:#fff
    style Repository fill:#9B59B6,color:#fff
    style PostgreSQL fill:#336791,color:#fff
```

### Detailed Flow

1. **Client** sends an HTTP request to the Gin router
2. **Middleware chain** processes the request in order:
   - `RequestLogger` — logs method, path, client IP
   - `MetricsMiddleware` — starts the duration timer
   - `RateLimiter` — checks the token bucket for the client IP
   - `CORSMiddleware` — adds CORS headers
   - `AuthMiddleware` — validates JWT, injects `user_id` and `role` into context (protected routes only)
3. **Handler** parses and validates the request body, calls the repository, and writes the JSON response
4. **Repository** executes the database query via GORM, returns domain models
5. **PostgreSQL** executes the query and returns rows
6. The response travels back up the chain — `MetricsMiddleware` records the final status code and duration

## Project Structure

```
go-vibe/
├── cmd/
│   └── server/
│       └── main.go           # Entry point — wires everything together
├── internal/
│   ├── handlers/
│   │   ├── user_handler.go   # HTTP handlers for /v1/users
│   │   └── user_handler_test.go
│   ├── middleware/
│   │   ├── auth.go           # JWT validation middleware
│   │   ├── cors.go           # CORS headers middleware
│   │   ├── logger.go         # Zap request logging middleware
│   │   ├── metrics.go        # Prometheus instrumentation
│   │   └── rate_limit.go     # Token-bucket rate limiter
│   ├── models/
│   │   └── user.go           # User domain model (GORM + JSON tags)
│   ├── repository/
│   │   ├── user_repository.go          # UserRepository interface
│   │   ├── postgres_user_repository.go # PostgreSQL implementation
│   │   └── user_repository_mock.go     # gomock-generated mock
│   └── routes/
│       └── routes.go         # Route registration
├── pkg/
│   ├── config/
│   │   └── config.go         # Layered YAML + env-var config loading (Viper)
│   ├── logger/
│   │   └── logger.go         # Zap logger initialisation
│   └── utils/
│       └── password.go       # bcrypt helpers
├── config/
│   ├── base.yaml             # Shared defaults for all stages
│   ├── development.yaml      # Overrides for local development
│   ├── staging.yaml          # Overrides for staging
│   └── production.yaml       # Overrides for production
├── migrations/               # SQL migration files
├── helm/myapp/               # Helm chart for Kubernetes
├── docs-site/                # This VitePress documentation
├── Dockerfile                # Multi-stage container build
├── docker-compose.yml        # Local development stack
└── go.mod
```

## Layer Responsibilities

### `cmd/server/main.go`

The composition root. Responsibilities:
- Load config from environment
- Initialise the logger
- Open the database connection and run `AutoMigrate`
- Construct repositories, handlers, and middleware
- Register routes and start the HTTP server

```go
func main() {
    cfg    := config.Load()
    logger := logger.New(cfg.Env)
    db     := database.Connect(cfg.DatabaseURL)
    repo   := repository.NewPostgresUserRepository(db)
    h      := handlers.NewUserHandler(repo, logger)

    r := routes.Setup(h, cfg, logger)
    r.Run(":" + cfg.ServerPort)
}
```

### `internal/handlers/`

Thin HTTP adapters. Each handler:
1. Parses and validates the request (Gin binding)
2. Calls one repository method
3. Writes the JSON response with the correct status code

Handlers never contain business logic or SQL — they only translate between HTTP and the domain layer.

### `internal/middleware/`

Reusable Gin middleware functions. Each middleware calls `c.Next()` to continue the chain or `c.AbortWithStatusJSON()` to short-circuit with an error response.

### `internal/repository/`

The data access layer, hidden behind an interface:

```go
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByID(ctx context.Context, id uint) (*models.User, error)
    FindByEmail(ctx context.Context, email string) (*models.User, error)
    FindAll(ctx context.Context) ([]models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uint) error
}
```

The interface means:
- **Tests** use a `MockUserRepository` — no real database needed
- **Future implementations** (Redis cache, alternate database) require zero changes to handlers

### `internal/models/`

Plain Go structs with GORM and JSON tags. No methods, no business logic — just data shapes.

### `pkg/`

Reusable, domain-agnostic packages that could be extracted into separate modules:
- `config` — layered YAML config loading with stage-specific overrides and env-var bindings (Viper)
- `logger` — Zap logger factory
- `utils` — bcrypt password helpers

## Dependency Graph

```
main.go
  ├── pkg/config        (no internal deps)
  ├── pkg/logger        (no internal deps)
  ├── internal/models   (no internal deps)
  ├── internal/repository
  │     └── internal/models
  ├── internal/handlers
  │     ├── internal/repository  (interface only)
  │     └── internal/models
  ├── internal/middleware
  │     └── pkg/logger
  └── internal/routes
        ├── internal/handlers
        └── internal/middleware
```

Dependencies always flow **inward** toward `models` — never outward. This is the core invariant of Clean Architecture.

## Testing Strategy

| Layer | Test Type | Database? | Key Tool |
|-------|-----------|-----------|----------|
| Handlers | Unit | ❌ (mock repo) | gomock |
| Middleware | Unit | ❌ | httptest |
| Repository | Integration | ✅ (real PostgreSQL) | testify |
| Routes | Integration | ❌ | httptest + mock repo |

Run the full suite with:

```bash
go test ./... -race -v
```
