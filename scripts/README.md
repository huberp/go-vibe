# Build and Deployment Scripts

This directory contains shell scripts (.sh) for Linux/macOS and PowerShell scripts (.ps1) for Windows to manage the application.

## Available Scripts

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

## Notes

- All scripts check for errors and exit with appropriate status codes
- The run-background script saves the server PID to `server.pid`
- Server logs are written to `server.log` (and `server-error.log` on Windows)
- The stop script cleans up PID files automatically
