package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		password := "testpassword123"
		hash, err := HashPassword(password)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("should verify correct password", func(t *testing.T) {
		password := "testpassword123"
		hash, _ := HashPassword(password)
		
		result := CheckPasswordHash(password, hash)
		assert.True(t, result)
	})

	t.Run("should reject incorrect password", func(t *testing.T) {
		password := "testpassword123"
		hash, _ := HashPassword(password)
		
		result := CheckPasswordHash("wrongpassword", hash)
		assert.False(t, result)
	})
}

func TestGenerateJWT(t *testing.T) {
	t.Run("should generate valid JWT token", func(t *testing.T) {
		userID := uint(123)
		role := "admin"
		secret := "test-secret"
		
		token, err := GenerateJWT(userID, role, secret)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
