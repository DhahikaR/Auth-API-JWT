package service

import (
	"auth-api-jwt/helper"
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
	"auth-api-jwt/repository"
	"auth-api-jwt/utils"
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	AuthRepository repository.AuthRepository
	UserRepository repository.UserRepository
	DB             *gorm.DB
	Validate       *validator.Validate
}

func NewAuthService(authRepository repository.AuthRepository, userRepository repository.UserRepository, DB *gorm.DB, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		AuthRepository: authRepository,
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *AuthServiceImpl) Register(ctx context.Context, request web.AuthRegisterRequest) (domain.User, error) {
	if err := service.Validate.Struct(request); err != nil {
		return domain.User{}, err
	}

	hashed, _ := utils.HashPassword(request.Password)

	user := domain.User{
		Email:        request.Email,
		PasswordHash: hashed,
		FullName:     request.FullName,
		Role:         "user",
		IsVerified:   false,
	}

	return service.AuthRepository.Create(ctx, service.DB, user)
}

func (service *AuthServiceImpl) Login(ctx context.Context, request web.AuthLoginRequest) (string, error) {
	if err := service.Validate.Struct(request); err != nil {
		return "", err
	}

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := service.AuthRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !utils.CheckPassword(request.Password, user.PasswordHash) {
		return "", errors.New("invalid email or password")
	}

	err = service.UserRepository.UpdateLastLogin(ctx, tx, user.Id.String(), time.Now())
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateJWT(user.Id.String(), user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
