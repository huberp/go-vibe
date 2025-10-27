package examples

import (
	"myapp/internal/models"
	"myapp/pkg/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUserRepositoryWithFixtures demonstrates using fixtures for integration testing
func TestUserRepositoryWithFixtures(t *testing.T) {
	t.Run("should query users loaded from fixtures", func(t *testing.T) {
		// Setup test database with fixtures
		tdb := testutil.NewTestDB(t)
		defer tdb.Cleanup()
		tdb.LoadFixtures(t, "../../../testdata/fixtures/minimal")

		// Query all users
		var users []models.User
		err := tdb.DB.Find(&users).Error
		assert.NoError(t, err)
		assert.Len(t, users, 2) // minimal fixture has 2 users

		// Verify user details
		var testUser models.User
		err = tdb.DB.Where("email = ?", "test@example.com").First(&testUser).Error
		assert.NoError(t, err)
		assert.Equal(t, "Test User", testUser.Name)
		assert.Equal(t, "user", testUser.Role)

		// Verify admin user
		var adminUser models.User
		err = tdb.DB.Where("email = ?", "admin@example.com").First(&adminUser).Error
		assert.NoError(t, err)
		assert.Equal(t, "Test Admin", adminUser.Name)
		assert.Equal(t, "admin", adminUser.Role)
	})

	t.Run("should support data modification and reload", func(t *testing.T) {
		// Setup test database with fixtures
		tdb := testutil.NewTestDB(t)
		defer tdb.Cleanup()
		tdb.LoadFixtures(t, "../../../testdata/fixtures/minimal")

		// Modify a user
		result := tdb.DB.Model(&models.User{}).
			Where("email = ?", "test@example.com").
			Update("name", "Modified Name")
		assert.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)

		// Verify modification
		var user models.User
		tdb.DB.Where("email = ?", "test@example.com").First(&user)
		assert.Equal(t, "Modified Name", user.Name)

		// Reload fixtures to restore original data
		tdb.ReloadFixtures(t)

		// Verify original data is restored
		var restoredUser models.User
		tdb.DB.Where("email = ?", "test@example.com").First(&restoredUser)
		assert.Equal(t, "Test User", restoredUser.Name)
	})

	t.Run("should support full fixture set", func(t *testing.T) {
		// Setup test database with full fixtures
		tdb := testutil.NewTestDB(t)
		defer tdb.Cleanup()
		tdb.LoadFixtures(t, "../../../testdata/fixtures/full")

		// Count total users
		var count int64
		tdb.DB.Model(&models.User{}).Count(&count)
		assert.Equal(t, int64(5), count)

		// Count admin users
		var adminCount int64
		tdb.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
		assert.Equal(t, int64(2), adminCount)

		// Count regular users
		var userCount int64
		tdb.DB.Model(&models.User{}).Where("role = ?", "user").Count(&userCount)
		assert.Equal(t, int64(3), userCount)
	})

	t.Run("should support transaction-based testing", func(t *testing.T) {
		// Setup test database
		tdb := testutil.NewTestDB(t)
		defer tdb.Cleanup()
		tdb.LoadFixtures(t, "../../../testdata/fixtures/minimal")

		// Start a transaction
		tx := tdb.BeginTransaction()
		defer tx.Rollback()

		// Create a new user in the transaction
		newUser := &models.User{
			Name:         "Transaction User",
			Email:        "transaction@example.com",
			PasswordHash: "hash",
			Role:         "user",
		}
		err := tx.Create(newUser).Error
		assert.NoError(t, err)

		// Verify user exists in transaction
		var foundInTx models.User
		err = tx.Where("email = ?", "transaction@example.com").First(&foundInTx).Error
		assert.NoError(t, err)
		assert.Equal(t, "Transaction User", foundInTx.Name)

		// Rollback transaction
		tx.Rollback()

		// Verify user doesn't exist in main DB
		var foundInDB models.User
		err = tdb.DB.Where("email = ?", "transaction@example.com").First(&foundInDB).Error
		assert.Error(t, err) // Should be record not found
	})
}
