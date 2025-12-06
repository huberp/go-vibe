package routes

import (
	"encoding/json"
	"myapp/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouterWithDB(t *testing.T) (*gin.Engine, *gorm.DB) {
	// Create test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Migrate schema
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create test logger
	logger, _ := zap.NewDevelopment()

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupRoutes(router, db, logger, "test-secret")

	return router, db
}

func TestInfoEndpoint(t *testing.T) {
	t.Run("should return info with build and user stats", func(t *testing.T) {
		router, db := setupTestRouterWithDB(t)

		// Add some test users
		users := []models.User{
			{Name: "User1", Email: "user1@test.com", Role: "user", PasswordHash: "hash1"},
			{Name: "User2", Email: "user2@test.com", Role: "user", PasswordHash: "hash2"},
			{Name: "Admin1", Email: "admin1@test.com", Role: "admin", PasswordHash: "hash3"},
		}
		for _, user := range users {
			db.Create(&user)
		}

		// Make request to /info endpoint
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info", nil)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check build info
		assert.Contains(t, response, "build")
		buildInfo := response["build"].(map[string]any)
		assert.Contains(t, buildInfo, "version")
		assert.Contains(t, buildInfo, "commit")
		assert.Contains(t, buildInfo, "build_time")
		assert.Contains(t, buildInfo, "go_version")

		// Check user stats
		assert.Contains(t, response, "users")
		userStats := response["users"].(map[string]any)
		assert.Equal(t, float64(3), userStats["total"])
		assert.Equal(t, float64(1), userStats["admins"])
		assert.Equal(t, float64(2), userStats["regular"])
	})

	t.Run("should return info with zero users when database is empty", func(t *testing.T) {
		router, _ := setupTestRouterWithDB(t)

		// Make request to /info endpoint
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info", nil)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check user stats are zero
		assert.Contains(t, response, "users")
		userStats := response["users"].(map[string]any)
		assert.Equal(t, float64(0), userStats["total"])
		assert.Equal(t, float64(0), userStats["admins"])
		assert.Equal(t, float64(0), userStats["regular"])
	})
}
