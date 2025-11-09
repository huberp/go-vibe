package handlers

import (
	"myapp/pkg/health"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints for Kubernetes probes
type HealthHandler struct {
	registry *health.Registry
}

// NewHealthHandler creates a new health handler with the given registry
func NewHealthHandler(registry *health.Registry) *HealthHandler {
	return &HealthHandler{
		registry: registry,
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns overall health status with all component checks
// @Tags health
// @Produce json
// @Success 200 {object} health.Response "All components healthy"
// @Failure 503 {object} health.Response "One or more components unhealthy"
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check all providers (scope = nil means check all, but each provider only once)
	statusCode, response := h.registry.BuildResponse(nil)
	c.JSON(statusCode, response)
}

// StartupProbe godoc
// @Summary Kubernetes startup probe
// @Description Indicates if the application has started successfully
// @Tags health
// @Produce json
// @Success 200 {object} health.Response "Application started"
// @Failure 503 {object} health.Response "Application not started"
// @Router /health/startup [get]
func (h *HealthHandler) StartupProbe(c *gin.Context) {
	// Check only startup scope providers
	scope := health.ScopeStartup
	statusCode, response := h.registry.BuildResponse(&scope)
	c.JSON(statusCode, response)
}

// LivenessProbe godoc
// @Summary Kubernetes liveness probe
// @Description Indicates if the application is running and should not be restarted
// @Tags health
// @Produce json
// @Success 200 {object} health.Response "Application is alive"
// @Router /health/liveness [get]
func (h *HealthHandler) LivenessProbe(c *gin.Context) {
	// Check only liveness scope providers
	scope := health.ScopeLive
	_, response := h.registry.BuildResponse(&scope)
	c.JSON(http.StatusOK, response)
}

// ReadinessProbe godoc
// @Summary Kubernetes readiness probe
// @Description Indicates if the application is ready to accept traffic
// @Tags health
// @Produce json
// @Success 200 {object} health.Response "Application ready to accept traffic"
// @Failure 503 {object} health.Response "Application not ready"
// @Router /health/readiness [get]
func (h *HealthHandler) ReadinessProbe(c *gin.Context) {
	// Check only readiness scope providers
	scope := health.ScopeReady
	statusCode, response := h.registry.BuildResponse(&scope)
	c.JSON(statusCode, response)
}
