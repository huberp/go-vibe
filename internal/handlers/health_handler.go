package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints for Kubernetes probes
type HealthHandler struct {
	db *gorm.DB
}

// HealthStatus represents the health status of a component
type HealthStatus string

const (
	StatusUP   HealthStatus = "UP"
	StatusDown HealthStatus = "DOWN"
)

// ComponentHealth represents the health of a single component
type ComponentHealth struct {
	Status  HealthStatus           `json:"status"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HealthResponse represents the overall health response
type HealthResponse struct {
	Status     HealthStatus               `json:"status"`
	Components map[string]ComponentHealth `json:"components,omitempty"`
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns overall health status with all component checks
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "All components healthy"
// @Failure 503 {object} HealthResponse "One or more components unhealthy"
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	components := make(map[string]ComponentHealth)
	overallStatus := StatusUP

	// Check database
	dbStatus, dbDetails := h.checkDatabase()
	components["database"] = ComponentHealth{
		Status:  dbStatus,
		Details: dbDetails,
	}
	
	if dbStatus == StatusDown {
		overallStatus = StatusDown
	}

	response := HealthResponse{
		Status:     overallStatus,
		Components: components,
	}

	statusCode := http.StatusOK
	if overallStatus == StatusDown {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// StartupProbe godoc
// @Summary Kubernetes startup probe
// @Description Indicates if the application has started successfully
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "Application started"
// @Failure 503 {object} HealthResponse "Application not started"
// @Router /health/startup [get]
func (h *HealthHandler) StartupProbe(c *gin.Context) {
	// Startup probe checks if the application has started
	// This includes basic database connectivity check
	dbStatus, dbDetails := h.checkDatabase()
	
	response := HealthResponse{
		Status: dbStatus,
		Components: map[string]ComponentHealth{
			"database": {
				Status:  dbStatus,
				Details: dbDetails,
			},
		},
	}

	statusCode := http.StatusOK
	if dbStatus == StatusDown {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// LivenessProbe godoc
// @Summary Kubernetes liveness probe
// @Description Indicates if the application is running and should not be restarted
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "Application is alive"
// @Router /health/liveness [get]
func (h *HealthHandler) LivenessProbe(c *gin.Context) {
	// Liveness probe is simple - if we can respond, we're alive
	// This should not check external dependencies to avoid unnecessary restarts
	response := HealthResponse{
		Status: StatusUP,
	}

	c.JSON(http.StatusOK, response)
}

// ReadinessProbe godoc
// @Summary Kubernetes readiness probe
// @Description Indicates if the application is ready to accept traffic
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "Application ready to accept traffic"
// @Failure 503 {object} HealthResponse "Application not ready"
// @Router /health/readiness [get]
func (h *HealthHandler) ReadinessProbe(c *gin.Context) {
	// Readiness probe checks if the application can serve traffic
	// This includes checking database and other critical dependencies
	components := make(map[string]ComponentHealth)
	overallStatus := StatusUP

	// Check database
	dbStatus, dbDetails := h.checkDatabase()
	components["database"] = ComponentHealth{
		Status:  dbStatus,
		Details: dbDetails,
	}
	
	if dbStatus == StatusDown {
		overallStatus = StatusDown
	}

	response := HealthResponse{
		Status:     overallStatus,
		Components: components,
	}

	statusCode := http.StatusOK
	if overallStatus == StatusDown {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// checkDatabase checks if the database connection is healthy
func (h *HealthHandler) checkDatabase() (HealthStatus, map[string]interface{}) {
	if h.db == nil {
		return StatusDown, map[string]interface{}{
			"error": "database not initialized",
		}
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		return StatusDown, map[string]interface{}{
			"error": err.Error(),
		}
	}

	// Ping database with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return StatusDown, map[string]interface{}{
			"error": err.Error(),
		}
	}

	// Get database stats
	stats := sqlDB.Stats()
	return StatusUP, map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
	}
}
