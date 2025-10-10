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
