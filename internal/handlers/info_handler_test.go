package handlers

import (
	"encoding/json"
	"myapp/pkg/info"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockInfoProvider is a simple mock implementation for testing
type mockInfoProvider struct {
	name string
	data map[string]interface{}
}

func (m *mockInfoProvider) Name() string {
	return m.name
}

func (m *mockInfoProvider) Info() (map[string]interface{}, error) {
	return m.data, nil
}

func TestNewInfoHandler(t *testing.T) {
	t.Run("should create info handler", func(t *testing.T) {
		registry := info.NewRegistry()
		handler := NewInfoHandler(registry)
		
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.registry)
		assert.NotNil(t, handler.limiter)
	})
}

func TestGetInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return empty info when no providers registered", func(t *testing.T) {
		registry := info.NewRegistry()
		handler := NewInfoHandler(registry)
		
		router := gin.New()
		router.GET("/info", handler.GetInfo)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Empty(t, response)
	})

	t.Run("should return aggregated info from single provider", func(t *testing.T) {
		registry := info.NewRegistry()
		provider := &mockInfoProvider{
			name: "build",
			data: map[string]interface{}{
				"version": "1.0.0",
				"commit":  "abc123",
			},
		}
		registry.Register(provider)
		
		handler := NewInfoHandler(registry)
		router := gin.New()
		router.GET("/info", handler.GetInfo)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Contains(t, response, "build")
		buildInfo := response["build"].(map[string]interface{})
		assert.Equal(t, "1.0.0", buildInfo["version"])
		assert.Equal(t, "abc123", buildInfo["commit"])
	})

	t.Run("should return aggregated info from multiple providers", func(t *testing.T) {
		registry := info.NewRegistry()
		
		buildProvider := &mockInfoProvider{
			name: "build",
			data: map[string]interface{}{
				"version": "1.0.0",
			},
		}
		
		statsProvider := &mockInfoProvider{
			name: "stats",
			data: map[string]interface{}{
				"total": 42,
			},
		}
		
		registry.Register(buildProvider)
		registry.Register(statsProvider)
		
		handler := NewInfoHandler(registry)
		router := gin.New()
		router.GET("/info", handler.GetInfo)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Len(t, response, 2)
		assert.Contains(t, response, "build")
		assert.Contains(t, response, "stats")
	})

	t.Run("should enforce rate limit bulkhead protection", func(t *testing.T) {
		registry := info.NewRegistry()
		handler := NewInfoHandler(registry)
		
		router := gin.New()
		router.GET("/info", handler.GetInfo)
		
		// Make requests up to the burst limit (20)
		successCount := 0
		for i := 0; i < 25; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/info", nil)
			router.ServeHTTP(w, req)
			
			if w.Code == http.StatusOK {
				successCount++
			}
		}
		
		// Should have some successful requests but not all 25
		assert.Greater(t, successCount, 0, "Should have some successful requests")
		assert.Less(t, successCount, 25, "Should block some requests due to rate limit")
	})

	t.Run("should return 429 when rate limit exceeded", func(t *testing.T) {
		registry := info.NewRegistry()
		handler := NewInfoHandler(registry)
		
		router := gin.New()
		router.GET("/info", handler.GetInfo)
		
		// Exhaust the rate limit
		for i := 0; i < 25; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/info", nil)
			router.ServeHTTP(w, req)
		}
		
		// Wait a tiny bit to ensure limiter state is settled
		time.Sleep(10 * time.Millisecond)
		
		// Next request should be rate limited
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info", nil)
		router.ServeHTTP(w, req)
		
		// Should eventually hit 429
		if w.Code == http.StatusTooManyRequests {
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["error"], "rate limit exceeded")
		}
	})
}

