# Health Check Package

A reusable, extensible health check framework for Go applications with support for Kubernetes probes.

## Overview

The `health` package provides a flexible system for monitoring application health across multiple scopes (startup, liveness, readiness). It follows the provider pattern, making it easy to add custom health checks.

## Features

- ✅ **Multiple Health Scopes**: Support for startup, liveness, readiness, and base health checks
- ✅ **Kubernetes Ready**: Built-in support for K8s startup, liveness, and readiness probes
- ✅ **Extensible**: Easy to add custom health check providers
- ✅ **Thread-Safe**: Registry supports concurrent access
- ✅ **Minimal Dependencies**: Core types have no external dependencies; optional providers use GORM
- ✅ **JSON Response Types**: Ready-to-use HTTP response structures

## Quick Start

### 1. Create a Registry

```go
import "yourapp/pkg/health"

registry := health.NewRegistry()
```

### 2. Register Health Check Providers

```go
// Database health check (startup + readiness scopes)
registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup, health.ScopeReady))

// Simple liveness check
registry.Register(health.NewSimpleHealthCheckProvider("app", health.ScopeLive))
```

### 3. Check Health

```go
// Check all providers
allResults := registry.Check(nil)

// Check only startup scope
scope := health.ScopeStartup
startupResults := registry.Check(&scope)

// Determine overall status
overallStatus := health.OverallStatus(allResults)
```

### 4. Build HTTP Response

```go
components := make(map[string]health.ComponentHealth)
for name, result := range checkResults {
    components[name] = health.ComponentHealth{
        Status:  result.Status,
        Details: result.Details,
    }
}

response := health.Response{
    Status:     health.OverallStatus(checkResults),
    Components: components,
}

// Marshal to JSON
jsonData, _ := json.Marshal(response)
```

## Health Check Scopes

The package supports four health check scopes:

| Scope | Constant | Description | Use Case |
|-------|----------|-------------|----------|
| **Base** | `ScopeBase` | Only appears in `/health` | Additional checks not needed for probes |
| **Startup** | `ScopeStartup` | Appears in `/health/startup` and `/health` | One-time initialization checks |
| **Liveness** | `ScopeLive` | Appears in `/health/liveness` and `/health` | Checks if app is running (restart if down) |
| **Readiness** | `ScopeReady` | Appears in `/health/readiness` and `/health` | Checks if app can accept traffic |

### Scope Behavior

- A provider can be registered in **multiple scopes**
- When checking all scopes (nil), each provider is **checked only once**
- Scoped checks only include providers with matching scope

Example:
```go
// Provider appears in startup and readiness
dbProvider := health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup, health.ScopeReady)
registry.Register(dbProvider)

// Check startup scope - includes database
startupScope := health.ScopeStartup
startupResults := registry.Check(&startupScope)

// Check liveness scope - does NOT include database
liveScope := health.ScopeLive
liveResults := registry.Check(&liveScope)

// Check all - includes database once
allResults := registry.Check(nil)
```

## Built-in Providers

### DatabaseHealthCheckProvider

Checks PostgreSQL database connectivity with connection pool statistics.

```go
provider := health.NewDatabaseHealthCheckProvider(
    db,                      // *gorm.DB
    health.ScopeStartup,     // appears in startup probe
    health.ScopeReady,       // appears in readiness probe
)
```

**Response when healthy:**
```json
{
  "status": "UP",
  "details": {
    "max_open_connections": 100,
    "open_connections": 5,
    "in_use": 2,
    "idle": 3
  }
}
```

**Response when unhealthy:**
```json
{
  "status": "DOWN",
  "details": {
    "error": "database connection failed"
  }
}
```

### SimpleHealthCheckProvider

A basic health check that always returns UP status with a timestamp.

```go
provider := health.NewSimpleHealthCheckProvider(
    "app",              // name
    health.ScopeLive,   // scope(s)
)
```

**Response:**
```json
{
  "status": "UP",
  "details": {
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

## Creating Custom Providers

Implement the `HealthCheckProvider` interface:

```go
type CustomHealthProvider struct {
    name   string
    scopes []health.Scope
}

func (c *CustomHealthProvider) Name() string {
    return c.name
}

