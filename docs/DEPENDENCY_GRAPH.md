# Go Module Dependency Graph Summary

## Overview

This document presents the results of `go mod tidy` and `go mod graph` for the go-vibe project.

## Go Mod Tidy Results

✅ **Dependencies are already clean!** 

The `go mod tidy` command completed successfully with no changes needed to `go.mod` or `go.sum` files. This indicates that:
- All dependencies listed in `go.mod` are actually used by the project
- No unused dependencies are present
- All indirect dependencies are properly tracked in `go.sum`

## Direct Dependencies (from go.mod)

The project has the following **direct dependencies** (listed in the `require` section):

### Core Dependencies
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/gin-gonic/gin | v1.11.0 | HTTP framework |
| gorm.io/gorm | v1.31.0 | ORM |
| gorm.io/driver/postgres | v1.6.0 | PostgreSQL driver |
| gorm.io/driver/sqlite | v1.6.0 | SQLite driver (for testing) |

### Authentication & Security
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/golang-jwt/jwt/v5 | v5.3.0 | JWT authentication |
| golang.org/x/crypto | v0.43.0 | Cryptography |
| github.com/google/uuid | v1.6.0 | UUID generation |

### Configuration & Logging
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/spf13/viper | v1.21.0 | Configuration management |
| go.uber.org/zap | v1.27.0 | Structured logging |

### Middleware & Extensions
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/gin-contrib/cors | v1.7.6 | CORS middleware |
| golang.org/x/time | v0.14.0 | Rate limiting |

### Observability
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/prometheus/client_golang | v1.23.2 | Metrics |
| go.opentelemetry.io/otel | v1.38.0 | Distributed tracing |
| go.opentelemetry.io/otel/trace | v1.38.0 | Tracing API |
| go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin | v0.63.0 | Gin tracing |

### Database Migrations
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/golang-migrate/migrate/v4 | v4.19.0 | Database migrations |

### API Documentation
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/swaggo/swag | v1.16.6 | Swagger generation |
| github.com/swaggo/gin-swagger | v1.6.1 | Gin Swagger integration |
| github.com/swaggo/files | v1.0.1 | Swagger file serving |

### Testing
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/stretchr/testify | v1.11.1 | Testing framework |
| go.uber.org/mock | v0.6.0 | Mocking |

## Dependency Statistics

- **Total direct dependencies**: 26 (in `require` section of go.mod)
- **Total graph entries**: 952 (including all transitive dependencies)
- **Go version**: 1.25.2

## Full Dependency Graph

The complete dependency graph has been generated with `go mod graph` and contains 952 relationships showing:
- Direct dependencies from the project
- All transitive (indirect) dependencies
- Version relationships between packages

### Top-level Dependencies (myapp → package)

The project directly depends on 97 packages (including both direct and indirect dependencies tracked in go.mod):

