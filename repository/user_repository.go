package repository

import (
	"auth-api-jwt/models/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	Save(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error)
	Update(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error)
	Delete(ctx context.Context, tx *gorm.DB, userId string) error
	FindById(ctx context.Context, tx *gorm.DB, userId string) (domain.User, error)
	FindAll(ctx context.Context, tx *gorm.DB) ([]domain.User, error)
	UpdateLastLogin(ctx context.Context, tx *gorm.DB, userId string, loginAt time.Time) error
}
