# User Management Example

This example demonstrates how to build a complete user management microservice using the go-vibe template. It showcases authentication, authorization, CRUD operations, and best practices for building production-ready APIs.

## Overview

This example implements a full-featured user management system with:

- ✅ **User CRUD operations** - Create, Read, Update, Delete users
- ✅ **JWT Authentication** - Secure token-based authentication
- ✅ **Role-Based Access Control (RBAC)** - Admin and user roles
- ✅ **Password Security** - Bcrypt hashing with cost factor 12
- ✅ **Input Validation** - Comprehensive request validation
- ✅ **Database Migrations** - Version-controlled schema management
- ✅ **Complete Test Coverage** - Unit and integration tests
- ✅ **API Documentation** - OpenAPI/Swagger annotations

## What This Example Demonstrates

### 1. Domain Models

See `internal/models/user.go` for:
- GORM model definition with proper tags
- JSON serialization control (hiding sensitive fields)
- Validation rules
- Database constraints

### 2. Repository Pattern

See `internal/repository/` for:
- Interface-based design for testability
- PostgreSQL implementation
- Mock implementation for testing
- Context propagation

### 3. HTTP Handlers

See `internal/handlers/` for:
- Request/response DTOs
- Input validation with Gin binding
- Error handling patterns
- Swagger annotations for API documentation
- Authentication and authorization logic

### 4. Database Migrations

See `migrations/` for:
- SQL migration files (up/down)
- Schema version control
- Index optimization

### 5. Testing Strategy

All files include comprehensive tests demonstrating:
- Table-driven tests
- Mock usage
- Edge case coverage
- Error handling validation

## API Endpoints

### Authentication

| Method | Endpoint       | Auth | Description          |
|--------|----------------|------|----------------------|
| POST   | `/v1/login`    | None | Authenticate user    |
| POST   | `/v1/users`    | None | Create user (signup) |

### User Management

| Method | Endpoint         | Auth              | Description       |
|--------|------------------|-------------------|-------------------|
| GET    | `/v1/users`      | JWT (admin)       | List all users    |
| GET    | `/v1/users/{id}` | JWT (owner/admin) | Get user by ID    |
| PUT    | `/v1/users/{id}` | JWT (owner/admin) | Update user       |
| DELETE | `/v1/users/{id}` | JWT (admin)       | Delete user       |

## How to Use This Example

### Option 1: Reference Implementation

Use this as a reference when building your own domain models. Study the patterns and adapt them to your use case.

### Option 2: Copy and Modify

1. Copy the relevant files to your go-vibe project
2. Rename "User" to your domain entity (e.g., "Product", "Order")
3. Modify fields and validation rules for your use case
4. Update the migrations to match your schema
5. Adjust business logic in handlers

### Option 3: Extend This Example

1. Copy this entire directory to your project
2. Add additional domain models alongside User
3. Create relationships between models
4. Add more complex business logic

## Files in This Example

```
examples/user-management/
├── README.md                           # This file
├── internal/
│   ├── models/
│   │   ├── user.go                    # User domain model
│   │   └── user_test.go               # Model tests
│   ├── handlers/
│   │   ├── user_handler.go            # User CRUD handlers
│   │   ├── user_handler_test.go       # Handler tests
│   │   ├── auth_handler.go            # Login/authentication
│   │   └── auth_handler_test.go       # Auth tests
│   ├── repository/
│   │   ├── user_repository.go         # Repository interface
│   │   ├── postgres_user_repository.go # PostgreSQL implementation
│   │   └── user_repository_mock.go    # Mock for testing
│   └── user_stats_provider.go         # Info endpoint integration
├── migrations/
│   ├── 000001_create_users_table.up.sql   # Create users table
│   └── 000001_create_users_table.down.sql # Drop users table
└── scripts/
    ├── test-api.sh                    # API testing script (Linux/macOS)
    └── test-api.ps1                   # API testing script (Windows)
```

## Integration Points

### Adding to Routes

In your `internal/routes/routes.go`, integrate user management:

```go
import (
    userHandlers "path/to/examples/user-management/internal/handlers"
    userRepo "path/to/examples/user-management/internal/repository"
)

// Create repository
userRepository := userRepo.NewPostgresUserRepository(db)

// Create handlers
userHandler := userHandlers.NewUserHandler(userRepository)
authHandler := userHandlers.NewAuthHandler(db, jwtSecret, logger)

// Setup routes
v1 := router.Group("/v1")
{
    v1.POST("/login", authHandler.Login)
    v1.POST("/users", userHandler.CreateUser)
    
    protected := v1.Group("/")
    protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
    {
        admin := protected.Group("/")
        admin.Use(middleware.RequireRole("admin"))
        {
            admin.GET("/users", userHandler.GetUsers)
            admin.DELETE("/users/:id", userHandler.DeleteUser)
        }
        
        protected.GET("/users/:id", userHandler.GetUserByID)
        protected.PUT("/users/:id", userHandler.UpdateUser)
    }
}
```

