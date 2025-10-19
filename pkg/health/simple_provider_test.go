package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSimpleHealthCheckProvider(t *testing.T) {
	t.Run("should create provider with default scopes", func(t *testing.T) {
		provider := NewSimpleHealthCheckProvider("app")

		assert.NotNil(t, provider)
		assert.Equal(t, "app", provider.Name())
		assert.Equal(t, []Scope{ScopeLive}, provider.Scopes())
	})

	t.Run("should create provider with custom scopes", func(t *testing.T) {
		provider := NewSimpleHealthCheckProvider("app", ScopeBase, ScopeReady)

		assert.NotNil(t, provider)
		assert.Equal(t, []Scope{ScopeBase, ScopeReady}, provider.Scopes())
	})
}

func TestSimpleHealthCheckProvider_Name(t *testing.T) {
	t.Run("should return configured name", func(t *testing.T) {
		provider := NewSimpleHealthCheckProvider("my-service")

		assert.Equal(t, "my-service", provider.Name())
	})
}

func TestSimpleHealthCheckProvider_Check(t *testing.T) {
	t.Run("should always return UP status", func(t *testing.T) {
		provider := NewSimpleHealthCheckProvider("app")

		result, err := provider.Check()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, StatusUp, result.Status)
		assert.NotNil(t, result.Details)
		assert.Contains(t, result.Details, "timestamp")
	})
}

func TestSimpleHealthCheckProvider_Scopes(t *testing.T) {
	t.Run("should return configured scopes", func(t *testing.T) {
		provider := NewSimpleHealthCheckProvider("app", ScopeStartup, ScopeReady, ScopeLive)

		scopes := provider.Scopes()

		assert.Len(t, scopes, 3)
		assert.Contains(t, scopes, ScopeStartup)
		assert.Contains(t, scopes, ScopeReady)
		assert.Contains(t, scopes, ScopeLive)
	})
}
