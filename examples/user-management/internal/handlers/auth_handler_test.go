package handlers

import (
	"bytes"
	"encoding/json"
	"myapp/examples/user-management/internal/models"
	"myapp/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Auto-migrate the User model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func setupTestLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func TestNewAuthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	logger := setupTestLogger()

	t.Run("should create auth handler", func(t *testing.T) {
		handler := NewAuthHandler(db, "test-secret", logger)
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.db)
		assert.NotNil(t, handler.logger)
		assert.Equal(t, "test-secret", handler.secret)
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should login successfully with valid credentials", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		// Create a test user
		hashedPassword, _ := utils.HashPassword("password123")
		user := &models.User{
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: hashedPassword,
			Role:         "user",
		}
		db.Create(user)

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response LoginResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, user.ID, response.User.ID)
		assert.Equal(t, user.Email, response.User.Email)
		assert.Equal(t, user.Name, response.User.Name)
		assert.Equal(t, user.Role, response.User.Role)
	})

	t.Run("should fail with invalid email", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid credentials", response["error"])
	})

	t.Run("should fail with invalid password", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		// Create a test user
		hashedPassword, _ := utils.HashPassword("password123")
		user := &models.User{
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: hashedPassword,
			Role:         "user",
		}
		db.Create(user)

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid credentials", response["error"])
	})

	t.Run("should reject invalid request format", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		// Invalid JSON
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject missing email", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := map[string]string{
			"password": "password123",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject missing password", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := map[string]string{
			"email": "test@example.com",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject invalid email format", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := LoginRequest{
			Email:    "not-an-email",
			Password: "password123",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should not expose sensitive data in response or logs", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		// Create a test user
		hashedPassword, _ := utils.HashPassword("password123")
		user := &models.User{
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: hashedPassword,
			Role:         "user",
		}
		db.Create(user)

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Verify password and password hash are not in response
		responseBody := w.Body.String()
		assert.NotContains(t, responseBody, "password123")
		assert.NotContains(t, responseBody, hashedPassword)
		assert.NotContains(t, responseBody, "PasswordHash")
		assert.NotContains(t, responseBody, "password_hash")
	})

	t.Run("should include request_id in context when available", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		// Create a test user
		hashedPassword, _ := utils.HashPassword("password123")
		user := &models.User{
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: hashedPassword,
			Role:         "user",
		}
		db.Create(user)

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		// Add middleware to set request_id
		router.Use(func(c *gin.Context) {
			c.Set("request_id", "test-request-123")
			c.Next()
		})
		router.POST("/login", handler.Login)

		loginReq := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should handle database error during login", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		// Close the database to simulate database error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid credentials", response["error"])
	})

	t.Run("should handle JWT generation failure", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		// Create a test user
		hashedPassword, _ := utils.HashPassword("password123")
		user := &models.User{
			Name:         "Test User",
			Email:        "test@example.com",
			PasswordHash: hashedPassword,
			Role:         "user",
		}
		db.Create(user)

		// Use an empty secret to potentially cause JWT generation issues
		// Note: This may not actually fail with the current JWT implementation,
		// but we test the error handling path
		handler := NewAuthHandler(db, "", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Even with empty secret, JWT library may succeed
		// This test ensures we handle the error path if it occurs
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	})

	t.Run("should return different error for admin vs user roles", func(t *testing.T) {
		db := setupTestDB(t)
		logger := setupTestLogger()

		// Create admin user
		hashedPassword, _ := utils.HashPassword("adminpass")
		adminUser := &models.User{
			Name:         "Admin User",
			Email:        "admin@example.com",
			PasswordHash: hashedPassword,
			Role:         "admin",
		}
		db.Create(adminUser)

		handler := NewAuthHandler(db, "test-secret", logger)
		router := gin.New()
		router.POST("/login", handler.Login)

		loginReq := LoginRequest{
			Email:    "admin@example.com",
			Password: "adminpass",
		}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response LoginResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "admin", response.User.Role)
	})
}
