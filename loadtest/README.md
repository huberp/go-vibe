# Load Testing

This directory contains load testing infrastructure for the go-vibe microservice using [k6](https://k6.io/), a modern load testing tool built on Go.

## Directory Structure

```
loadtest/
├── scripts/              # k6 test scripts
│   ├── smoke-test.js    # Basic smoke test (1 VU, 30s)
│   ├── auth-load-test.js # Authentication load test (10 VUs)
│   ├── stress-test.js   # Stress test (up to 100 VUs)
│   └── user-crud-test.js # Full CRUD operations test
├── data/                # Data generation and cleanup
│   ├── generate/        # Generate test users
│   │   └── main.go
│   └── cleanup/         # Clean up test data
│       └── main.go
├── run-smoke-test.sh    # Run smoke test
├── run-auth-test.sh     # Run authentication test
└── run-stress-test.sh   # Run stress test
```

## Prerequisites

### Install k6

**macOS (Homebrew):**
```bash
brew install k6
```

**Linux:**
```bash
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

**Windows (Chocolatey):**
```powershell
choco install k6
```

For other installation methods, see: https://k6.io/docs/getting-started/installation/

## Quick Start

### 1. Prepare Test Data (Optional)

Load test fixtures into your database:

```bash
# Using Docker Compose database
export DATABASE_URL="postgres://user:password@localhost:5432/myapp?sslmode=disable"

# Generate 100 test users
cd loadtest/data
go run generate-users.go -db "$DATABASE_URL" -users 100

# Later, clean up test data
go run cleanup-users.go -db "$DATABASE_URL"
```

### 2. Start the Application

Make sure your application is running:

```bash
# From project root
./server
# or
docker-compose up
```

### 3. Run Load Tests

**Smoke Test** (1 user, 30 seconds):
```bash
cd loadtest
./run-smoke-test.sh http://localhost:8080
```

**Authentication Load Test** (10 users, 2 minutes):
```bash
./run-auth-test.sh http://localhost:8080
```

**Stress Test** (up to 100 users, 12 minutes):
```bash
./run-stress-test.sh http://localhost:8080
```

**CRUD Test** (20 users, 5 minutes):
```bash
k6 run --env BASE_URL=http://localhost:8080 scripts/user-crud-test.js
```

## Test Scripts

### smoke-test.js
**Purpose:** Quick validation that the system works under minimal load.

**Configuration:**
- Virtual Users: 1
- Duration: 30 seconds
- Endpoints: `/health`, `/info`

**Thresholds:**
- 95% of requests < 500ms
- Error rate < 1%

**Usage:**
```bash
k6 run --env BASE_URL=http://localhost:8080 scripts/smoke-test.js
```

### auth-load-test.js
**Purpose:** Test authentication endpoints under realistic load.

**Configuration:**
- Virtual Users: Ramps from 0 → 10 → 0
- Duration: 2 minutes
- Endpoints: `/v1/login`, `/v1/users`, `/v1/users/{id}`

**Thresholds:**
- 95% of requests < 1s
- Login requests < 1.5s
- Error rate < 5%

**Features:**
- Tests login with multiple user roles
- Uses fixture user credentials
- Validates JWT tokens

**Usage:**
```bash
k6 run --env BASE_URL=http://localhost:8080 scripts/auth-load-test.js
```

### stress-test.js
**Purpose:** Push the system beyond normal capacity to find breaking points.

**Configuration:**
- Virtual Users: Ramps 0 → 50 → 100 → 0
- Duration: 12 minutes
- Endpoints: `/health`, `/metrics`, `/info`

**Thresholds:**
- 95% of requests < 2s
- Error rate < 10%

**Usage:**
```bash
k6 run --env BASE_URL=http://localhost:8080 scripts/stress-test.js
```

### user-crud-test.js
**Purpose:** Test full user lifecycle (create, login, read).

**Configuration:**
- Virtual Users: 20 for 3 minutes
- Endpoints: All user CRUD operations

**Thresholds:**
- 95% of requests < 1.5s
- Create user < 2s
- Error rate < 5%

**Features:**
- Creates users with random data
- Tests login after user creation
- Admin operations (requires admin credentials)

**Usage:**
```bash
k6 run --env BASE_URL=http://localhost:8080 \
       --env ADMIN_EMAIL=admin@example.com \
       --env ADMIN_PASSWORD=password123 \
       scripts/user-crud-test.js
```

## Test Data Management

### Generating Test Users

The `generate` command creates test users in the database:

```bash
cd loadtest/data/generate

# Generate 100 users
go run main.go -db "postgres://..." -users 100

# Generate 500 users without admin
go run main.go -db "postgres://..." -users 500 -admin=false

# Generate 1000 users
go run main.go -db "postgres://..." -users 1000
```

**Options:**
- `-db`: Database connection URL (required)
- `-users`: Number of users to generate (default: 100)
- `-admin`: Include admin user (default: true)

**Generated Users:**
- Email pattern: `loadtest-user-N@example.com`
- Password: `password123` (bcrypt hashed)
- Role: `user` (or `admin` for the first user if `-admin=true`)

### Cleaning Up Test Data

The `cleanup` command removes all load test users:

```bash
cd loadtest/data/cleanup
go run main.go -db "postgres://..."
```

This deletes all users with emails matching `loadtest-%@example.com`.

## Analyzing Results

### Understanding k6 Output

k6 provides detailed metrics after each test:

```
     ✓ health check status is 200
     ✓ login status is 200

     checks.........................: 100.00% ✓ 1200  ✗ 0   
     data_received..................: 2.4 MB  40 kB/s
     data_sent......................: 240 kB  4.0 kB/s
     http_req_blocked...............: avg=1.2ms   min=1µs   med=5µs    max=100ms p(95)=10ms  
     http_req_duration..............: avg=250ms   min=50ms  med=200ms  max=1s    p(95)=500ms 
     http_req_failed................: 0.00%   ✓ 0     ✗ 1200
     http_reqs......................: 1200    20/s
     iteration_duration.............: avg=2.5s    min=1s    med=2s     max=5s    p(95)=4s    
```

**Key Metrics:**
- `checks`: Percentage of assertion checks that passed
- `http_req_duration`: Response time statistics (avg, p95, max)
- `http_req_failed`: Percentage of failed requests
- `http_reqs`: Total requests and requests per second
- `iteration_duration`: Time to complete one iteration (one virtual user loop)

### Exporting Results

**JSON output:**
```bash
k6 run --out json=results.json scripts/smoke-test.js
```

**InfluxDB (for Grafana dashboards):**
```bash
k6 run --out influxdb=http://localhost:8086/k6 scripts/smoke-test.js
```

**Cloud (k6 Cloud service):**
```bash
k6 run --out cloud scripts/smoke-test.js
```

## Environment Variables

All scripts support the following environment variables:

- `BASE_URL`: Base URL of the application (default: `http://localhost:8080`)
- `ADMIN_EMAIL`: Admin user email for CRUD tests (default: `admin@example.com`)
- `ADMIN_PASSWORD`: Admin user password (default: `password123`)

**Example:**
```bash
BASE_URL=http://staging.example.com \
ADMIN_EMAIL=admin@staging.com \
ADMIN_PASSWORD=securepassword \
k6 run scripts/user-crud-test.js
```

## Best Practices

1. **Start Small:** Always run smoke tests before larger load tests
2. **Use Test Data:** Load fixtures or generate test users before running CRUD tests
3. **Monitor Resources:** Watch CPU, memory, and database connections during tests
4. **Clean Up:** Remove test data after load testing
5. **Gradual Ramp-Up:** Use stages to gradually increase load
6. **Set Thresholds:** Define acceptable performance criteria
7. **Test Realistic Scenarios:** Mimic actual user behavior patterns

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Load Test

on:
  schedule:
    - cron: '0 2 * * *' # Daily at 2 AM
  workflow_dispatch:

jobs:
  load-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install k6
        run: |
          sudo gpg -k
          sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6
      
      - name: Start application
        run: docker-compose up -d
      
      - name: Wait for application
        run: sleep 10
      
      - name: Run smoke test
        run: cd loadtest && ./run-smoke-test.sh http://localhost:8080
      
      - name: Run auth test
        run: cd loadtest && ./run-auth-test.sh http://localhost:8080
```

## Troubleshooting

### k6 is not installed
```bash
# Install k6 using your package manager
brew install k6  # macOS
sudo apt install k6  # Ubuntu/Debian
```

### Connection refused errors
- Make sure the application is running
- Check the `BASE_URL` is correct
- Verify firewalls aren't blocking connections

### High error rates
- Check application logs for errors
- Verify database is running and accessible
- Ensure sufficient resources (CPU, memory)
- Reduce the number of virtual users

### Slow response times
- Check database performance
- Monitor application resource usage
- Review database query performance
- Consider scaling up resources

## Additional Resources

- [k6 Documentation](https://k6.io/docs/)
- [k6 Examples](https://k6.io/docs/examples/)
- [k6 Best Practices](https://k6.io/docs/testing-guides/test-types/)
- [Load Testing Best Practices](https://k6.io/docs/testing-guides/load-testing-best-practices/)

## See Also

- `/testdata/fixtures/` - Test fixtures for database testing
- `/pkg/testutil/` - Test utilities and helpers
- `/scripts/test.sh` - Unit and integration test script
