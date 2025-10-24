package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should reject request without token", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuthMiddleware("secret"))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should reject request with invalid token", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuthMiddleware("secret"))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should allow user with correct role", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_role", "admin")
		})
		router.Use(RequireRole("admin"))
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should reject user with wrong role", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_role", "user")
		})
		router.Use(RequireRole("admin"))
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestIsOwnerOrAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should allow admin to access any resource", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_role", "admin")
		c.Set("user_id", uint(1))

		result := IsOwnerOrAdmin(c, uint(999))
		assert.True(t, result)
	})

	t.Run("should allow user to access their own resource", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_role", "user")
		c.Set("user_id", uint(5))

		result := IsOwnerOrAdmin(c, uint(5))
		assert.True(t, result)
	})

	t.Run("should deny user from accessing another user's resource", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_role", "user")
		c.Set("user_id", uint(5))

		result := IsOwnerOrAdmin(c, uint(10))
		assert.False(t, result)
	})

	t.Run("should deny access when user_role is missing", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", uint(5))

		result := IsOwnerOrAdmin(c, uint(5))
		assert.False(t, result)
	})

	t.Run("should deny access when user_id is missing", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_role", "user")

		result := IsOwnerOrAdmin(c, uint(5))
		assert.False(t, result)
	})
}
