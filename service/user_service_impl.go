package service

import (
	"auth-api-jwt/helper"
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
	"auth-api-jwt/repository"
	"auth-api-jwt/utils"

	"context"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *gorm.DB
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, DB *gorm.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *UserServiceImpl) Create(ctx context.Context, request web.UserCreateRequest) (domain.User, error) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Email:        request.Email,
		PasswordHash: request.Password,
		FullName:     request.FullName,
		Role:         request.Role,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	created, err := service.UserRepository.Save(ctx, tx, user)
	if err != nil {
		return domain.User{}, err
	}

	return created, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, request web.UserUpdateRequest) (domain.User, error) {
	if err := service.Validate.Struct(request); err != nil {
		return domain.User{}, err
	}

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, request.Id.String())
	if err != nil {
		return domain.User{}, err
	}

	user.FullName = request.FullName
	user.Email = request.Email
	user.Role = request.Role

	if request.PasswordHash != "" {
		hashed, _ := utils.HashPassword(request.PasswordHash)
		user.PasswordHash = hashed
	}

	updated, err := service.UserRepository.Update(ctx, tx, user)
	if err != nil {
		return domain.User{}, err
	}

	return updated, nil
}

func (service *UserServiceImpl) UpdateMe(ctx context.Context, request web.UserUpdateRequest) (domain.User, error) {
	if err := service.Validate.Struct(request); err != nil {
		return domain.User{}, err
	}

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, request.Id.String())
	if err != nil {
		return domain.User{}, err
	}

	user.FullName = request.FullName
	user.Email = request.Email

	if request.PasswordHash != "" {
		hashed, _ := utils.HashPassword(request.PasswordHash)
		user.PasswordHash = hashed
	}

	updated, err := service.UserRepository.Update(ctx, tx, user)
	if err != nil {
		return domain.User{}, err
	}

	return updated, nil
}

func (service *UserServiceImpl) Delete(ctx context.Context, targetUserId string) error {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := service.UserRepository.FindById(ctx, tx, targetUserId)
	if err != nil {
		return err
	}

	if err := service.UserRepository.Delete(ctx, tx, targetUserId); err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) FindById(ctx context.Context, targetUserId string) (domain.User, error) {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, targetUserId)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (service *UserServiceImpl) Me(ctx context.Context, targetUserId string) (domain.User, error) {
	return service.FindById(ctx, targetUserId)
}

func (service *UserServiceImpl) FindAll(ctx context.Context) ([]domain.User, error) {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		return []domain.User{}, err
	}

	return user, nil
}
