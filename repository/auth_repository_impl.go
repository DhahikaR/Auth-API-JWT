package repository

import (
	"auth-api-jwt/models/domain"
	"context"

	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		DB: db,
	}
}

func (repository *AuthRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	err := tx.WithContext(ctx).Create(&user).Error
	return user, err
}

func (repository *AuthRepositoryImpl) FindByEmail(ctx context.Context, tx *gorm.DB, userEmail string) (domain.User, error) {
	var user domain.User
	result := tx.Where("email = ?", userEmail).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, result.Error
}
