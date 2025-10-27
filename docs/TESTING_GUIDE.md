# Test Fixtures and Load Testing Guide

This guide covers the test fixtures and load testing infrastructure for go-vibe.

## Overview

The project now includes comprehensive testing capabilities:

1. **Test Fixtures** - YAML-based database fixtures for integration testing
2. **Load Testing** - k6-based load testing scripts for performance validation

## Table of Contents

- [Test Fixtures](#test-fixtures)
  - [Quick Start](#test-fixtures-quick-start)
  - [Usage Examples](#test-fixtures-usage-examples)
  - [Available Fixture Sets](#available-fixture-sets)
- [Load Testing](#load-testing)
  - [Quick Start](#load-testing-quick-start)
  - [Available Tests](#available-load-tests)
  - [Test Data Generation](#test-data-generation)
- [Integration with CI/CD](#integration-with-cicd)
- [Best Practices](#best-practices)

## Test Fixtures

Test fixtures provide a way to load predefined test data into your database for integration testing. This is useful for:

- Integration tests that need realistic database state
- Testing complex queries and relationships
- Ensuring consistent test data across test runs
- Load testing with known data sets

### Test Fixtures Quick Start

1. **Import the test utility package:**

```go
import "myapp/pkg/testutil"
```

2. **Create a test database and load fixtures:**

```go
func TestWithFixtures(t *testing.T) {
    // Setup test database
    tdb := testutil.NewTestDB(t)
    defer tdb.Cleanup()

    // Load fixtures
    tdb.LoadFixtures(t, "../../testdata/fixtures/minimal")

    // Your test code here
    var users []models.User
    tdb.DB.Find(&users)
    assert.Len(t, users, 2) // minimal fixture has 2 users
}
```

### Test Fixtures Usage Examples

**Example 1: Testing Repository with Fixtures**

```go
func TestUserRepository(t *testing.T) {
    tdb := testutil.NewTestDB(t)
    defer tdb.Cleanup()
    tdb.LoadFixtures(t, "../../testdata/fixtures/full")

    repo := repository.NewPostgresUserRepository(tdb.DB)

    // Test finding all users
    users, err := repo.FindAll(context.Background())
    assert.NoError(t, err)
    assert.Len(t, users, 5) // full fixture has 5 users

    // Test finding specific user
    user, err := repo.FindByEmail(context.Background(), "john@example.com")
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
}
```

**Example 2: Testing with Modified Data**

```go
func TestDataModification(t *testing.T) {
    tdb := testutil.NewTestDB(t)
    defer tdb.Cleanup()
    tdb.LoadFixtures(t, "../../testdata/fixtures/minimal")

    // Modify data
    tdb.DB.Model(&models.User{}).Where("id = ?", 1).Update("name", "Modified")

    // Reload fixtures to restore original state
    tdb.ReloadFixtures(t)

    // Verify original data is restored
    var user models.User
    tdb.DB.First(&user, 1)
    assert.Equal(t, "Test User", user.Name) // Original name restored
}
```

**Example 3: Using Transactions for Test Isolation**

```go
func TestWithTransaction(t *testing.T) {
    tdb := testutil.NewTestDB(t)
    defer tdb.Cleanup()
    tdb.LoadFixtures(t, "../../testdata/fixtures/minimal")

    // Start transaction
    tx := tdb.BeginTransaction()
    defer tx.Rollback()

    // Make changes in transaction
    newUser := &models.User{Name: "New User", Email: "new@example.com"}
    tx.Create(newUser)

    // Changes visible in transaction
    var count int64
    tx.Model(&models.User{}).Count(&count)
    assert.Equal(t, int64(3), count) // 2 from fixtures + 1 new

    // Rollback - changes not persisted
    tx.Rollback()

    // Verify no changes in main DB
    tdb.DB.Model(&models.User{}).Count(&count)
    assert.Equal(t, int64(2), count) // Only original 2 from fixtures
}
```

### Available Fixture Sets

#### Minimal Fixtures (`testdata/fixtures/minimal/`)

- **Users:** 2 (1 user, 1 admin)
- **Purpose:** Quick unit tests, fast setup
- **Credentials:**
  - `test@example.com` / `password123` (role: user)
  - `admin@example.com` / `password123` (role: admin)

#### Full Fixtures (`testdata/fixtures/full/`)

- **Users:** 5 (3 users, 2 admins)
- **Purpose:** Integration tests, comprehensive scenarios
- **Credentials:**
  - `john@example.com` / `password123` (role: user)
  - `jane@example.com` / `password123` (role: admin)
  - `bob@example.com` / `password123` (role: user)
  - `alice@example.com` / `password123` (role: admin)
  - `charlie@example.com` / `password123` (role: user)

For more details, see [`testdata/fixtures/README.md`](testdata/fixtures/README.md)

## Load Testing

Load testing validates that the application performs well under realistic and stress conditions.

### Load Testing Quick Start

1. **Install k6:**

```bash
# macOS
brew install k6

# Ubuntu/Debian
sudo apt-get install k6

# Windows
choco install k6
```

2. **Start your application:**

```bash
./server
# or
docker-compose up
```

3. **Run a smoke test:**

```bash
cd loadtest
./run-smoke-test.sh http://localhost:8080
```

### Available Load Tests

#### 1. Smoke Test

**Purpose:** Quick validation that the system works.

**Configuration:**
- 1 virtual user
- 30 seconds duration
- Tests `/health` and `/info` endpoints

**Run:**
```bash
cd loadtest
./run-smoke-test.sh http://localhost:8080
```

#### 2. Authentication Load Test

**Purpose:** Test login and authenticated endpoints under realistic load.

**Configuration:**
- Ramps from 0 to 10 virtual users
- 2 minutes total duration
- Tests login, get users, get user by ID

**Run:**
```bash
cd loadtest
./run-auth-test.sh http://localhost:8080
```

#### 3. Stress Test

**Purpose:** Push the system beyond normal capacity.

**Configuration:**
- Ramps up to 100 virtual users
- 12 minutes total duration
- Tests health, metrics, and info endpoints

**Run:**
```bash
cd loadtest
./run-stress-test.sh http://localhost:8080
```

#### 4. User CRUD Test

**Purpose:** Test full user lifecycle operations.

**Configuration:**
- 20 virtual users
- 5 minutes duration
- Tests create, login, read operations

**Run:**
```bash
cd loadtest
k6 run --env BASE_URL=http://localhost:8080 scripts/user-crud-test.js
```

### Test Data Generation

For load testing with many users, use the data generation script:

**Generate Test Users:**
```bash
cd loadtest/data/generate
export DATABASE_URL="postgres://user:password@localhost:5432/myapp?sslmode=disable"

# Generate 100 test users
go run main.go -db "$DATABASE_URL" -users 100

# Generate 1000 test users
go run main.go -db "$DATABASE_URL" -users 1000
```

**Clean Up Test Users:**
```bash
cd loadtest/data/cleanup
go run main.go -db "$DATABASE_URL"
```

For more details, see [`loadtest/README.md`](loadtest/README.md)

## Integration with CI/CD

### Using Fixtures in CI

Test fixtures work automatically in CI environments since they use SQLite in-memory databases by default:

```yaml
# .github/workflows/test.yml
- name: Run tests with fixtures
  run: go test ./... -v
```

For PostgreSQL integration tests in CI:

```yaml
- name: Setup PostgreSQL
  run: |
    docker run -d -p 5432:5432 \
      -e POSTGRES_PASSWORD=password \
      postgres:13

- name: Run integration tests
  run: |
    export TEST_DATABASE_URL="postgres://postgres:password@localhost/postgres?sslmode=disable"
    go test ./... -v -tags=integration
```

### Load Testing in CI

Add scheduled load tests to monitor performance over time:

```yaml
# .github/workflows/load-test.yml
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
          sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6
      
      - name: Start services
        run: docker-compose up -d
      
      - name: Wait for services
        run: sleep 10
      
      - name: Run smoke test
        run: cd loadtest && ./run-smoke-test.sh http://localhost:8080
```

## Best Practices

### Test Fixtures

1. ✅ **Use minimal fixtures** for unit tests to keep tests fast
2. ✅ **Use full fixtures** for integration tests requiring comprehensive data
3. ✅ **Always clean up** using `defer tdb.Cleanup()`
4. ✅ **Reload fixtures** when tests modify data: `tdb.ReloadFixtures(t)`
5. ✅ **Use transactions** for test isolation when testing data modifications
6. ❌ **Don't use fixtures** for simple mocked tests - use `gomock` instead
7. ❌ **Don't commit** sensitive data in fixtures - use hashed passwords

### Load Testing

1. ✅ **Start with smoke tests** before running larger load tests
2. ✅ **Load test data** before running CRUD tests
3. ✅ **Monitor resources** (CPU, memory, database) during tests
4. ✅ **Clean up test data** after load testing
5. ✅ **Use realistic scenarios** that mimic actual user behavior
6. ✅ **Set thresholds** to automatically fail on performance regressions
7. ❌ **Don't run stress tests** against production environments
8. ❌ **Don't skip cleanup** - always remove test data

## Troubleshooting

### Fixtures Not Loading

**Problem:** Fixtures fail to load with "no such table" error.

**Solution:**
- Ensure `AutoMigrate()` is called before loading fixtures
- Check that fixture YAML files are in the correct directory
- Verify file structure matches database schema

### Load Test Connection Errors

**Problem:** k6 reports connection refused errors.

**Solution:**
- Verify application is running: `curl http://localhost:8080/health`
- Check correct BASE_URL is specified
- Ensure no firewall is blocking connections

### High Error Rates in Load Tests

**Problem:** Load test shows high percentage of failed requests.

**Solution:**
- Check application logs for errors
- Verify database is running and accessible
- Reduce number of virtual users to find breaking point
- Check resource limits (connections, memory, CPU)

## Additional Resources

- [Test Fixtures README](testdata/fixtures/README.md) - Detailed fixture documentation
- [Load Testing README](loadtest/README.md) - Comprehensive load testing guide
- [Test Utilities](pkg/testutil/) - TestDB implementation and examples
- [Example Tests](pkg/testutil/examples/) - Practical examples using fixtures

## Summary

This implementation provides:

✅ **YAML-based test fixtures** for database testing  
✅ **Fixture management utilities** with load/reload/cleanup  
✅ **Transaction-based test isolation**  
✅ **k6 load testing scripts** for smoke, load, and stress testing  
✅ **Test data generators** for large-scale load testing  
✅ **Automation scripts** for easy test execution  
✅ **Comprehensive documentation** and examples  
✅ **CI/CD integration** examples  

All tests are passing, and the infrastructure is ready for production use! 🚀
