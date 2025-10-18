package info

import (
	"context"
	"myapp/internal/models"

	"gorm.io/gorm"
)

// UserStatsProvider provides user-related statistics.
type UserStatsProvider struct {
	db *gorm.DB
}

// NewUserStatsProvider creates a new UserStatsProvider.
func NewUserStatsProvider(db *gorm.DB) *UserStatsProvider {
	return &UserStatsProvider{
		db: db,
	}
}

// Name returns the name of this provider.
func (u *UserStatsProvider) Name() string {
	return "users"
}

// Info returns user statistics.
func (u *UserStatsProvider) Info() (map[string]interface{}, error) {
	ctx := context.Background()
	
	var totalUsers int64
	if err := u.db.WithContext(ctx).Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	var adminUsers int64
	if err := u.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", "admin").Count(&adminUsers).Error; err != nil {
		return nil, err
	}

	var regularUsers int64
	if err := u.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", "user").Count(&regularUsers).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":   totalUsers,
		"admins":  adminUsers,
		"regular": regularUsers,
	}, nil
}
