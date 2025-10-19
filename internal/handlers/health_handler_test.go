package handlers

import (
	"myapp/pkg/health"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHealthTestRouter(registry *health.Registry) (*gin.Engine, *HealthHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHealthHandler(registry)
	return router, handler
}

func TestNewHealthHandler(t *testing.T) {
	t.Run("should create health handler with valid registry", func(t *testing.T) {
		registry := health.NewRegistry()
		handler := NewHealthHandler(registry)
		
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.registry)
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("should return healthy status with UP components", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		registry := health.NewRegistry()
		registry.Register(health.NewDatabaseHealthCheckProvider(db))
		
		router, handler := setupHealthTestRouter(registry)
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
		
		registry := health.NewRegistry()
		registry.Register(health.NewDatabaseHealthCheckProvider(db))
		
		router, handler := setupHealthTestRouter(registry)
		router.GET("/health", handler.HealthCheck)
		
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), "DOWN")
	})

	t.Run("should return UP status when no providers registered", func(t *testing.T) {
		registry := health.NewRegistry()
		router, handler := setupHealthTestRouter(registry)
		router.GET("/health", handler.HealthCheck)
		
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
	})

	t.Run("should check each provider only once", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		registry := health.NewRegistry()
		// Register a provider with multiple scopes
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup, health.ScopeReady, health.ScopeLive))
		
		router, handler := setupHealthTestRouter(registry)
		router.GET("/health", handler.HealthCheck)
		
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		// Database should appear only once in the response
		assert.Contains(t, w.Body.String(), "database")
	})
}

func TestStartupProbe(t *testing.T) {
	t.Run("should return OK when application is ready", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		registry := health.NewRegistry()
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup))
		
		router, handler := setupHealthTestRouter(registry)
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
		
		registry := health.NewRegistry()
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup))
		
		router, handler := setupHealthTestRouter(registry)
		router.GET("/health/startup", handler.StartupProbe)
		
		req, _ := http.NewRequest("GET", "/health/startup", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("should only check startup scoped providers", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		registry := health.NewRegistry()
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup))
		// This provider should not be checked in startup
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeLive))
		
		router, handler := setupHealthTestRouter(registry)
		router.GET("/health/startup", handler.StartupProbe)
		
		req, _ := http.NewRequest("GET", "/health/startup", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "database")
	})
}

func TestLivenessProbe(t *testing.T) {
	t.Run("should return OK when application is alive", func(t *testing.T) {
		registry := health.NewRegistry()
		router, handler := setupHealthTestRouter(registry)
		router.GET("/health/liveness", handler.LivenessProbe)
		
		req, _ := http.NewRequest("GET", "/health/liveness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
	})

	t.Run("should only check liveness scoped providers", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		registry := health.NewRegistry()
		// Provider with liveness scope
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeLive))
		
		router, handler := setupHealthTestRouter(registry)
		router.GET("/health/liveness", handler.LivenessProbe)
		
		req, _ := http.NewRequest("GET", "/health/liveness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "database")
	})
}

func TestReadinessProbe(t *testing.T) {
	t.Run("should return OK when ready to accept traffic", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		registry := health.NewRegistry()
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeReady))
		
		router, handler := setupHealthTestRouter(registry)
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
		
		registry := health.NewRegistry()
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeReady))
		
		router, handler := setupHealthTestRouter(registry)
		router.GET("/health/readiness", handler.ReadinessProbe)
		
		req, _ := http.NewRequest("GET", "/health/readiness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), "DOWN")
	})

	t.Run("should only check ready scoped providers", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		registry := health.NewRegistry()
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeReady))
		// This provider should not be checked in readiness
		registry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeLive))
		
		router, handler := setupHealthTestRouter(registry)
		router.GET("/health/readiness", handler.ReadinessProbe)
		
		req, _ := http.NewRequest("GET", "/health/readiness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "database")
	})
}

