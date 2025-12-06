package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// Store user info in context
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("user_role", role)
		}

		c.Next()
	}
}

// RequireRole checks if user has the required role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found in token"})
			return
		}

		if userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}

// IsOwnerOrAdmin checks if the authenticated user is either the resource owner or an admin
// resourceUserID is the ID of the user resource being accessed
func IsOwnerOrAdmin(c *gin.Context, resourceUserID uint) bool {
	// Get the authenticated user's role
	userRole, roleExists := c.Get("user_role")
	if !roleExists {
		return false
	}

	// Admin users can access any resource
	if userRole == "admin" {
		return true
	}

	// Get the authenticated user's ID
	userID, idExists := c.Get("user_id")
	if !idExists {
		return false
	}

	// Regular users can only access their own resources
	return userID == resourceUserID
}
