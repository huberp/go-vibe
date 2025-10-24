package middleware

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		assert.Contains(t, w.Body.String(), "authorization header required")
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
		assert.Contains(t, w.Body.String(), "invalid token")
	})

	t.Run("should reject request with malformed authorization header", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuthMiddleware("secret"))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		testCases := []struct {
			name   string
			header string
		}{
			{"missing Bearer prefix", "token.without.bearer"},
			{"only Bearer", "Bearer"},
			{"Bearer with extra parts", "Bearer token extra"},
			{"wrong prefix", "Basic sometoken"},
			{"empty Bearer", "Bearer "},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/protected", nil)
				req.Header.Set("Authorization", tc.header)
				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusUnauthorized, w.Code)
			})
		}
	})

	t.Run("should accept valid token with claims", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		router.Use(JWTAuthMiddleware(secret))
		router.GET("/protected", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			assert.True(t, exists)
			assert.Equal(t, uint(123), userID)

			userRole, exists := c.Get("user_role")
			assert.True(t, exists)
			assert.Equal(t, "admin", userRole)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create a valid token
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"role":    "admin",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should reject token with wrong signing method", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		router.Use(JWTAuthMiddleware(secret))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create a token with RS256 instead of HS256
		// This will fail validation because we expect HMAC
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"role":    "admin",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		// Manually craft a token that looks like it uses a different signing method
		// For this test, we'll just use an invalid token which will also fail
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9."+tokenString[len("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."):])
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should reject expired token", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		router.Use(JWTAuthMiddleware(secret))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create an expired token
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"role":    "admin",
			"exp":     time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should reject token with wrong secret", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuthMiddleware("correct-secret"))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create a token with wrong secret
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"role":    "admin",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("wrong-secret"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should handle token with missing user_id claim", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		router.Use(JWTAuthMiddleware(secret))
		router.GET("/protected", func(c *gin.Context) {
			// Should still succeed, but user_id won't be in context
			_, exists := c.Get("user_id")
			assert.False(t, exists)

			userRole, exists := c.Get("user_role")
			assert.True(t, exists)
			assert.Equal(t, "admin", userRole)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create a token without user_id
		claims := jwt.MapClaims{
			"role": "admin",
			"exp":  time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should handle token with missing role claim", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		router.Use(JWTAuthMiddleware(secret))
		router.GET("/protected", func(c *gin.Context) {
			// Should still succeed, but role won't be in context
			userID, exists := c.Get("user_id")
			assert.True(t, exists)
			assert.Equal(t, uint(123), userID)

			_, exists = c.Get("user_role")
			assert.False(t, exists)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create a token without role
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should handle token with invalid claim types", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		router.Use(JWTAuthMiddleware(secret))
		router.GET("/protected", func(c *gin.Context) {
			// Claims with wrong types should be ignored
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create a token with wrong claim types
		claims := jwt.MapClaims{
			"user_id": "not-a-number", // Should be float64
			"role":    123,            // Should be string
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
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
