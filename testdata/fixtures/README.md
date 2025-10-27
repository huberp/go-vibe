# Test Fixtures

This directory contains test fixtures for database testing in go-vibe.

## Directory Structure

```
testdata/fixtures/
├── minimal/          # Minimal fixture set (2 users: 1 user, 1 admin)
│   └── users.yml
└── full/            # Full fixture set (5 users with varied roles)
    └── users.yml
```

## Fixture Sets

### Minimal (`minimal/`)
Contains a minimal set of test data for quick tests:
- 2 users (1 regular user, 1 admin)
- All passwords are `password123` (bcrypt hashed)

**Users:**
- `test@example.com` - Test User (role: user)
- `admin@example.com` - Test Admin (role: admin)

### Full (`full/`)
Contains a comprehensive set of test data for integration tests:
- 5 users (3 regular users, 2 admins)
- All passwords are `password123` (bcrypt hashed)

**Users:**
- `john@example.com` - John Doe (role: user)
- `jane@example.com` - Jane Smith (role: admin)
- `bob@example.com` - Bob Johnson (role: user)
- `alice@example.com` - Alice Williams (role: admin)
- `charlie@example.com` - Charlie Brown (role: user)

## Usage

### Using Fixtures in Tests

```go
package mypackage

import (
    "myapp/pkg/testutil"
    "testing"
)

func TestWithFixtures(t *testing.T) {
    // Create test database
    tdb := testutil.NewTestDB(t)
    defer tdb.Cleanup()

    // Load minimal fixtures
    tdb.LoadFixtures(t, "../../testdata/fixtures/minimal")

    // Your test code here...
}
```

### Loading Different Fixture Sets

```go
// Load minimal fixtures (faster, good for unit tests)
tdb.LoadFixtures(t, "../../testdata/fixtures/minimal")

// Load full fixtures (comprehensive, good for integration tests)
tdb.LoadFixtures(t, "../../testdata/fixtures/full")
```

### Reloading Fixtures

If your test modifies data and you want to reset to the original state:

```go
tdb.LoadFixtures(t, "../../testdata/fixtures/minimal")

// Modify some data...
tdb.DB.Model(&models.User{}).Where("id = ?", 1).Update("name", "Modified")

// Reload fixtures to restore original data
tdb.ReloadFixtures(t)
```

### Truncating Tables

To completely clear all data:

```go
tdb.TruncateTables(t)
```

## Adding New Fixtures

1. Create a new YAML file in the appropriate directory (e.g., `users.yml`)
2. Follow the existing format:

```yaml
- id: 1
  name: User Name
  email: user@example.com
  password_hash: $2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYL1P6L7BqG
  role: user
  created_at: 2024-01-01T00:00:00Z
  updated_at: 2024-01-01T00:00:00Z
  deleted_at: null
```

3. Generate password hashes using bcrypt with cost factor 12

## Password Hash for Testing

The default password hash in fixtures is for `password123`:
```
$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYL1P6L7BqG
```

To generate a new hash in Go:
```go
import "golang.org/x/crypto/bcrypt"

hash, _ := bcrypt.GenerateFromPassword([]byte("yourpassword"), 12)
```

## Best Practices

1. **Use minimal fixtures** for unit tests to keep tests fast
2. **Use full fixtures** for integration tests that need comprehensive data
3. **Always clean up** after tests using `defer tdb.Cleanup()`
4. **Reload fixtures** if tests modify data and need a clean state
5. **Don't use fixtures** for simple mocked tests - use `gomock` instead

## See Also

- `/pkg/testutil/testdb.go` - Test database utilities
- `/pkg/testutil/testdb_test.go` - Example usage of test utilities
