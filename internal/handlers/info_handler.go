package handlers

import (
	"myapp/pkg/info"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// InfoHandler handles requests to the /info endpoint.
type InfoHandler struct {
	registry *info.Registry
	limiter  *rate.Limiter
	mu       sync.Mutex
}

// NewInfoHandler creates a new InfoHandler with the given registry.
// The handler includes built-in rate limiting (bulkhead pattern) to prevent DOS attacks.
func NewInfoHandler(registry *info.Registry) *InfoHandler {
	return &InfoHandler{
		registry: registry,
		// Bulkhead: Limit to 10 requests per second with burst of 20
		// This is more restrictive than global rate limit to protect against DOS
		limiter: rate.NewLimiter(10, 20),
	}
}

// GetInfo returns aggregated information from all registered providers.
// This endpoint has additional rate limiting (bulkhead) to prevent abuse.
// @Summary Get application information
// @Description Get aggregated information from all registered info providers
// @Tags info
// @Produce json
// @Success 200 {object} map[string]interface{} "Aggregated information"
// @Failure 429 {object} map[string]string "Rate limit exceeded"
// @Router /info [get]
func (h *InfoHandler) GetInfo(c *gin.Context) {
	// Bulkhead rate limiting - protect against DOS attacks
	h.mu.Lock()
	allowed := h.limiter.Allow()
	h.mu.Unlock()

	if !allowed {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "rate limit exceeded for info endpoint",
		})
		return
	}

	info := h.registry.GetAll()
	c.JSON(http.StatusOK, info)
}
