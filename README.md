# Production-Ready User Management Microservice

A production-ready microservice built with Go 1.25.2, Gin v1.11.0, following TDD principles and designed for Kubernetes deployment.

## Features

- ✅ RESTful API for user management (CRUD operations)
- ✅ JWT-based authentication and authorization
- ✅ Role-based access control (admin/user)
- ✅ PostgreSQL database with GORM
- ✅ Structured logging with Zap
- ✅ W3C trace context support for distributed tracing
- ✅ OpenTelemetry (OTEL) tracing integration
- ✅ Prometheus metrics (including user count)
- ✅ Rate limiting middleware (100 req/s per IP)
- ✅ CORS middleware for cross-origin requests
- ✅ API versioning (/v1/...) for backward compatibility
- ✅ Enhanced bcrypt security (cost factor: 12)
- ✅ 100% test coverage for handlers and middleware
- ✅ Multi-stage Docker build (Alpine-based)
- ✅ Helm chart for Kubernetes deployment
- ✅ GitHub Actions CI/CD pipelines

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
│   └── utils/                  # Utilities (JWT, hashing)
├── helm/                       # Helm chart
├── .github/workflows/          # CI/CD pipelines
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
| JWT | v5.3.0 | Authentication |
| Viper | v1.21.0 | Configuration management |
| Zap | v1.27.0 | Structured logging |
| Testify | v1.11.1 | Testing framework |
| Prometheus | v1.23.2 | Metrics |
| OpenTelemetry | v1.33.0 | Distributed tracing |
| CORS | v1.7.0 | Cross-origin resource sharing |
| Rate Limiter | golang.org/x/time | Request rate limiting |

## Prerequisites

- Go 1.25.2+
- PostgreSQL 13+
- Docker (for containerization)
- Kubernetes cluster (for deployment)
- Helm 3+ (for deployment)

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

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/myapp?sslmode=disable"
export JWT_SECRET="your-secret-key-change-in-production"
export SERVER_PORT="8080"
```

### 4. Run the application

```bash
go run ./cmd/server
```

The server will start on `http://localhost:8080`

## API Documentation

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
      summary: Health check
      responses:
        '200':
          description: Service is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: healthy

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

### 7. Health check

```bash
curl http://localhost:8080/health
```

### 8. Prometheus metrics (includes users_total)

```bash
curl http://localhost:8080/metrics | grep -E "(http_requests_total|users_total)"
```

Example output:
```
http_requests_total{method="GET",path="/v1/users",status="200"} 42
users_total 156
```

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

```bash
helm install myapp ./helm/myapp \
  --namespace production \
  --create-namespace
```

3. **Update deployment:**

```bash
helm upgrade myapp ./helm/myapp \
  --set image.tag=v1.0.1 \
  -n production
```

4. **Uninstall:**

```bash
helm uninstall myapp -n production
```

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

**Troubleshooting**: See [Local Kubernetes Troubleshooting Guide](docs/LOCAL_K8S_TROUBLESHOOTING.md) for common issues and solutions.

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

The project includes three GitHub Actions workflows:

### 1. Build (`build.yml`)
- Triggers on push to main/develop or PR to main
- Builds the application
- Uploads binary artifact

### 2. Test (`test.yml`)
- Triggers on push to main/develop or PR to main
- Runs all tests with race detection
- Generates coverage report
- Uploads to Codecov

### 3. Deploy (`deploy.yml`)
- Triggers on push to main or version tags
- Builds and pushes Docker image to GHCR
- Deploys to Kubernetes using Helm

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
- Health checks with liveness/readiness probes
- Multi-stage Docker build for minimal image size

## License

MIT

## Support

For issues and questions, please open an issue on GitHub.
