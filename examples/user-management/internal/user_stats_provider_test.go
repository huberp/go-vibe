package info

import (
	"myapp/examples/user-management/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestNewUserStatsProvider(t *testing.T) {
	t.Run("should create provider", func(t *testing.T) {
		db := setupTestDB(t)
		provider := NewUserStatsProvider(db)

		assert.NotNil(t, provider)
		assert.NotNil(t, provider.db)
	})
}

func TestUserStatsProvider_Name(t *testing.T) {
	t.Run("should return users as name", func(t *testing.T) {
		db := setupTestDB(t)
		provider := NewUserStatsProvider(db)
		assert.Equal(t, "users", provider.Name())
	})
}

func TestUserStatsProvider_Info(t *testing.T) {
	t.Run("should return zero counts for empty database", func(t *testing.T) {
		db := setupTestDB(t)
		provider := NewUserStatsProvider(db)

		info, err := provider.Info()

		assert.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, int64(0), info["total"])
		assert.Equal(t, int64(0), info["admins"])
		assert.Equal(t, int64(0), info["regular"])
	})

	t.Run("should return correct user counts", func(t *testing.T) {
		db := setupTestDB(t)

		// Create test users
		users := []models.User{
			{Name: "User1", Email: "user1@example.com", Role: "user", PasswordHash: "hash1"},
			{Name: "User2", Email: "user2@example.com", Role: "user", PasswordHash: "hash2"},
			{Name: "Admin1", Email: "admin1@example.com", Role: "admin", PasswordHash: "hash3"},
			{Name: "Admin2", Email: "admin2@example.com", Role: "admin", PasswordHash: "hash4"},
			{Name: "User3", Email: "user3@example.com", Role: "user", PasswordHash: "hash5"},
		}

		for _, user := range users {
			if err := db.Create(&user).Error; err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}
		}

		provider := NewUserStatsProvider(db)
		info, err := provider.Info()

		assert.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, int64(5), info["total"])
		assert.Equal(t, int64(2), info["admins"])
		assert.Equal(t, int64(3), info["regular"])
	})

	t.Run("should handle only admin users", func(t *testing.T) {
		db := setupTestDB(t)

		// Create only admin users
		users := []models.User{
			{Name: "Admin1", Email: "admin1@example.com", Role: "admin", PasswordHash: "hash1"},
			{Name: "Admin2", Email: "admin2@example.com", Role: "admin", PasswordHash: "hash2"},
		}

		for _, user := range users {
			if err := db.Create(&user).Error; err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}
		}

		provider := NewUserStatsProvider(db)
		info, err := provider.Info()

		assert.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, int64(2), info["total"])
		assert.Equal(t, int64(2), info["admins"])
		assert.Equal(t, int64(0), info["regular"])
	})
}
