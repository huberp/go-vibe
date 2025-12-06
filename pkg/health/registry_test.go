package health

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockHealthCheckProvider is a simple mock implementation for testing
type mockHealthCheckProvider struct {
	name   string
	result *CheckResult
	err    error
	scopes []Scope
}

func (m *mockHealthCheckProvider) Name() string {
	return m.name
}

func (m *mockHealthCheckProvider) Check() (*CheckResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func (m *mockHealthCheckProvider) Scopes() []Scope {
	return m.scopes
}

func TestNewRegistry(t *testing.T) {
	t.Run("should create empty registry", func(t *testing.T) {
		registry := NewRegistry()
		assert.NotNil(t, registry)
		result := registry.Check(nil)
		assert.Empty(t, result)
	})
}

func TestRegister(t *testing.T) {
	t.Run("should register provider", func(t *testing.T) {
		registry := NewRegistry()
		provider := &mockHealthCheckProvider{
			name: "test",
			result: &CheckResult{
				Status:  StatusUp,
				Details: map[string]any{"key": "value"},
			},
			scopes: []Scope{ScopeBase},
		}

		registry.Register(provider)
		result := registry.Check(nil)

		assert.Len(t, result, 1)
		assert.Contains(t, result, "test")
		assert.Equal(t, StatusUp, result["test"].Status)
	})

	t.Run("should register multiple providers", func(t *testing.T) {
		registry := NewRegistry()
		provider1 := &mockHealthCheckProvider{
			name:   "provider1",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeBase},
		}
		provider2 := &mockHealthCheckProvider{
			name:   "provider2",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeBase},
		}

		registry.Register(provider1)
		registry.Register(provider2)
		result := registry.Check(nil)

		assert.Len(t, result, 2)
		assert.Contains(t, result, "provider1")
		assert.Contains(t, result, "provider2")
	})
}

func TestCheck(t *testing.T) {
	t.Run("should check all providers when scope is nil", func(t *testing.T) {
		registry := NewRegistry()
		provider1 := &mockHealthCheckProvider{
			name:   "db",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeStartup, ScopeReady},
		}
		provider2 := &mockHealthCheckProvider{
			name:   "cache",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeLive},
		}

		registry.Register(provider1)
		registry.Register(provider2)
		result := registry.Check(nil)

		assert.Len(t, result, 2)
		assert.Contains(t, result, "db")
		assert.Contains(t, result, "cache")
	})

	t.Run("should only check providers with matching scope", func(t *testing.T) {
		registry := NewRegistry()
		provider1 := &mockHealthCheckProvider{
			name:   "db",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeStartup},
		}
		provider2 := &mockHealthCheckProvider{
			name:   "cache",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeLive},
		}

		registry.Register(provider1)
		registry.Register(provider2)

		scope := ScopeStartup
		result := registry.Check(&scope)

		assert.Len(t, result, 1)
		assert.Contains(t, result, "db")
		assert.NotContains(t, result, "cache")
	})

	t.Run("should check providers with multiple scopes", func(t *testing.T) {
		registry := NewRegistry()
		provider := &mockHealthCheckProvider{
			name:   "db",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeStartup, ScopeReady, ScopeLive},
		}

		registry.Register(provider)

		// Check startup scope
		startupScope := ScopeStartup
		startupResult := registry.Check(&startupScope)
		assert.Len(t, startupResult, 1)
		assert.Contains(t, startupResult, "db")

		// Check ready scope
		readyScope := ScopeReady
		readyResult := registry.Check(&readyScope)
		assert.Len(t, readyResult, 1)
		assert.Contains(t, readyResult, "db")

		// Check live scope
		liveScope := ScopeLive
		liveResult := registry.Check(&liveScope)
		assert.Len(t, liveResult, 1)
		assert.Contains(t, liveResult, "db")
	})

	t.Run("should omit providers that return errors", func(t *testing.T) {
		registry := NewRegistry()
		provider1 := &mockHealthCheckProvider{
			name:   "working",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeBase},
		}
		provider2 := &mockHealthCheckProvider{
			name:   "failing",
			err:    errors.New("test error"),
			scopes: []Scope{ScopeBase},
		}

		registry.Register(provider1)
		registry.Register(provider2)
		result := registry.Check(nil)

		assert.Len(t, result, 1)
		assert.Contains(t, result, "working")
		assert.NotContains(t, result, "failing")
	})

	t.Run("should return empty map when no providers registered", func(t *testing.T) {
		registry := NewRegistry()
		result := registry.Check(nil)

		assert.Empty(t, result)
	})

	t.Run("should check each provider only once for base scope", func(t *testing.T) {
		registry := NewRegistry()
		// Provider that appears in multiple scopes
		provider := &mockHealthCheckProvider{
			name:   "db",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeStartup, ScopeReady, ScopeLive},
		}

		registry.Register(provider)

		// When checking all (scope = nil), should appear only once
		result := registry.Check(nil)

		assert.Len(t, result, 1)
		assert.Contains(t, result, "db")
	})

	t.Run("should filter by base scope", func(t *testing.T) {
		registry := NewRegistry()
		baseProvider := &mockHealthCheckProvider{
			name:   "base-only",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeBase},
		}
		startupProvider := &mockHealthCheckProvider{
			name:   "startup-only",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeStartup},
		}

		registry.Register(baseProvider)
		registry.Register(startupProvider)

		baseScope := ScopeBase
		result := registry.Check(&baseScope)

		assert.Len(t, result, 1)
		assert.Contains(t, result, "base-only")
		assert.NotContains(t, result, "startup-only")
	})
}

