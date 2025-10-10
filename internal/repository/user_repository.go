package repository

import (
	"context"
	"myapp/internal/models"
)

//go:generate mockgen -source=user_repository.go -destination=user_repository_mock.go -package=repository

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindAll(ctx context.Context) ([]models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}
