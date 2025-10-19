# Health Check System

This document describes the registrable health check system and how to use it.

## Overview

The health check system is based on a provider pattern similar to the info endpoint. Components can register health checks with specific scopes that determine which endpoints will execute them.

## Scopes

Health checks can be registered in one or more of the following scopes:

- **`base`**: Only appears in `/health` endpoint
- **`startup`**: Appears in `/health/startup`, `/health/readiness`, and `/health`
- **`ready`**: Appears in `/health/readiness` and `/health`
- **`live`**: Appears in `/health/liveness` and `/health`

When `/health` is called, each provider is checked only once, regardless of how many scopes it's registered in.

## Endpoints

- **`/health`**: Returns all registered health checks (each checked only once)
- **`/health/startup`**: Returns only health checks registered with `ScopeStartup`
- **`/health/readiness`**: Returns only health checks registered with `ScopeReady`
- **`/health/liveness`**: Returns only health checks registered with `ScopeLive`

## Creating a Custom Health Check Provider

To create a custom health check, implement the `HealthCheckProvider` interface:

```go
package health

type HealthCheckProvider interface {
    Name() string
    Check() (*CheckResult, error)
    Scopes() []Scope
}
```

### Example: Cache Health Check

```go
package myapp

import (
    "myapp/pkg/health"
    "time"
)

type CacheHealthCheckProvider struct {
    cache  CacheClient
    scopes []Scope
}

func NewCacheHealthCheckProvider(cache CacheClient, scopes ...health.Scope) *CacheHealthCheckProvider {
    // Default to ready scope if none provided
    if len(scopes) == 0 {
        scopes = []health.Scope{health.ScopeReady}
    }
    return &CacheHealthCheckProvider{
        cache:  cache,
        scopes: scopes,
    }
}

func (c *CacheHealthCheckProvider) Name() string {
    return "cache"
}

func (c *CacheHealthCheckProvider) Check() (*health.CheckResult, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := c.cache.Ping(ctx)
    if err != nil {
        return &health.CheckResult{
            Status: health.StatusDown,
            Details: map[string]interface{}{
                "error": err.Error(),
            },
        }, nil
    }

    return &health.CheckResult{
        Status: health.StatusUp,
        Details: map[string]interface{}{
            "connected": true,
        },
    }, nil
}

func (c *CacheHealthCheckProvider) Scopes() []health.Scope {
    return c.scopes
}
```

### Registering the Health Check

In your `routes.go` or main setup:

```go
// Setup health check providers
healthRegistry := health.NewRegistry()

// Database check - appears in startup and readiness
healthRegistry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup, health.ScopeReady))

// Cache check - appears in readiness only
healthRegistry.Register(NewCacheHealthCheckProvider(cache, health.ScopeReady))

// Custom check - appears in all scopes
healthRegistry.Register(NewCustomHealthCheck(health.ScopeStartup, health.ScopeReady, health.ScopeLive))

// Base check - only appears in /health
healthRegistry.Register(NewBasicHealthCheck(health.ScopeBase))

healthHandler := handlers.NewHealthHandler(healthRegistry)
```

## Response Format

All health check endpoints return a response in this format:

```json
{
  "status": "UP",
  "components": {
    "database": {
      "status": "UP",
      "details": {
        "max_open_connections": 100,
        "open_connections": 2,
        "in_use": 0,
        "idle": 2
      }
    },
    "cache": {
      "status": "UP",
      "details": {
        "connected": true
      }
    }
  }
}
```

When any component has a status of `DOWN`, the overall status becomes `DOWN` and the HTTP status code is `503 Service Unavailable`.

## Scope Behavior

### Example Configuration

```go
// Provider A: registered in startup and ready scopes
providerA := NewProviderA(health.ScopeStartup, health.ScopeReady)

// Provider B: registered in live scope only
providerB := NewProviderB(health.ScopeLive)

// Provider C: registered in base scope only
providerC := NewProviderC(health.ScopeBase)

healthRegistry.Register(providerA)
healthRegistry.Register(providerB)
healthRegistry.Register(providerC)
```

### Endpoint Results

- `/health`: Returns A, B, C (all providers, each checked once)
- `/health/startup`: Returns A only
- `/health/readiness`: Returns A only
- `/health/liveness`: Returns B only

## Best Practices

1. **Use appropriate scopes**:
   - `ScopeStartup`: Critical dependencies needed for application startup (e.g., database)
   - `ScopeReady`: Dependencies needed to serve traffic (e.g., database, cache, message queue)
   - `ScopeLive`: Light checks that won't cause restarts (typically none or very basic checks)
   - `ScopeBase`: Informational checks that are expensive or not critical for K8s probes

2. **Add timeouts**: Always use context with timeout in your health checks to prevent hanging

3. **Keep it simple**: Health checks should be fast and lightweight

4. **Avoid cascading failures**: Liveness checks should not depend on external services to prevent unnecessary pod restarts

5. **Thread safety**: Health checks can be called concurrently, ensure your implementation is thread-safe

## Migration from Old System

The old system with direct database checks is now replaced with a provider-based system. The migration is backward compatible:

**Old code:**
```go
healthHandler := handlers.NewHealthHandler(db)
```

**New code:**
```go
healthRegistry := health.NewRegistry()
healthRegistry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup, health.ScopeReady))
healthHandler := handlers.NewHealthHandler(healthRegistry)
```

The response format and endpoints remain the same.
