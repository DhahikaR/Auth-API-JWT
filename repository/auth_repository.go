package repository

import (
	"auth-api-jwt/models/domain"
	"context"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, tx *gorm.DB, userEmail string) (domain.User, error)
}