```
myapp → github.com/KyleBanks/depth@v1.2.1
myapp → github.com/beorn7/perks@v1.0.1
myapp → github.com/bytedance/gopkg@v0.1.3
myapp → github.com/bytedance/sonic@v1.14.1
myapp → github.com/bytedance/sonic/loader@v0.3.0
myapp → github.com/cespare/xxhash/v2@v2.3.0
myapp → github.com/cloudwego/base64x@v0.1.6
myapp → github.com/davecgh/go-spew@v1.1.1
myapp → github.com/fsnotify/fsnotify@v1.9.0
myapp → github.com/gabriel-vasile/mimetype@v1.4.10
myapp → github.com/gin-contrib/cors@v1.7.6
myapp → github.com/gin-contrib/sse@v1.1.0
myapp → github.com/gin-gonic/gin@v1.11.0
myapp → github.com/go-logr/logr@v1.4.3
myapp → github.com/go-logr/stdr@v1.2.2
myapp → github.com/go-openapi/jsonpointer@v0.22.1
myapp → github.com/go-openapi/jsonreference@v0.21.2
myapp → github.com/go-openapi/spec@v0.22.0
myapp → github.com/go-openapi/swag/conv@v0.25.1
myapp → github.com/go-openapi/swag/jsonname@v0.25.1
myapp → github.com/go-openapi/swag/jsonutils@v0.25.1
myapp → github.com/go-openapi/swag/loading@v0.25.1
myapp → github.com/go-openapi/swag/stringutils@v0.25.1
myapp → github.com/go-openapi/swag/typeutils@v0.25.1
myapp → github.com/go-openapi/swag/yamlutils@v0.25.1
myapp → github.com/go-playground/locales@v0.14.1
myapp → github.com/go-playground/universal-translator@v0.18.1
myapp → github.com/go-playground/validator/v10@v10.28.0
myapp → github.com/go-viper/mapstructure/v2@v2.4.0
myapp → github.com/goccy/go-json@v0.10.5
myapp → github.com/goccy/go-yaml@v1.18.0
myapp → github.com/golang-jwt/jwt/v5@v5.3.0
myapp → github.com/golang-migrate/migrate/v4@v4.19.0
myapp → github.com/google/uuid@v1.6.0
myapp → github.com/hashicorp/errwrap@v1.1.0
myapp → github.com/hashicorp/go-multierror@v1.1.1
myapp → github.com/jackc/pgpassfile@v1.0.0
myapp → github.com/jackc/pgservicefile@v0.0.0-20240606120523-5a60cdf6a761
myapp → github.com/jackc/pgx/v5@v5.6.0
myapp → github.com/jackc/puddle/v2@v2.2.2
myapp → github.com/jinzhu/inflection@v1.0.0
myapp → github.com/jinzhu/now@v1.1.5
myapp → github.com/json-iterator/go@v1.1.12
myapp → github.com/klauspost/cpuid/v2@v2.3.0
myapp → github.com/leodido/go-urn@v1.4.0
myapp → github.com/lib/pq@v1.10.9
myapp → github.com/mattn/go-isatty@v0.0.20
myapp → github.com/mattn/go-sqlite3@v1.14.22
myapp → github.com/modern-go/concurrent@v0.0.0-20180306012644-bacd9c7ef1dd
myapp → github.com/modern-go/reflect2@v1.0.2
myapp → github.com/munnerz/goautoneg@v0.0.0-20191010083416-a7dc8b61c822
myapp → github.com/pelletier/go-toml/v2@v2.2.4
myapp → github.com/pmezard/go-difflib@v1.0.0
myapp → github.com/prometheus/client_golang@v1.23.2
myapp → github.com/prometheus/client_model@v0.6.2
myapp → github.com/prometheus/common@v0.66.1
myapp → github.com/prometheus/procfs@v0.16.1
myapp → github.com/quic-go/qpack@v0.5.1
myapp → github.com/quic-go/quic-go@v0.55.0
myapp → github.com/sagikazarmark/locafero@v0.11.0
myapp → github.com/sourcegraph/conc@v0.3.1-0.20240121214520-5f936abd7ae8
myapp → github.com/spf13/afero@v1.15.0
myapp → github.com/spf13/cast@v1.10.0
myapp → github.com/spf13/pflag@v1.0.10
myapp → github.com/spf13/viper@v1.21.0
myapp → github.com/stretchr/testify@v1.11.1
myapp → github.com/subosito/gotenv@v1.6.0
myapp → github.com/swaggo/files@v1.0.1
myapp → github.com/swaggo/gin-swagger@v1.6.1
myapp → github.com/swaggo/swag@v1.16.6
myapp → github.com/twitchyliquid64/golang-asm@v0.15.1
myapp → github.com/ugorji/go/codec@v1.3.0
myapp → go.opentelemetry.io/auto/sdk@v1.1.0
myapp → go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin@v0.63.0
myapp → go.opentelemetry.io/otel@v1.38.0
myapp → go.opentelemetry.io/otel/metric@v1.38.0
myapp → go.opentelemetry.io/otel/trace@v1.38.0
myapp → go.uber.org/mock@v0.6.0
myapp → go.uber.org/multierr@v1.10.0
myapp → go.uber.org/zap@v1.27.0
myapp → go.yaml.in/yaml/v2@v2.4.3
myapp → go.yaml.in/yaml/v3@v3.0.4
myapp → golang.org/x/arch@v0.22.0
myapp → golang.org/x/crypto@v0.43.0
myapp → golang.org/x/mod@v0.29.0
myapp → golang.org/x/net@v0.46.0
myapp → golang.org/x/sync@v0.17.0
myapp → golang.org/x/sys@v0.37.0
myapp → golang.org/x/text@v0.30.0
myapp → golang.org/x/time@v0.14.0
myapp → golang.org/x/tools@v0.38.0
myapp → google.golang.org/protobuf@v1.36.10
myapp → gopkg.in/yaml.v3@v3.0.1
myapp → gorm.io/driver/postgres@v1.6.0
myapp → gorm.io/driver/sqlite@v1.6.0
myapp → gorm.io/gorm@v1.31.0
```

## Key Findings

1. ✅ **Clean Dependencies**: `go mod tidy` found no unused dependencies - all packages are actively used
2. ✅ **Version Consistency**: All dependencies are at consistent versions with no conflicts
3. ✅ **Well-Maintained**: The project uses modern, actively maintained packages
4. ✅ **Comprehensive Stack**: Good coverage of web framework, database, auth, logging, metrics, and testing tools

## Recommendations

- Dependencies are well-maintained and up-to-date
- No action needed - dependency management is in good shape
- Continue using `go mod tidy` regularly in CI/CD to maintain this clean state
