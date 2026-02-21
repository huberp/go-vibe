---
layout: home

hero:
  name: "go-vibe"
  text: "Stop writing boilerplate."
  tagline: "Start building features. A production-ready Go microservice template with JWT auth, PostgreSQL, Prometheus metrics, and Kubernetes deployment — all wired up and ready to go."
  image:
    src: /hero-gopher.svg
    alt: go-vibe gopher
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/huberp/go-vibe

features:
  - icon: 🧪
    title: TDD-First Development
    details: Built with Test-Driven Development from day one. Every handler, middleware, and repository method has comprehensive tests using testify and gomock. Write tests before code — it's not optional here.
    link: /guide/features
    linkText: Learn more

  - icon: 🔐
    title: JWT Auth & RBAC
    details: Stateless JWT authentication with role-based access control (user/admin). Middleware validates tokens, extracts claims, and enforces permissions on every protected route.
    link: /guide/features
    linkText: Learn more

  - icon: 🐘
    title: PostgreSQL + GORM
    details: Repository pattern with GORM v2 for type-safe, parameterized queries. AutoMigrate, connection pooling, and a clean interface layer make swapping databases trivial.
    link: /guide/features
    linkText: Learn more

  - icon: 📊
    title: Prometheus Metrics
    details: Built-in /metrics endpoint tracking HTTP request counts, durations, and user totals. Drop in a Grafana dashboard and get production visibility in minutes.
    link: /guide/observability
    linkText: Learn more

  - icon: 🐳
    title: Docker + Kubernetes
    details: Multi-stage Dockerfile, docker-compose for local dev, and a production-ready Helm chart with HPA, resource limits, liveness/readiness probes, and configurable replicas.
    link: /guide/deployment
    linkText: Learn more

  - icon: ⚙️
    title: CI/CD Pipelines
    details: GitHub Actions workflows for build, test, and deploy. Includes cross-platform PostgreSQL testing (Linux + Windows), codecov integration, and automated Docker image publishing.
    link: /guide/ci-cd
    linkText: Learn more

  - icon: 📝
    title: Structured Logging
    details: Uber Zap for zero-allocation, structured JSON logging. Every request logs method, path, status code, duration, and client IP. Never debug blind in production again.
    link: /guide/observability
    linkText: Learn more

  - icon: 🛡️
    title: Rate Limiting & CORS
    details: Token-bucket rate limiting middleware and configurable CORS headers keep your API safe from abuse and ready for browser-based frontends out of the box.
    link: /guide/features
    linkText: Learn more
---

<div class="home-quick-start">

## Quick Start

Clone the repo, set two environment variables, and you have a running user management API in under a minute.

::: code-group

```bash [Clone & Run]
# 1. Clone the repository
git clone https://github.com/huberp/go-vibe.git
cd go-vibe

# 2. Set required environment variables
export DATABASE_URL="postgres://user:password@localhost:5432/govibe?sslmode=disable"
export JWT_SECRET="your-super-secret-key-change-in-production"
export SERVER_PORT="8080"

# 3. Run database migrations and start the server
go run ./cmd/server
```

```bash [Docker Compose]
# Spin up the full stack (app + PostgreSQL) with one command
docker compose up --build

# The API is now available at http://localhost:8080
# Metrics at http://localhost:8080/metrics
# Swagger UI at http://localhost:8080/swagger/index.html
```

```bash [Test the API]
# Create your first user
curl -s -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123","role":"admin"}' \
  | jq .

# Login and grab the JWT token
TOKEN=$(curl -s -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}' \
  | jq -r .token)

# List all users (admin only)
curl -s http://localhost:8080/v1/users \
  -H "Authorization: Bearer $TOKEN" | jq .
```

:::

</div>
