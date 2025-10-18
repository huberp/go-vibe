# Info Package

The `info` package provides an extensible framework for aggregating and exposing application information through a unified API. It's designed to be easily extracted as an external Go module.

## Overview

This package implements a pluggable provider pattern inspired by Spring Boot's actuator info endpoint. It allows applications to expose various types of information (build details, statistics, custom metrics) through a simple and extensible interface.

## Core Components

### InfoProvider Interface

The `InfoProvider` interface is the foundation of the info system:

```go
type InfoProvider interface {
    // Name returns the unique name of this provider
    Name() string
    
    // Info returns the information provided by this provider
    Info() (map[string]interface{}, error)
}
```

### Registry

The `Registry` manages multiple `InfoProvider` instances and aggregates their data:

```go
registry := info.NewRegistry()
registry.Register(provider1)
registry.Register(provider2)

// Get aggregated information from all providers
allInfo := registry.GetAll()
```

## Built-in Providers

### BuildInfoProvider

Provides build-time information about the application:

```go
buildProvider := info.NewBuildInfoProvider(
    "1.0.0",           // version
    "abc123def",       // commit SHA
    "2024-01-01T...",  // build time
    "go1.25.2",        // Go version
)
```

Returns:
```json
{
  "build": {
    "version": "1.0.0",
    "commit": "abc123def",
    "build_time": "2024-01-01T12:00:00Z",
    "go_version": "go1.25.2"
  }
}
```

### UserStatsProvider

Provides user-related statistics from a GORM database:

```go
userStatsProvider := info.NewUserStatsProvider(db)
```

Returns:
```json
{
  "users": {
    "total": 100,
    "admins": 5,
    "regular": 95
  }
}
```

## Usage Example

### Basic Setup

```go
import (
    "myapp/pkg/info"
    "runtime"
)

// Create registry
registry := info.NewRegistry()

// Register providers
registry.Register(info.NewBuildInfoProvider(
    "v1.0.0",
    "abc123",
    "",  // will use current time
    runtime.Version(),
))

registry.Register(info.NewUserStatsProvider(db))

// Get all info
allInfo := registry.GetAll()
```

### HTTP Handler Integration

```go
import (
    "github.com/gin-gonic/gin"
    "myapp/pkg/info"
)

func setupInfoEndpoint(router *gin.Engine, registry *info.Registry) {
    router.GET("/info", func(c *gin.Context) {
        c.JSON(200, registry.GetAll())
    })
}
```

## Creating Custom Providers

Implement the `InfoProvider` interface:

```go
type DatabaseStatsProvider struct {
    db *gorm.DB
}

func NewDatabaseStatsProvider(db *gorm.DB) *DatabaseStatsProvider {
    return &DatabaseStatsProvider{db: db}
}

func (p *DatabaseStatsProvider) Name() string {
    return "database"
}

func (p *DatabaseStatsProvider) Info() (map[string]interface{}, error) {
    var stats struct {
        MaxConnections int
        OpenConnections int
        InUse int
        Idle int
    }
    
    sqlDB, err := p.db.DB()
    if err != nil {
        return nil, err
    }
    
    dbStats := sqlDB.Stats()
    stats.MaxConnections = dbStats.MaxOpenConnections
    stats.OpenConnections = dbStats.OpenConnections
    stats.InUse = dbStats.InUse
    stats.Idle = dbStats.Idle
    
    return map[string]interface{}{
        "max_connections": stats.MaxConnections,
        "open_connections": stats.OpenConnections,
        "in_use": stats.InUse,
        "idle": stats.Idle,
    }, nil
}

// Register the custom provider
registry.Register(NewDatabaseStatsProvider(db))
```

## Thread Safety

The `Registry` is thread-safe and can be safely accessed from multiple goroutines. All built-in providers are also designed to be thread-safe.

## Error Handling

If a provider returns an error, it will be silently omitted from the aggregated results. This ensures that one failing provider doesn't break the entire info endpoint.

## Extracting as External Module

This package is structured to be easily extracted as a standalone Go module:

1. **Self-contained**: Core interfaces (`InfoProvider`, `Registry`) have no external dependencies
2. **Optional providers**: Built-in providers (like `UserStatsProvider`) can be separated into sub-packages
3. **Clean interfaces**: Well-defined interfaces make it easy to create custom implementations

To extract as a module:

```bash
# Create new repository
mkdir go-info-provider
cd go-info-provider
go mod init github.com/yourorg/go-info-provider

# Copy core files
cp pkg/info/provider.go .
cp pkg/info/registry.go .
cp pkg/info/build_provider.go ./providers/

# Update imports
# Publish as go module
```

## Best Practices

1. **Provider Names**: Use descriptive, unique names for providers (e.g., "build", "users", "database")
2. **Error Handling**: Return errors from `Info()` method when data cannot be retrieved
3. **Performance**: Keep `Info()` methods lightweight - avoid expensive operations
4. **Data Format**: Return JSON-serializable data (primitives, maps, slices)
5. **Security**: Don't expose sensitive information (credentials, secrets, PII)

## Testing

All components include comprehensive unit tests:

```bash
go test ./pkg/info/... -v
```

Example test pattern:

```go
func TestCustomProvider(t *testing.T) {
    provider := NewCustomProvider(dependencies)
    
    info, err := provider.Info()
    assert.NoError(t, err)
    assert.Contains(t, info, "expected_key")
}
```

## API Response Format

The aggregated info response follows this structure:

```json
{
  "provider_name_1": {
    "key1": "value1",
    "key2": "value2"
  },
  "provider_name_2": {
    "metric1": 123,
    "metric2": 456
  }
}
```

Each provider's data is nested under its name, preventing key collisions.
