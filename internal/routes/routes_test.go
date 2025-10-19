package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Use in-memory SQLite for testing
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	logger := zap.NewNop()
	
	SetupRoutes(router, db, logger, "test-secret")
	return router
}

func TestMetricsEndpoint(t *testing.T) {
	router := setupTestRouter()

	t.Run("should expose metrics in Prometheus format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/metrics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		body := w.Body.String()
		
		// Verify Prometheus format with TYPE and HELP comments
		assert.Contains(t, body, "# TYPE", "Should contain Prometheus TYPE declarations")
		assert.Contains(t, body, "# HELP", "Should contain Prometheus HELP descriptions")
	})

	t.Run("should expose runtime.MemStats metrics", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/metrics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		body := w.Body.String()
		
		// Verify key runtime.MemStats metrics are exposed
		expectedMemStatsMetrics := []string{
			"go_memstats_alloc_bytes",
			"go_memstats_sys_bytes",
			"go_memstats_heap_alloc_bytes",
			"go_memstats_heap_sys_bytes",
			"go_memstats_heap_idle_bytes",
			"go_memstats_heap_inuse_bytes",
			"go_memstats_heap_released_bytes",
			"go_memstats_heap_objects",
			"go_memstats_mallocs_total",
			"go_memstats_frees_total",
			"go_memstats_gc_sys_bytes",
		}
		
		for _, metric := range expectedMemStatsMetrics {
			assert.Contains(t, body, metric, "Should expose %s metric", metric)
		}
	})

	t.Run("should expose HTTP request metrics", func(t *testing.T) {
		// Make a request to generate metrics
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Now check metrics
		req, _ = http.NewRequest("GET", "/metrics", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		body := w.Body.String()
		
		// Verify custom HTTP metrics are exposed
		assert.Contains(t, body, "http_requests_total", "Should expose http_requests_total metric")
		assert.Contains(t, body, "http_request_duration_seconds", "Should expose http_request_duration_seconds metric")
	})

	t.Run("should expose Go runtime metrics", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/metrics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		body := w.Body.String()
		
		// Verify Go runtime metrics
		goRuntimeMetrics := []string{
			"go_goroutines",
			"go_threads",
			"go_gc_duration_seconds",
			"go_info",
		}
		
		for _, metric := range goRuntimeMetrics {
			assert.Contains(t, body, metric, "Should expose %s metric", metric)
		}
	})

	t.Run("metrics should be in valid Prometheus exposition format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/metrics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		body := w.Body.String()
		lines := strings.Split(body, "\n")
		
		// Verify format rules
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			
			// Comments start with #
			if strings.HasPrefix(line, "#") {
				// TYPE or HELP comments
				assert.True(t, 
					strings.Contains(line, "# TYPE") || strings.Contains(line, "# HELP"),
					"Comment lines should be TYPE or HELP declarations",
				)
			} else {
				// Metric lines should have metric name, optional labels, and value
				// Format: metric_name{label="value"} value
				// Or: metric_name value
				parts := strings.Fields(line)
				if len(parts) >= 1 {
					// First part should be metric name or metric with labels
					assert.True(t, 
						strings.Contains(parts[0], "_") || strings.Contains(parts[0], "{"),
						"Metric lines should contain metric name with underscores or labels",
					)
				}
			}
		}
	})
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	t.Run("should return healthy status with components", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
		assert.Contains(t, w.Body.String(), "database")
	})
}

func TestHealthStartupEndpoint(t *testing.T) {
	router := setupTestRouter()

	t.Run("should return startup probe status", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/startup", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
	})
}

func TestHealthLivenessEndpoint(t *testing.T) {
	router := setupTestRouter()

	t.Run("should return liveness probe status", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/liveness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
	})
}

func TestHealthReadinessEndpoint(t *testing.T) {
	router := setupTestRouter()

	t.Run("should return readiness probe status", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/readiness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "UP")
		assert.Contains(t, w.Body.String(), "database")
	})
}
