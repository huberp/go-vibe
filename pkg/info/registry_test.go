package info

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockProvider is a simple mock implementation of InfoProvider for testing
type mockProvider struct {
	name string
	data map[string]interface{}
	err  error
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Info() (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

func TestNewRegistry(t *testing.T) {
	t.Run("should create empty registry", func(t *testing.T) {
		registry := NewRegistry()
		assert.NotNil(t, registry)
		result := registry.GetAll()
		assert.Empty(t, result)
	})
}

func TestRegister(t *testing.T) {
	t.Run("should register provider", func(t *testing.T) {
		registry := NewRegistry()
		provider := &mockProvider{
			name: "test",
			data: map[string]interface{}{"key": "value"},
		}

		registry.Register(provider)
		result := registry.GetAll()

		assert.Len(t, result, 1)
		assert.Contains(t, result, "test")
		assert.Equal(t, map[string]interface{}{"key": "value"}, result["test"])
	})

	t.Run("should register multiple providers", func(t *testing.T) {
		registry := NewRegistry()
		provider1 := &mockProvider{
			name: "provider1",
			data: map[string]interface{}{"key1": "value1"},
		}
		provider2 := &mockProvider{
			name: "provider2",
			data: map[string]interface{}{"key2": "value2"},
		}

		registry.Register(provider1)
		registry.Register(provider2)
		result := registry.GetAll()

		assert.Len(t, result, 2)
		assert.Contains(t, result, "provider1")
		assert.Contains(t, result, "provider2")
	})
}

func TestGetAll(t *testing.T) {
	t.Run("should aggregate all providers", func(t *testing.T) {
		registry := NewRegistry()
		provider1 := &mockProvider{
			name: "build",
			data: map[string]interface{}{"version": "1.0.0"},
		}
		provider2 := &mockProvider{
			name: "stats",
			data: map[string]interface{}{"count": 42},
		}

		registry.Register(provider1)
		registry.Register(provider2)
		result := registry.GetAll()

		assert.Equal(t, map[string]interface{}{
			"build": map[string]interface{}{"version": "1.0.0"},
			"stats": map[string]interface{}{"count": 42},
		}, result)
	})

	t.Run("should omit providers that return errors", func(t *testing.T) {
		registry := NewRegistry()
		provider1 := &mockProvider{
			name: "working",
			data: map[string]interface{}{"key": "value"},
		}
		provider2 := &mockProvider{
			name: "failing",
			err:  errors.New("test error"),
		}

		registry.Register(provider1)
		registry.Register(provider2)
		result := registry.GetAll()

		assert.Len(t, result, 1)
		assert.Contains(t, result, "working")
		assert.NotContains(t, result, "failing")
	})

	t.Run("should return empty map when no providers registered", func(t *testing.T) {
		registry := NewRegistry()
		result := registry.GetAll()

		assert.Empty(t, result)
	})
}
