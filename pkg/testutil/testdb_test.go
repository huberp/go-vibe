package testutil

import (
	"myapp/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestDB(t *testing.T) {
	t.Run("should create test database successfully", func(t *testing.T) {
		tdb := NewTestDB(t)
		defer tdb.Cleanup()

		assert.NotNil(t, tdb.DB)
		assert.NotNil(t, tdb.sqlDB)
	})
}

func TestLoadFixtures(t *testing.T) {
	t.Run("should load fixtures from directory", func(t *testing.T) {
		tdb := NewTestDB(t)
		defer tdb.Cleanup()

		// Load minimal fixtures
		tdb.LoadFixtures(t, "../../testdata/fixtures/minimal")

		// Verify fixtures were loaded
		var count int64
		tdb.DB.Model(&models.User{}).Count(&count)
		assert.Greater(t, count, int64(0))

		// Verify specific user
		var user models.User
		err := tdb.DB.Where("email = ?", "test@example.com").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "user", user.Role)
	})
}

func TestReloadFixtures(t *testing.T) {
	t.Run("should reload fixtures after modification", func(t *testing.T) {
		tdb := NewTestDB(t)
		defer tdb.Cleanup()

		// Load fixtures
		tdb.LoadFixtures(t, "../../testdata/fixtures/minimal")

		// Modify a user
		tdb.DB.Model(&models.User{}).Where("email = ?", "test@example.com").Update("name", "Modified Name")

		// Verify modification
		var user models.User
		tdb.DB.Where("email = ?", "test@example.com").First(&user)
		assert.Equal(t, "Modified Name", user.Name)

		// Reload fixtures
		tdb.ReloadFixtures(t)

		// Verify original data is restored
		tdb.DB.Where("email = ?", "test@example.com").First(&user)
		assert.Equal(t, "Test User", user.Name)
	})
}

func TestTruncateTables(t *testing.T) {
	t.Run("should truncate all tables", func(t *testing.T) {
		tdb := NewTestDB(t)
		defer tdb.Cleanup()

		// Load fixtures
		tdb.LoadFixtures(t, "../../testdata/fixtures/minimal")

		// Verify data exists
		var countBefore int64
		tdb.DB.Model(&models.User{}).Count(&countBefore)
		assert.Greater(t, countBefore, int64(0))

		// Truncate
		tdb.TruncateTables(t)

		// Verify data is removed
		var countAfter int64
		tdb.DB.Model(&models.User{}).Count(&countAfter)
		assert.Equal(t, int64(0), countAfter)
	})
}

func TestBeginTransaction(t *testing.T) {
	t.Run("should create transaction for test isolation", func(t *testing.T) {
		tdb := NewTestDB(t)
		defer tdb.Cleanup()

		// Start transaction
		tx := tdb.BeginTransaction()
		defer tx.Rollback()

		// Create user in transaction
		user := &models.User{
			Name:         "Transaction Test",
			Email:        "tx@example.com",
			PasswordHash: "hash",
			Role:         "user",
		}
		err := tx.Create(user).Error
		assert.NoError(t, err)

		// Verify user exists in transaction
		var foundInTx models.User
		err = tx.Where("email = ?", "tx@example.com").First(&foundInTx).Error
		assert.NoError(t, err)

		// Rollback transaction
		tx.Rollback()

		// Verify user does not exist in main DB
		var foundInDB models.User
		err = tdb.DB.Where("email = ?", "tx@example.com").First(&foundInDB).Error
		assert.Error(t, err) // Should not find the user
	})
}
