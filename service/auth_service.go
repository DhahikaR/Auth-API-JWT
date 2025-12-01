package service

import (
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, request web.AuthRegisterRequest) (domain.User, error)
	Login(ctx context.Context, request web.AuthLoginRequest) (string, error)
}
