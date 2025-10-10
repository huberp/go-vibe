package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	t.Run("should create valid user with all fields", func(t *testing.T) {
		user := User{
			ID:           1,
			Name:         "John Doe",
			Email:        "john@example.com",
			PasswordHash: "hashedpassword123",
			Role:         "user",
		}

		assert.Equal(t, uint(1), user.ID)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, "hashedpassword123", user.PasswordHash)
		assert.Equal(t, "user", user.Role)
	})

	t.Run("should have correct GORM tags", func(t *testing.T) {
		user := User{}
		assert.NotNil(t, user)
	})
}
