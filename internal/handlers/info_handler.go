package handlers

import (
	"myapp/pkg/info"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InfoHandler handles requests to the /info endpoint.
type InfoHandler struct {
	registry *info.Registry
}

// NewInfoHandler creates a new InfoHandler with the given registry.
func NewInfoHandler(registry *info.Registry) *InfoHandler {
	return &InfoHandler{
		registry: registry,
	}
}

// GetInfo returns aggregated information from all registered providers.
// @Summary Get application information
// @Description Get aggregated information from all registered info providers
// @Tags info
// @Produce json
// @Success 200 {object} map[string]interface{} "Aggregated information"
// @Router /info [get]
func (h *InfoHandler) GetInfo(c *gin.Context) {
	info := h.registry.GetAll()
	c.JSON(http.StatusOK, info)
}