func TestOverallStatus(t *testing.T) {
	t.Run("should return UP when all components are up", func(t *testing.T) {
		components := map[string]*CheckResult{
			"db":    {Status: StatusUp},
			"cache": {Status: StatusUp},
		}

		status := OverallStatus(components)

		assert.Equal(t, StatusUp, status)
	})

	t.Run("should return DOWN when any component is down", func(t *testing.T) {
		components := map[string]*CheckResult{
			"db":    {Status: StatusUp},
			"cache": {Status: StatusDown},
		}

		status := OverallStatus(components)

		assert.Equal(t, StatusDown, status)
	})

	t.Run("should return UP for empty components", func(t *testing.T) {
		components := map[string]*CheckResult{}

		status := OverallStatus(components)

		assert.Equal(t, StatusUp, status)
	})

	t.Run("should return DOWN when all components are down", func(t *testing.T) {
		components := map[string]*CheckResult{
			"db":    {Status: StatusDown},
			"cache": {Status: StatusDown},
		}

		status := OverallStatus(components)

		assert.Equal(t, StatusDown, status)
	})
}

func TestBuildResponse(t *testing.T) {
	t.Run("should build response with all healthy components", func(t *testing.T) {
		registry := NewRegistry()
		provider := &mockHealthCheckProvider{
			name: "database",
			result: &CheckResult{
				Status:  StatusUp,
				Details: map[string]any{"connection": "active"},
			},
			scopes: []Scope{ScopeReady},
		}
		registry.Register(provider)

		statusCode, response := registry.BuildResponse(nil)

		assert.Equal(t, 200, statusCode)
		assert.Equal(t, StatusUp, response.Status)
		assert.Len(t, response.Components, 1)
		assert.Contains(t, response.Components, "database")
		assert.Equal(t, StatusUp, response.Components["database"].Status)
		assert.Equal(t, map[string]any{"connection": "active"}, response.Components["database"].Details)
	})

	t.Run("should return 503 when any component is down", func(t *testing.T) {
		registry := NewRegistry()
		healthyProvider := &mockHealthCheckProvider{
			name:   "cache",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeReady},
		}
		unhealthyProvider := &mockHealthCheckProvider{
			name: "database",
			result: &CheckResult{
				Status:  StatusDown,
				Details: map[string]any{"error": "connection refused"},
			},
			scopes: []Scope{ScopeReady},
		}
		registry.Register(healthyProvider)
		registry.Register(unhealthyProvider)

		statusCode, response := registry.BuildResponse(nil)

		assert.Equal(t, 503, statusCode)
		assert.Equal(t, StatusDown, response.Status)
		assert.Len(t, response.Components, 2)
		assert.Equal(t, StatusUp, response.Components["cache"].Status)
		assert.Equal(t, StatusDown, response.Components["database"].Status)
	})

	t.Run("should filter by scope", func(t *testing.T) {
		registry := NewRegistry()
		startupProvider := &mockHealthCheckProvider{
			name:   "startup-check",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeStartup},
		}
		readyProvider := &mockHealthCheckProvider{
			name:   "ready-check",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeReady},
		}
		registry.Register(startupProvider)
		registry.Register(readyProvider)

		scope := ScopeStartup
		statusCode, response := registry.BuildResponse(&scope)

		assert.Equal(t, 200, statusCode)
		assert.Equal(t, StatusUp, response.Status)
		assert.Len(t, response.Components, 1)
		assert.Contains(t, response.Components, "startup-check")
		assert.NotContains(t, response.Components, "ready-check")
	})

	t.Run("should return UP with empty components when no providers registered", func(t *testing.T) {
		registry := NewRegistry()

		statusCode, response := registry.BuildResponse(nil)

		assert.Equal(t, 200, statusCode)
		assert.Equal(t, StatusUp, response.Status)
		assert.Empty(t, response.Components)
	})

	t.Run("should handle multiple scopes correctly", func(t *testing.T) {
		registry := NewRegistry()
		multiScopeProvider := &mockHealthCheckProvider{
			name:   "db",
			result: &CheckResult{Status: StatusUp},
			scopes: []Scope{ScopeStartup, ScopeReady, ScopeLive},
		}
		registry.Register(multiScopeProvider)

		// Check each scope individually
		startupScope := ScopeStartup
		startupCode, startupResp := registry.BuildResponse(&startupScope)
		assert.Equal(t, 200, startupCode)
		assert.Contains(t, startupResp.Components, "db")

		readyScope := ScopeReady
		readyCode, readyResp := registry.BuildResponse(&readyScope)
		assert.Equal(t, 200, readyCode)
		assert.Contains(t, readyResp.Components, "db")

		liveScope := ScopeLive
		liveCode, liveResp := registry.BuildResponse(&liveScope)
		assert.Equal(t, 200, liveCode)
		assert.Contains(t, liveResp.Components, "db")

		// Check all scopes (nil)
		allCode, allResp := registry.BuildResponse(nil)
		assert.Equal(t, 200, allCode)
		assert.Len(t, allResp.Components, 1) // Should appear only once
		assert.Contains(t, allResp.Components, "db")
	})
}
