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
				Details: map[string]interface{}{"key": "value"},
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