func (c *CustomHealthProvider) Check() (*health.CheckResult, error) {
    // Perform your health check
    isHealthy := performCheck()
    
    if isHealthy {
        return &health.CheckResult{
            Status: health.StatusUp,
            Details: map[string]interface{}{
                "message": "all systems operational",
            },
        }, nil
    }
    
    return &health.CheckResult{
        Status: health.StatusDown,
        Details: map[string]interface{}{
            "error": "system degraded",
        },
    }, nil
}

func (c *CustomHealthProvider) Scopes() []health.Scope {
    return c.scopes
}
```

## Response Types

### CheckResult

Internal result from a single health check:

```go
type CheckResult struct {
    Status  Status                 `json:"status"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```

### ComponentHealth

Health status of a single component (for HTTP responses):

```go
type ComponentHealth struct {
    Status  Status                 `json:"status"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```

### Response

Overall health response (for HTTP endpoints):

```go
type Response struct {
    Status     Status                     `json:"status"`
    Components map[string]ComponentHealth `json:"components,omitempty"`
}
```

## Kubernetes Integration

Use the health package with Kubernetes probes:

```yaml
livenessProbe:
  httpGet:
    path: /health/liveness
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30

readinessProbe:
  httpGet:
    path: /health/readiness
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10

startupProbe:
  httpGet:
    path: /health/startup
    port: 8080
  initialDelaySeconds: 0
  periodSeconds: 5
  failureThreshold: 30
```

## Example HTTP Handler

```go
func HealthCheckHandler(c *gin.Context) {
    checkResults := registry.Check(nil)
    
    components := make(map[string]health.ComponentHealth)
    for name, result := range checkResults {
        components[name] = health.ComponentHealth{
            Status:  result.Status,
            Details: result.Details,
        }
    }
    
    overallStatus := health.OverallStatus(checkResults)
    response := health.Response{
        Status:     overallStatus,
        Components: components,
    }
    
    statusCode := http.StatusOK
    if overallStatus == health.StatusDown {
        statusCode = http.StatusServiceUnavailable
    }
    
    c.JSON(statusCode, response)
}
```

## Complete Example

```go
package main

import (
    "yourapp/pkg/health"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupHealthChecks(router *gin.Engine, db *gorm.DB) {
    // Create registry
    registry := health.NewRegistry()
    
    // Register providers
    registry.Register(health.NewDatabaseHealthCheckProvider(
        db,
        health.ScopeStartup,
        health.ScopeReady,
    ))
    
    registry.Register(health.NewSimpleHealthCheckProvider(
        "app",
        health.ScopeLive,
    ))
    
    // Setup routes
    router.GET("/health", func(c *gin.Context) {
        handleHealthCheck(c, registry, nil)
    })
    
    router.GET("/health/startup", func(c *gin.Context) {
        scope := health.ScopeStartup
        handleHealthCheck(c, registry, &scope)
    })
    
    router.GET("/health/liveness", func(c *gin.Context) {
        scope := health.ScopeLive
        handleHealthCheck(c, registry, &scope)
    })
    
    router.GET("/health/readiness", func(c *gin.Context) {
        scope := health.ScopeReady
        handleHealthCheck(c, registry, &scope)
    })
}

func handleHealthCheck(c *gin.Context, registry *health.Registry, scope *health.Scope) {
    checkResults := registry.Check(scope)
    
    components := make(map[string]health.ComponentHealth)
    for name, result := range checkResults {
        components[name] = health.ComponentHealth{
            Status:  result.Status,
            Details: result.Details,
        }
    }
    
    overallStatus := health.OverallStatus(checkResults)
    response := health.Response{
        Status:     overallStatus,
        Components: components,
    }
    
    statusCode := http.StatusOK
    if overallStatus == health.StatusDown {
        statusCode = http.StatusServiceUnavailable
    }
    
    c.JSON(statusCode, response)
}
```

## Best Practices

1. **Use appropriate scopes**:
   - Startup: Database migrations, one-time initialization
   - Liveness: Simple checks that indicate the app is running
   - Readiness: Checks that determine if the app can serve traffic

2. **Keep liveness checks simple**: Complex liveness checks can cause unnecessary restarts

3. **Make readiness checks comprehensive**: Include all dependencies needed to serve traffic

4. **Return meaningful details**: Include connection counts, latencies, error messages

5. **Use timeouts**: Implement timeouts in health checks to prevent hanging

## License

Apache 2.0
