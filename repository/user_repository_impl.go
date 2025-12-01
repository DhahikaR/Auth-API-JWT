package repository

import (
	"auth-api-jwt/models/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	err := tx.WithContext(ctx).Create(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	err := tx.WithContext(ctx).Model(domain.User{}).Where("id = ?", user.Id).Updates(map[string]interface{}{
		"email":         user.Email,
		"password_hash": user.PasswordHash,
		"full_name":     user.FullName,
		"role":          user.Role,
	}).Error

	return user, err
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, userId string) error {
	return tx.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userId).Update("deleted_at", gorm.DeletedAt{Valid: true}).Error
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *gorm.DB, userId string) (domain.User, error) {
	var user domain.User
	result := tx.Where("id = ?", userId).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, result.Error
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *gorm.DB) ([]domain.User, error) {
	var users []domain.User
	err := tx.WithContext(ctx).Find(&users).Error

	return users, err
}

func (repository *UserRepositoryImpl) UpdateLastLogin(ctx context.Context, tx *gorm.DB, userId string, loginAt time.Time) error {
	return tx.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userId).Update("last_login_at", loginAt).Error
}
