# Configuration

go-vibe uses a **layered configuration** strategy: shared defaults live in `config/base.yaml`, stage-specific files override them, and environment variables override everything. This keeps common settings in one place while letting each environment stay independently tunable without duplicating values.

## How It Works

Configuration is loaded at startup by `pkg/config/config.go` using [Viper](https://github.com/spf13/viper):

1. `config/base.yaml` is read first — it provides defaults for all keys.
2. The stage file (`config/<stage>.yaml`) is merged on top, overriding only the keys it defines.
3. Environment variables are applied last and always win.

```
┌─────────────────────────────────────────────────────┐
│              Environment Variables                   │  ← highest priority
│  DATABASE_URL, JWT_SECRET, SERVER_PORT, …            │
└──────────────────────┬──────────────────────────────┘
                       │ overrides
┌──────────────────────▼──────────────────────────────┐
│          config/<stage>.yaml                         │
│  development.yaml / staging.yaml / production.yaml   │
└──────────────────────┬──────────────────────────────┘
                       │ overrides
┌──────────────────────▼──────────────────────────────┐
│              config/base.yaml                        │  ← lowest priority
│  Shared defaults for all environments                │
└─────────────────────────────────────────────────────┘
```

### Selecting the Stage

Set the `APP_STAGE` environment variable before starting the server (defaults to `development`):

```bash
export APP_STAGE="production"   # loads config/production.yaml on top of base
```

Supported stages out of the box: `development`, `staging`, `production`.

## Config Files

### `config/base.yaml`

Shared defaults applied to every stage:

```yaml
# config/base.yaml
server:
  port: "8080"

database:
  url: "postgres://user:password@localhost:5432/myapp?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: 30  # seconds

jwt:
  secret: "your-secret-key"

rate_limit:
  requests_per_second: 100
  burst: 200

observability:
  otel: false
```

### `config/development.yaml`

Overrides for local development — only keys that differ from `base.yaml` need to appear here:

```yaml
# config/development.yaml
database:
  url: "postgres://myapp:myapp@localhost:5432/myapp?sslmode=disable"

jwt:
  secret: "dev-secret-key"

observability:
  otel: true
```

### `config/staging.yaml`

Overrides for staging — tighter rate limits and secrets from environment variables:

```yaml
# config/staging.yaml
database:
  url: "postgres://user:password@staging-db:5432/myapp?sslmode=disable"
  max_open_conns: 50

jwt:
  secret: "${JWT_SECRET}"   # injected via environment variable

rate_limit:
  requests_per_second: 50
  burst: 100

observability:
  otel: true
```

### `config/production.yaml`

Production overrides — all sensitive values come from environment variables:

```yaml
# config/production.yaml
database:
  url: "${DATABASE_URL}"
  max_open_conns: 100
  max_idle_conns: 25
  conn_max_lifetime: 60

jwt:
  secret: "${JWT_SECRET}"

server:
  port: "${SERVER_PORT:8080}"   # default to 8080 if not set

rate_limit:
  requests_per_second: 50
  burst: 100

observability:
  otel: true
```

## Environment Variables

Environment variables always override YAML values and are the recommended way to inject secrets. Viper maps them to nested keys using `_` as a separator:

| Environment Variable | Config Key | Description |
|---|---|---|
| `APP_STAGE` | — | Active stage (`development`, `staging`, `production`) |
| `SERVER_PORT` | `server.port` | HTTP server port |
| `DATABASE_URL` | `database.url` | PostgreSQL connection string |
| `DB_MAX_OPEN_CONNS` | `database.max_open_conns` | Max open DB connections |
| `DB_MAX_IDLE_CONNS` | `database.max_idle_conns` | Max idle DB connections |
| `DB_CONN_MAX_LIFETIME` | `database.conn_max_lifetime` | Connection lifetime in seconds |
| `JWT_SECRET` | `jwt.secret` | JWT signing secret |
| `RATE_LIMIT_REQUESTS_PER_SECOND` | `rate_limit.requests_per_second` | Allowed requests per second per IP |
| `RATE_LIMIT_BURST` | `rate_limit.burst` | Burst size for the token-bucket limiter |
| `OBSERVABILITY_OTEL` | `observability.otel` | Enable OpenTelemetry (`true`/`false`) |

::: warning Security
Never commit `JWT_SECRET` or `DATABASE_URL` to source control. Use Kubernetes Secrets or a secrets manager in production.
:::

## Adding a New Stage

1. Create `config/<stage>.yaml` with only the keys that differ from `base.yaml`.
2. Set `APP_STAGE=<stage>` in the target environment.

No code changes are required — the loader picks up the new file automatically.

## Configuration in Code

The `pkg/config` package exposes a typed `Config` struct. Access it anywhere by calling `config.Load()`:

```go
// pkg/config/config.go
type Config struct {
    Server        ServerConfig
    Database      DatabaseConfig
    JWT           JWTConfig
    RateLimit     RateLimitConfig
    Observability ObservabilityConfig
}

// Load reads base.yaml, merges the stage file, then applies env-var overrides.
func Load() *Config {
    return LoadWithStage(getStage())
}
```

The active stage is determined by:

```go
func getStage() string {
    stage := os.Getenv("APP_STAGE")
    if stage == "" {
        stage = "development"
    }
    return stage
}
```
