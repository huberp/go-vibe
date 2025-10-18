package info

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuildInfoProvider(t *testing.T) {
	t.Run("should create provider with specified values", func(t *testing.T) {
		provider := NewBuildInfoProvider("1.0.0", "abc123", "2024-01-01T00:00:00Z", "go1.25.2")
		
		assert.NotNil(t, provider)
		assert.Equal(t, "1.0.0", provider.Version)
		assert.Equal(t, "abc123", provider.CommitSHA)
		assert.Equal(t, "2024-01-01T00:00:00Z", provider.BuildTime)
		assert.Equal(t, "go1.25.2", provider.GoVersion)
	})

	t.Run("should use defaults for empty values", func(t *testing.T) {
		provider := NewBuildInfoProvider("", "", "", "")
		
		assert.NotNil(t, provider)
		assert.Equal(t, "dev", provider.Version)
		assert.Equal(t, "unknown", provider.CommitSHA)
		assert.NotEmpty(t, provider.BuildTime)
		assert.Equal(t, "unknown", provider.GoVersion)
	})
}

func TestBuildInfoProvider_Name(t *testing.T) {
	t.Run("should return build as name", func(t *testing.T) {
		provider := NewBuildInfoProvider("1.0.0", "abc123", "2024-01-01T00:00:00Z", "go1.25.2")
		assert.Equal(t, "build", provider.Name())
	})
}

func TestBuildInfoProvider_Info(t *testing.T) {
	t.Run("should return build information", func(t *testing.T) {
		provider := NewBuildInfoProvider("1.0.0", "abc123", "2024-01-01T00:00:00Z", "go1.25.2")
		
		info, err := provider.Info()
		
		assert.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, "1.0.0", info["version"])
		assert.Equal(t, "abc123", info["commit"])
		assert.Equal(t, "2024-01-01T00:00:00Z", info["build_time"])
		assert.Equal(t, "go1.25.2", info["go_version"])
	})

	t.Run("should return default values in info", func(t *testing.T) {
		provider := NewBuildInfoProvider("", "", "", "")
		
		info, err := provider.Info()
		
		assert.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, "dev", info["version"])
		assert.Equal(t, "unknown", info["commit"])
		assert.NotEmpty(t, info["build_time"])
		assert.Equal(t, "unknown", info["go_version"])
	})
}
