package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHealthTestRouter(db *gorm.DB) (*gin.Engine, *HealthHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHealthHandler(db)
	return router, handler
}

func TestNewHealthHandler(t *testing.T) {
	t.Run("should create health handler with valid database", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		handler := NewHealthHandler(db)
		
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.db)
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("should return healthy status with UP components", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		router, handler := setupHealthTestRouter(db)
		router.GET("/health", handler.HealthCheck)
		
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
		assert.Contains(t, w.Body.String(), "database")
	})
	
	t.Run("should return unhealthy status with database error", func(t *testing.T) {
		// Create a database and close it to simulate connection error
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		sqlDB, _ := db.DB()
		sqlDB.Close()
		
		router, handler := setupHealthTestRouter(db)
		router.GET("/health", handler.HealthCheck)
		
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), "DOWN")
	})
}

func TestStartupProbe(t *testing.T) {
	t.Run("should return OK when application is ready", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		router, handler := setupHealthTestRouter(db)
		router.GET("/health/startup", handler.StartupProbe)
		
		req, _ := http.NewRequest("GET", "/health/startup", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
	})
	
	t.Run("should return service unavailable when database not ready", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		sqlDB, _ := db.DB()
		sqlDB.Close()
		
		router, handler := setupHealthTestRouter(db)
		router.GET("/health/startup", handler.StartupProbe)
		
		req, _ := http.NewRequest("GET", "/health/startup", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})
}

func TestLivenessProbe(t *testing.T) {
	t.Run("should return OK when application is alive", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		router, handler := setupHealthTestRouter(db)
		router.GET("/health/liveness", handler.LivenessProbe)
		
		req, _ := http.NewRequest("GET", "/health/liveness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
	})
}

func TestReadinessProbe(t *testing.T) {
	t.Run("should return OK when ready to accept traffic", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		router, handler := setupHealthTestRouter(db)
		router.GET("/health/readiness", handler.ReadinessProbe)
		
		req, _ := http.NewRequest("GET", "/health/readiness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
		assert.Contains(t, w.Body.String(), "database")
	})
	
	t.Run("should return service unavailable when database not ready", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		sqlDB, _ := db.DB()
		sqlDB.Close()
		
		router, handler := setupHealthTestRouter(db)
		router.GET("/health/readiness", handler.ReadinessProbe)
		
		req, _ := http.NewRequest("GET", "/health/readiness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), "DOWN")
	})
}

func TestCheckDatabase(t *testing.T) {
	t.Run("should return UP status for healthy database", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		handler := NewHealthHandler(db)
		
		status, details := handler.checkDatabase()
		
		assert.Equal(t, StatusUP, status)
		assert.NotNil(t, details)
	})
	
	t.Run("should return DOWN status for closed database", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		sqlDB, _ := db.DB()
		sqlDB.Close()
		
		handler := NewHealthHandler(db)
		status, details := handler.checkDatabase()
		
		assert.Equal(t, StatusDown, status)
		assert.NotNil(t, details)
		assert.Contains(t, details["error"], "sql: database is closed")
	})
	
	t.Run("should handle nil database gracefully", func(t *testing.T) {
		handler := &HealthHandler{db: nil}
		
		status, details := handler.checkDatabase()
		
		assert.Equal(t, StatusDown, status)
		assert.NotNil(t, details)
		assert.Contains(t, details["error"], "database not initialized")
	})
}
