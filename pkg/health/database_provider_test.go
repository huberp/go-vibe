package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewDatabaseHealthCheckProvider(t *testing.T) {
	t.Run("should create provider with default scopes", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		provider := NewDatabaseHealthCheckProvider(db)

		assert.NotNil(t, provider)
		assert.Equal(t, "database", provider.Name())
		assert.Equal(t, []Scope{ScopeStartup, ScopeReady}, provider.Scopes())
	})

	t.Run("should create provider with custom scopes", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		provider := NewDatabaseHealthCheckProvider(db, ScopeBase, ScopeLive)

		assert.NotNil(t, provider)
		assert.Equal(t, []Scope{ScopeBase, ScopeLive}, provider.Scopes())
	})
}

func TestDatabaseHealthCheckProvider_Name(t *testing.T) {
	t.Run("should return database as name", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		provider := NewDatabaseHealthCheckProvider(db)

		assert.Equal(t, "database", provider.Name())
	})
}

func TestDatabaseHealthCheckProvider_Check(t *testing.T) {
	t.Run("should return UP status for healthy database", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		provider := NewDatabaseHealthCheckProvider(db)

		result, err := provider.Check()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, StatusUp, result.Status)
		assert.NotNil(t, result.Details)
		assert.Contains(t, result.Details, "max_open_connections")
	})

	t.Run("should return DOWN status for closed database", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		sqlDB, _ := db.DB()
		sqlDB.Close()

		provider := NewDatabaseHealthCheckProvider(db)
		result, err := provider.Check()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, StatusDown, result.Status)
		assert.Contains(t, result.Details["error"], "sql: database is closed")
	})

	t.Run("should return DOWN status for nil database", func(t *testing.T) {
		provider := NewDatabaseHealthCheckProvider(nil)
		result, err := provider.Check()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, StatusDown, result.Status)
		assert.Contains(t, result.Details["error"], "database not initialized")
	})
}

func TestDatabaseHealthCheckProvider_Scopes(t *testing.T) {
	t.Run("should return configured scopes", func(t *testing.T) {
		db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		provider := NewDatabaseHealthCheckProvider(db, ScopeStartup, ScopeReady, ScopeLive)

		scopes := provider.Scopes()

		assert.Len(t, scopes, 3)
		assert.Contains(t, scopes, ScopeStartup)
		assert.Contains(t, scopes, ScopeReady)
		assert.Contains(t, scopes, ScopeLive)
	})
}
