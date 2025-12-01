package service

import (
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
	"context"
)

type UserService interface {
	Create(ctx context.Context, request web.UserCreateRequest) (domain.User, error)
	Update(ctx context.Context, request web.UserUpdateRequest) (domain.User, error)
	UpdateMe(ctx context.Context, request web.UserUpdateRequest) (domain.User, error)
	Delete(ctx context.Context, targetUserId string) error
	FindById(ctx context.Context, targetUserId string) (domain.User, error)
	Me(ctx context.Context, targetUserId string) (domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
}
