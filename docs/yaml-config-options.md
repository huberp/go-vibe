# YAML-based Configuration Options for go-vibe

## Current State
- Configuration loaded from environment variables using Viper
- No YAML file support
- No stage/profile support
- Manual environment variable setting required

## Requirements
1. Central YAML config file for all stages
2. Support for including/importing other YAML files
3. Different deployment stages support
4. Easy stage selection for users
5. Backward compatible with existing environment variable approach

## Option 1: File-per-Stage with Base Config (Recommended)

### Structure
```
config/
├── base.yaml              # Base/shared configuration
├── development.yaml       # Development overrides
├── staging.yaml          # Staging overrides
└── production.yaml       # Production overrides
```

### Stage Selection
- Command line flag: `--stage=development` or `--stage=production`
- Environment variable: `APP_STAGE=development`
- Default to `development`

### Config Loading Order
1. Load base.yaml (shared config)
2. Load stage-specific config (e.g., production.yaml)
3. Override with environment variables (backward compatibility)

### Pros
- Clear separation of stage configs
- Easy to understand and maintain
- No complex include directives needed
- Simple file structure
- Each stage file only contains overrides

### Cons
- Potential duplication if not using base config properly
- Need to manage multiple files

### Implementation Approach
```yaml
# config/base.yaml
server:
  port: 8080

database:
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: 30

# config/production.yaml
database:
  url: postgres://prod-user:pass@prod-host:5432/myapp
  max_open_conns: 100

jwt:
  secret: ${JWT_SECRET}  # Can still use env vars for secrets
```

---

## Option 2: Profile-based Single File (Spring Boot style)

### Structure
```
config/
└── application.yaml       # Single file with all profiles
```

### Profile Definition
```yaml
# Default/shared configuration
server:
  port: 8080

database:
  max_open_conns: 25

---
# Development profile
spring:
  profiles: development

database:
  url: postgres://localhost:5432/myapp_dev

---
# Production profile  
spring:
  profiles: production

database:
  url: postgres://prod-host:5432/myapp
  max_open_conns: 100
```

### Stage Selection
- Command line: `--profile=production`
- Environment variable: `APP_PROFILE=production`

### Pros
- Single file to manage
- All configurations visible in one place
- Familiar to Spring Boot users

### Cons
- Single file can become large with many profiles
- Less clear separation
- Requires custom parsing logic for profile sections
- Viper doesn't natively support Spring-style profiles

---

## Option 3: Include-based Configuration (Most Flexible)

### Structure
```
config/
├── main.yaml              # Main config with includes
├── common/
│   ├── database.yaml
│   ├── server.yaml
│   └── logging.yaml
└── stages/
    ├── development.yaml
    ├── staging.yaml
    └── production.yaml
```

### Include Directive Format
```yaml
# config/main.yaml
includes:
  - common/database.yaml
  - common/server.yaml
  - stages/${APP_STAGE}.yaml

# Values can be overridden after includes
server:
  port: ${SERVER_PORT:8080}
```

### Stage Selection
- Environment variable: `APP_STAGE=production`
- Command line flag: `--stage=production`

### Pros
- Maximum flexibility
- Modular configuration
- Can reuse common configs
- Easy to organize complex configs

### Cons
- Most complex to implement
- Requires custom include processing
- Can be harder to understand the final merged config
- More files to manage

---

## Option 4: Hybrid Approach (File-per-Stage + Includes)

### Structure
```
config/
├── base.yaml              # Base config
├── common/
│   ├── database.yaml      # Database defaults
│   └── security.yaml      # Security defaults
└── stages/
    ├── development.yaml   # Dev overrides with includes
    ├── staging.yaml       # Staging overrides with includes
    └── production.yaml    # Prod overrides with includes
```

### Example Stage File with Includes
```yaml
# config/stages/production.yaml
includes:
  - ../base.yaml
  - ../common/database.yaml
  - ../common/security.yaml

# Production-specific overrides
database:
  url: postgres://prod-host:5432/myapp
  max_open_conns: 100

jwt:
  secret: ${JWT_SECRET}
```

### Pros
- Combines benefits of both approaches
- Modular and organized
- Clear stage separation

### Cons
- Most complex file structure
- Potentially confusing include paths
- Overhead of managing includes

---

## Recommendation: Option 1 (File-per-Stage with Base Config)

### Rationale
1. **Simplicity**: Easy to understand and implement
2. **Viper Support**: Viper has built-in support for merging configs
3. **Clear Separation**: Each stage is clearly defined
4. **Backward Compatible**: Can still use environment variables
5. **Maintainable**: Simple file structure
6. **Quick Implementation**: Can be done quickly with existing Viper features

### Implementation Plan
1. Create `config/` directory structure
2. Create base.yaml with defaults
3. Create stage-specific YAML files (development, staging, production)
4. Modify `pkg/config/config.go` to:
   - Accept stage parameter (from flag or env var)
   - Load base.yaml first
   - Merge stage-specific config
   - Allow environment variable overrides
5. Update main.go to pass stage selection
6. Update documentation
7. Add tests for each stage
8. Maintain backward compatibility

### Migration Path
- Existing deployments continue to work with env vars
- New deployments can use YAML configs
- Secrets can still be injected via env vars (recommended for production)

### Usage Examples
```bash
# Using stage flag
./server --stage=production

# Using environment variable
export APP_STAGE=production
./server

# Override specific values
export APP_STAGE=production
export DATABASE_URL=postgres://custom-host:5432/db
./server

# Development (default)
./server
```

---

## Questions for @huberp

1. **Which option do you prefer?** (Recommend Option 1)
   - Option 1: File-per-Stage with Base Config
   - Option 2: Profile-based Single File
   - Option 3: Include-based Configuration
   - Option 4: Hybrid Approach

2. **Stage names**: Do you want to keep `development`, `staging`, `production` or different names?

3. **Include directive priority**: If we add includes later, should we use a specific directive name like `includes:` or `import:`?

4. **Secrets handling**: Continue using environment variables for secrets (JWT_SECRET, DATABASE_URL) or support encrypted values in YAML?

5. **Default stage**: Should the default be `development` or require explicit stage selection?

Please confirm your preference and I'll create a detailed implementation plan!
