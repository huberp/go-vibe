package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
// @Description User account information
type User struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey" json:"id" example:"1"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name" example:"John Doe"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" example:"john@example.com"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:varchar(20);not null;default:'user'" json:"role" example:"user"`
	CreatedAt    time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
