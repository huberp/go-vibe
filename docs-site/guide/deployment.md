# Deployment

go-vibe ships with everything needed for production deployment: a multi-stage Dockerfile, a docker-compose file for local development, and a Helm chart for Kubernetes.

## Docker

### Build the Image

The `Dockerfile` uses a two-stage build — a `builder` stage compiles the binary and a minimal `alpine` runtime stage runs it. The final image is under 20 MB.

```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# Stage 2: Runtime
FROM alpine:3.19
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /app/server .
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1
ENTRYPOINT ["./server"]
```

```bash
# Build
docker build -t go-vibe:latest .

# Run (supply your own values)
docker run -d \
  --name go-vibe \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=disable" \
  -e JWT_SECRET="your-secret" \
  go-vibe:latest
```

### Multi-Architecture Build

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t ghcr.io/huberp/go-vibe:latest \
  --push .
```

## Docker Compose

The `docker-compose.yml` at the repo root starts the full stack:

```yaml
version: "3.9"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://govibe:govibe@db:5432/govibe?sslmode=disable
      JWT_SECRET: dev-secret-change-in-production
      SERVER_PORT: "8080"
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: govibe
      POSTGRES_PASSWORD: govibe
      POSTGRES_DB: govibe
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U govibe"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  pgdata:
```

```bash
# Start (detached)
docker compose up -d --build

# View logs
docker compose logs -f app

# Stop and remove volumes
docker compose down -v
```

## Kubernetes with Helm

The Helm chart lives in `helm/myapp/`. It provides a production-grade deployment with HPA, secrets, and health probes.

### Quick Deploy

```bash
# Add your secrets (never commit these values!)
kubectl create secret generic go-vibe-secrets \
  --from-literal=jwt-secret="$(openssl rand -hex 32)" \
  --from-literal=database-url="postgres://user:pass@pg-svc:5432/govibe?sslmode=disable"

# Install the Helm chart
helm install go-vibe ./helm/myapp \
  --set image.tag=latest \
  --set replicaCount=2
```

### Helm Values Reference

Key values in `helm/myapp/values.yaml`:

```yaml
replicaCount: 2

image:
  repository: ghcr.io/huberp/go-vibe
  pullPolicy: IfNotPresent
  tag: "latest"

service:
  type: ClusterIP
  port: 80
  targetPort: 8080

ingress:
  enabled: true
  className: nginx
  hosts:
    - host: api.example.com
      paths:
        - path: /
          pathType: Prefix

resources:
  limits:
    cpu: 500m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 64Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80

env:
  SERVER_PORT: "8080"

# Reference Kubernetes secrets for sensitive values
secretEnv:
  DATABASE_URL:
    secretName: go-vibe-secrets
    secretKey: database-url
  JWT_SECRET:
    secretName: go-vibe-secrets
    secretKey: jwt-secret
```

### Health Probes

The Helm chart configures Kubernetes liveness and readiness probes against `/health`:

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 15
  periodSeconds: 20
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  failureThreshold: 3
```

### Horizontal Pod Autoscaler

HPA scales the deployment automatically when CPU or memory exceed configured thresholds:

```bash
# Check HPA status
kubectl get hpa go-vibe

# NAME       REFERENCE               TARGETS         MINPODS  MAXPODS  REPLICAS
# go-vibe    Deployment/go-vibe      22%/70%         2        10       2
```

### Helm Upgrade & Rollback

```bash
# Upgrade to a new image tag
helm upgrade go-vibe ./helm/myapp --set image.tag=v1.2.0

# Roll back to the previous release
helm rollback go-vibe 1

# View release history
helm history go-vibe
```

## Database Migrations

go-vibe uses GORM `AutoMigrate` on startup. For production, consider using [golang-migrate](https://github.com/golang-migrate/migrate) for controlled, versioned migrations.

SQL migration files are stored in `migrations/`:

```
migrations/
  001_create_users_table.up.sql
  001_create_users_table.down.sql
```

```bash
# Run migrations manually
migrate -path ./migrations -database "$DATABASE_URL" up
```

## Production Checklist

- [ ] `JWT_SECRET` is at least 32 random bytes — stored in Kubernetes Secret
- [ ] `DATABASE_URL` uses SSL (`sslmode=require` or `sslmode=verify-full`)
- [ ] Resource requests/limits configured in Helm values
- [ ] HPA enabled with sensible thresholds
- [ ] Ingress configured with TLS termination
- [ ] Liveness and readiness probes verified
- [ ] Prometheus scraping configured (ServiceMonitor or annotation-based)
- [ ] Log aggregation set up (Loki, Datadog, CloudWatch, etc.)
