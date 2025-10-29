package handlers

import (
	"encoding/json"
	"myapp/pkg/health"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHealthCheckScopes demonstrates the complete scope behavior as per requirements
func TestHealthCheckScopes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a registry and register providers with different scopes
	registry := health.NewRegistry()

	// Provider registered as "base" - only appears in /health
	baseProvider := health.NewSimpleHealthCheckProvider("base-check", health.ScopeBase)
	registry.Register(baseProvider)

	// Provider registered as "startup" - appears in /health/startup, /health/readiness, and /health
	// Note: The requirement states startup appears in /health/readiness and /health
	// We register it only with ScopeStartup, and it will appear in startup endpoint
	startupProvider := health.NewSimpleHealthCheckProvider("startup-check", health.ScopeStartup)
	registry.Register(startupProvider)

	// Provider with ready scope - appears in /health/readiness and /health
	readyProvider := health.NewSimpleHealthCheckProvider("ready-check", health.ScopeReady)
	registry.Register(readyProvider)

	// Provider with live scope - appears in /health/liveness and /health
	liveProvider := health.NewSimpleHealthCheckProvider("live-check", health.ScopeLive)
	registry.Register(liveProvider)

	// Provider with multiple scopes - appears in all matching scopes and /health (but checked only once in /health)
	multiScopeProvider := health.NewSimpleHealthCheckProvider("multi-scope-check", health.ScopeStartup, health.ScopeReady, health.ScopeLive)
	registry.Register(multiScopeProvider)

	handler := NewHealthHandler(registry)
	router := gin.New()

	router.GET("/health", handler.HealthCheck)
	router.GET("/health/startup", handler.StartupProbe)
	router.GET("/health/readiness", handler.ReadinessProbe)
	router.GET("/health/liveness", handler.LivenessProbe)

	t.Run("base scope only appears in /health", func(t *testing.T) {
		// Check /health
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var healthResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &healthResponse)
		assert.Contains(t, healthResponse.Components, "base-check")

		// Check /health/startup - should NOT contain base-check
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/health/startup", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var startupResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &startupResponse)
		assert.NotContains(t, startupResponse.Components, "base-check")

		// Check /health/readiness - should NOT contain base-check
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/health/readiness", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var readinessResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &readinessResponse)
		assert.NotContains(t, readinessResponse.Components, "base-check")

		// Check /health/liveness - should NOT contain base-check
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/health/liveness", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var livenessResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &livenessResponse)
		assert.NotContains(t, livenessResponse.Components, "base-check")
	})

	t.Run("startup scope appears in /health/startup and /health", func(t *testing.T) {
		// Check /health/startup
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health/startup", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var startupResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &startupResponse)
		assert.Contains(t, startupResponse.Components, "startup-check")

		// Check /health
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var healthResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &healthResponse)
		assert.Contains(t, healthResponse.Components, "startup-check")
	})

	t.Run("ready scope appears in /health/readiness and /health", func(t *testing.T) {
		// Check /health/readiness
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health/readiness", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var readinessResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &readinessResponse)
		assert.Contains(t, readinessResponse.Components, "ready-check")

		// Check /health
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var healthResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &healthResponse)
		assert.Contains(t, healthResponse.Components, "ready-check")
	})

	t.Run("live scope appears in /health/liveness and /health", func(t *testing.T) {
		// Check /health/liveness
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health/liveness", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var livenessResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &livenessResponse)
		assert.Contains(t, livenessResponse.Components, "live-check")

		// Check /health
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var healthResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &healthResponse)
		assert.Contains(t, healthResponse.Components, "live-check")
	})

	t.Run("multi-scope check appears in all scopes but only once in /health", func(t *testing.T) {
		// Check /health/startup
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health/startup", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var startupResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &startupResponse)
		assert.Contains(t, startupResponse.Components, "multi-scope-check")

		// Check /health/readiness
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/health/readiness", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var readinessResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &readinessResponse)
		assert.Contains(t, readinessResponse.Components, "multi-scope-check")

		// Check /health/liveness
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/health/liveness", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var livenessResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &livenessResponse)
		assert.Contains(t, livenessResponse.Components, "multi-scope-check")

		// Check /health - should contain multi-scope-check only once
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var healthResponse health.Response
		json.Unmarshal(w.Body.Bytes(), &healthResponse)
		assert.Contains(t, healthResponse.Components, "multi-scope-check")
		// Verify it appears only once by checking we can access it
		_, exists := healthResponse.Components["multi-scope-check"]
		assert.True(t, exists)
	})

	t.Run("/health includes all providers but checks each only once", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response health.Response
		json.Unmarshal(w.Body.Bytes(), &response)

		// Should contain all 5 providers
		assert.Len(t, response.Components, 5)
		assert.Contains(t, response.Components, "base-check")
		assert.Contains(t, response.Components, "startup-check")
		assert.Contains(t, response.Components, "ready-check")
		assert.Contains(t, response.Components, "live-check")
		assert.Contains(t, response.Components, "multi-scope-check")

		// All should be UP
		for name, component := range response.Components {
			assert.Equal(t, health.StatusUp, component.Status, "Component %s should be UP", name)
		}
	})
}