### Running Migrations

```bash
# Linux/macOS
./scripts/migrate.sh up

# Windows PowerShell
.\scripts\migrate.ps1 up
```

### Adding User Count Metric

In your metrics middleware:

```go
import "path/to/examples/user-management/internal"

// Register user count collector
middleware.RegisterUserCountCollector(db)
```

## Key Patterns Demonstrated

### 1. Secure Password Handling

```go
// Hash password before storing
hashedPassword, err := utils.HashPassword(password)

// Never return password hash in responses
type User struct {
    PasswordHash string `json:"-" gorm:"not null"` // json:"-" prevents serialization
}
```

### 2. Role-Based Authorization

```go
// In middleware
if role != requiredRole {
    c.JSON(403, gin.H{"error": "Insufficient permissions"})
    c.Abort()
    return
}

// In handler
userRole := c.GetString("user_role")
if userRole != "admin" {
    c.JSON(403, gin.H{"error": "Admin access required"})
    return
}
```

### 3. Owner or Admin Pattern

```go
// Allow users to access/modify their own data, or admins to access anyone's
userID := c.GetUint("user_id")
userRole := c.GetString("user_role")

if userRole != "admin" && userID != requestedUserID {
    c.JSON(403, gin.H{"error": "Access denied"})
    return
}
```

### 4. Request Validation

```go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Role     string `json:"role" binding:"omitempty,oneof=user admin"`
}

if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(400, gin.H{"error": "Invalid request"})
    return
}
```

### 5. Consistent Error Responses

```go
// Not found
c.JSON(404, gin.H{"error": "User not found"})

// Validation error
c.JSON(400, gin.H{"error": "Invalid input"})

// Unauthorized
c.JSON(401, gin.H{"error": "Invalid credentials"})

// Forbidden
c.JSON(403, gin.H{"error": "Insufficient permissions"})

// Server error
c.JSON(500, gin.H{"error": "Internal server error"})
```

## Testing

Run the tests for this example:

```bash
# Test models
go test ./examples/user-management/internal/models/... -v

# Test handlers
go test ./examples/user-management/internal/handlers/... -v

# Test repository (requires database)
go test ./examples/user-management/internal/repository/... -v

# All tests with coverage
go test ./examples/user-management/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## API Testing

Use the provided scripts to test the API:

```bash
# Linux/macOS
cd examples/user-management/scripts
./test-api.sh

# Windows PowerShell
cd examples\user-management\scripts
.\test-api.ps1
```

## Security Considerations

This example demonstrates production-ready security practices:

1. **Password Security**
   - Bcrypt hashing with cost factor 12
   - Never log or return passwords
   - Minimum password length enforcement

2. **Authentication**
   - JWT tokens with HS256 algorithm
   - Token expiration (24 hours)
   - Secure secret key from environment

3. **Authorization**
   - Role-based access control
   - Owner-based permissions
   - Admin-only operations

4. **Input Validation**
   - Email format validation
   - Required field checking
   - Role enumeration (user/admin only)

5. **Database Security**
   - Parameterized queries via GORM
   - SQL injection prevention
   - Unique constraints on email

## Customization Guide

### Adding New Fields

1. Update `internal/models/user.go`:
```go
type User struct {
    gorm.Model
    Name         string `json:"name" gorm:"not null"`
    Email        string `json:"email" gorm:"uniqueIndex;not null"`
    PasswordHash string `json:"-" gorm:"not null"`
    Role         string `json:"role" gorm:"default:'user';not null"`
    PhoneNumber  string `json:"phone_number" gorm:"default:null"` // New field
}
```

2. Create migration:
```sql
-- migrations/000002_add_phone_number.up.sql
ALTER TABLE users ADD COLUMN phone_number VARCHAR(20);
```

3. Update request DTOs in handlers:
```go
type UpdateUserRequest struct {
    Name        string `json:"name" binding:"omitempty"`
    Email       string `json:"email" binding:"omitempty,email"`
    PhoneNumber string `json:"phone_number" binding:"omitempty"`
}
```

### Adding New Endpoints

1. Add handler method in `internal/handlers/user_handler.go`
2. Add route in your `routes.go`
3. Add Swagger annotations
4. Write tests
5. Update documentation

## Learning Resources

- **GORM Documentation**: https://gorm.io/docs/
- **Gin Framework**: https://gin-gonic.com/docs/
- **JWT Best Practices**: https://jwt.io/introduction
- **Go Testing**: https://go.dev/doc/tutorial/add-a-test

## License

This example is part of the go-vibe project and follows the same MIT license.

## Support

For questions or issues specific to this example:
1. Check the main project documentation
2. Review the test files for usage examples
3. Open an issue on the go-vibe repository
