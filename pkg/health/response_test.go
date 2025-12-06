package health

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentHealth(t *testing.T) {
	t.Run("should create component health with status and details", func(t *testing.T) {
		component := ComponentHealth{
			Status: StatusUp,
			Details: map[string]any{
				"version": "1.0.0",
			},
		}

		assert.Equal(t, StatusUp, component.Status)
		assert.Equal(t, "1.0.0", component.Details["version"])
	})

	t.Run("should marshal to JSON correctly", func(t *testing.T) {
		component := ComponentHealth{
			Status: StatusUp,
			Details: map[string]any{
				"message": "healthy",
			},
		}

		jsonData, err := json.Marshal(component)
		assert.NoError(t, err)
		assert.Contains(t, string(jsonData), `"status":"UP"`)
		assert.Contains(t, string(jsonData), `"message":"healthy"`)
	})

	t.Run("should omit empty details in JSON", func(t *testing.T) {
		component := ComponentHealth{
			Status: StatusUp,
		}

		jsonData, err := json.Marshal(component)
		assert.NoError(t, err)
		assert.NotContains(t, string(jsonData), "details")
	})
}

func TestResponse(t *testing.T) {
	t.Run("should create response with status and components", func(t *testing.T) {
		response := Response{
			Status: StatusUp,
			Components: map[string]ComponentHealth{
				"database": {
					Status: StatusUp,
					Details: map[string]any{
						"connections": 10,
					},
				},
			},
		}

		assert.Equal(t, StatusUp, response.Status)
		assert.Len(t, response.Components, 1)
		assert.Equal(t, StatusUp, response.Components["database"].Status)
	})

	t.Run("should marshal to JSON correctly", func(t *testing.T) {
		response := Response{
			Status: StatusUp,
			Components: map[string]ComponentHealth{
				"database": {
					Status: StatusUp,
				},
			},
		}

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)
		assert.Contains(t, string(jsonData), `"status":"UP"`)
		assert.Contains(t, string(jsonData), `"database"`)
	})

	t.Run("should omit empty components in JSON", func(t *testing.T) {
		response := Response{
			Status: StatusUp,
		}

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)
		assert.NotContains(t, string(jsonData), "components")
	})

	t.Run("should handle multiple components", func(t *testing.T) {
		response := Response{
			Status: StatusUp,
			Components: map[string]ComponentHealth{
				"database": {
					Status: StatusUp,
				},
				"cache": {
					Status: StatusUp,
				},
			},
		}

		assert.Len(t, response.Components, 2)
		assert.Contains(t, response.Components, "database")
		assert.Contains(t, response.Components, "cache")
	})

	t.Run("should reflect DOWN status when any component is down", func(t *testing.T) {
		components := map[string]*CheckResult{
			"database": {
				Status: StatusUp,
			},
			"cache": {
				Status: StatusDown,
			},
		}

		overallStatus := OverallStatus(components)
		assert.Equal(t, StatusDown, overallStatus)
	})
}

func TestResponseIntegration(t *testing.T) {
	t.Run("should work with registry and providers", func(t *testing.T) {
		// Create registry and register providers
		registry := NewRegistry()
		registry.Register(NewSimpleHealthCheckProvider("app", ScopeLive))

		// Check health
		checkResults := registry.Check(nil)

		// Build response
		components := make(map[string]ComponentHealth)
		for name, result := range checkResults {
			components[name] = ComponentHealth{
				Status:  result.Status,
				Details: result.Details,
			}
		}

		response := Response{
			Status:     OverallStatus(checkResults),
			Components: components,
		}

		// Verify
		assert.Equal(t, StatusUp, response.Status)
		assert.Len(t, response.Components, 1)
		assert.Equal(t, StatusUp, response.Components["app"].Status)
	})

	t.Run("should serialize complete health check response", func(t *testing.T) {
		// Create a realistic health response
		response := Response{
			Status: StatusUp,
			Components: map[string]ComponentHealth{
				"database": {
					Status: StatusUp,
					Details: map[string]any{
						"connections": 5,
						"max":         100,
					},
				},
				"app": {
					Status: StatusUp,
					Details: map[string]any{
						"timestamp": "2024-01-01T00:00:00Z",
					},
				},
			},
		}

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		// Verify JSON structure
		var parsed map[string]any
		err = json.Unmarshal(jsonData, &parsed)
		assert.NoError(t, err)
		assert.Equal(t, "UP", parsed["status"])

		components, ok := parsed["components"].(map[string]any)
		assert.True(t, ok)
		assert.Len(t, components, 2)
	})
}
