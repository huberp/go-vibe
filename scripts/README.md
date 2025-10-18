# Build and Deployment Scripts

This directory contains shell scripts (.sh) for Linux/macOS and PowerShell scripts (.ps1) for Windows to manage the application.

## Available Scripts

### Database Scripts

Start and stop the PostgreSQL database server.

**Linux/macOS:**
```bash
# Start database
./scripts/db-start.sh

# Stop database
./scripts/db-stop.sh
```

**Windows PowerShell:**
```powershell
# Start database
.\scripts\db-start.ps1

# Stop database
.\scripts\db-stop.ps1
```

**Note:** These scripts require PostgreSQL to be installed and the database to be initialized in the `./data` directory. See [docs/database/postgresql.md](../docs/database/postgresql.md) for setup instructions.

### Build Script
Builds the application binary.

**Linux/macOS:**
```bash
./scripts/build.sh
```

**Windows PowerShell:**
```powershell
.\scripts\build.ps1
```

### Test Script
Runs all tests.

**Linux/macOS:**
```bash
./scripts/test.sh
```

**Windows PowerShell:**
```powershell
.\scripts\test.ps1
```

### Test with Coverage
Runs tests and generates coverage reports (HTML and text).

**Linux/macOS:**
```bash
./scripts/test-coverage.sh
```

**Windows PowerShell:**
```powershell
.\scripts\test-coverage.ps1
```

This generates:
- `coverage.out` - Text coverage data
- `coverage.html` - Interactive HTML coverage report

### Generate Swagger Documentation
Generates OpenAPI/Swagger documentation from code annotations.

**Linux/macOS:**
```bash
./scripts/swagger.sh
```

**Windows PowerShell:**
```powershell
.\scripts\swagger.ps1
```

View the generated docs at `http://localhost:8080/swagger/index.html` when the server is running.

### Database Migrations
Manage database schema migrations using golang-migrate.

**Linux/macOS:**
```bash
# Apply all pending migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down

# Create new migration
./scripts/migrate.sh create add_user_profile

# Force to specific version (recovery)
./scripts/migrate.sh force 1

# Show current version
./scripts/migrate.sh version
```

**Windows PowerShell:**
```powershell
# Apply all pending migrations
.\scripts\migrate.ps1 up

# Rollback last migration
.\scripts\migrate.ps1 down

# Create new migration
.\scripts\migrate.ps1 create add_user_profile

# Force to specific version (recovery)
.\scripts\migrate.ps1 force 1

# Show current version
.\scripts\migrate.ps1 version
```

**Note:** Set `DATABASE_URL` environment variable before running migrations:
```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/myapp?sslmode=disable"
```

### Run Server in Background
Builds and starts the server in the background. The server PID is saved to `server.pid`.

**Linux/macOS:**
```bash
# Set environment variables first
export DATABASE_URL="postgres://user:password@localhost:5432/myapp?sslmode=disable"
export JWT_SECRET="your-secret-key"
export SERVER_PORT="8080"

# Start server
./scripts/run-background.sh
```

**Windows PowerShell:**
```powershell
# Set environment variables first
$env:DATABASE_URL="postgres://user:password@localhost:5432/myapp?sslmode=disable"
$env:JWT_SECRET="your-secret-key"
$env:SERVER_PORT="8080"

# Start server
.\scripts\run-background.ps1
```

**Logs:**
- Linux/macOS: `tail -f server.log`
- Windows: `Get-Content server.log -Wait`

### Stop Server
Stops the background server gracefully.

**Linux/macOS:**
```bash
./scripts/stop.sh
```

**Windows PowerShell:**
```powershell
.\scripts\stop.ps1
```

## Workflow Testing

These scripts are automatically tested in the `scripts-test.yml` GitHub Actions workflow, which verifies:
- Building the application
- Running tests
- Starting/stopping the server in background mode

## Local Kubernetes Testing

### Validate Deployment Configuration

Validates Helm charts and Kubernetes manifests without requiring a running cluster:

**Linux/macOS:**
```bash
./scripts/validate-k8s.sh
```

This script performs:
- Helm chart linting
- Kubernetes manifest generation
- Manifest syntax validation
- Dockerfile verification

### Setup Local Cluster
Creates a local Kind (Kubernetes in Docker) cluster for testing deployments.

**Linux/macOS:**
```bash
./scripts/local-k8s-setup.sh
```

**Windows PowerShell:**
```powershell
.\scripts\local-k8s-setup.ps1
```

### Deploy to Local Cluster
Builds the Docker image, loads it into Kind, and deploys with Helm (includes PostgreSQL).

**Linux/macOS:**
```bash
./scripts/local-k8s-deploy.sh
```

**Windows PowerShell:**
```powershell
.\scripts\local-k8s-deploy.ps1
```

### Access the Application
After deployment, port-forward to access the application:

```bash
kubectl port-forward -n production svc/myapp 8080:8080
```

Then test the API:
```bash
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

### View Logs
```bash
kubectl logs -n production -l app=myapp -f
```

### Cleanup Local Cluster
Provides options to clean up the deployment, namespace, or entire cluster.

**Linux/macOS:**
```bash
./scripts/local-k8s-cleanup.sh
```

**Windows PowerShell:**
```powershell
.\scripts\local-k8s-cleanup.ps1
```

## Notes

- All scripts check for errors and exit with appropriate status codes
- The run-background script saves the server PID to `server.pid`
- Server logs are written to `server.log` (and `server-error.log` on Windows)
- The stop script cleans up PID files automatically
- Local Kubernetes scripts use Kind cluster named `go-vibe-local`
- Local deployments use the `production` namespace for testing
