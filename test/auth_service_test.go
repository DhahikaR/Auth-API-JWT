package test

import (
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
	"auth-api-jwt/service"
	"auth-api-jwt/utils"
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func (m *AuthRepositoryMock) Create(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *AuthRepositoryMock) FindByEmail(ctx context.Context, tx *gorm.DB, userEmail string) (domain.User, error) {
	args := m.Called(ctx, tx, userEmail)
	return args.Get(0).(domain.User), args.Error(1)
}

func TestAuthService_Register_Success(t *testing.T) {
	authMock := new(AuthRepositoryMock)
	userMock := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	request := web.AuthRegisterRequest{
		Email:    "reg@example.com",
		Password: "secret",
		FullName: "Reg User",
	}

	expected := domain.User{
		Id:           uuid.New(),
		Email:        request.Email,
		PasswordHash: "hashed",
		FullName:     request.FullName,
		Role:         "user",
		IsVerified:   false,
	}

	authMock.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("domain.User")).Return(expected, nil)

	svc := service.NewAuthService(authMock, userMock, db, validate)
	got, err := svc.Register(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, expected.Email, got.Email)
	authMock.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	authMock := new(AuthRepositoryMock)
	userMock := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	password := "mypassword"
	hashed, err := utils.HashPassword(password)
	assert.NoError(t, err)

	user := domain.User{
		Id:           uuid.New(),
		Email:        "login@example.com",
		PasswordHash: hashed,
		FullName:     "Login User",
		Role:         "user",
	}

	authMock.On("FindByEmail", mock.Anything, mock.Anything, user.Email).Return(user, nil)

	userMock.On("UpdateLastLogin", mock.Anything, mock.Anything, user.Id.String(), mock.Anything).Return(nil)

	svc := service.NewAuthService(authMock, userMock, db, validate)
	token, err := svc.Login(context.Background(), web.AuthLoginRequest{Email: user.Email, Password: password})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	authMock.AssertExpectations(t)
	userMock.AssertExpectations(t)
}

func TestAuthService_Register_InvalidRequest(t *testing.T) {
	authMock := new(AuthRepositoryMock)
	userMock := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	req := web.AuthRegisterRequest{
		Email: "not-an-email",
	}

	svc := service.NewAuthService(authMock, userMock, db, validate)
	_, err := svc.Register(context.Background(), req)
	if err == nil {
		t.Fatalf("expected validation error for invalid register request")
	}
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	authMock := new(AuthRepositoryMock)
	userMock := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	authMock.On("FindByEmail", mock.Anything, mock.Anything, "noone@example.com").Return(domain.User{}, assert.AnError)

	svc := service.NewAuthService(authMock, userMock, db, validate)
	token, err := svc.Login(context.Background(), web.AuthLoginRequest{Email: "noone@example.com", Password: "whatever"})
	if err == nil {
		t.Fatalf("expected error for unknown email login")
	}
	if token != "" {
		t.Fatalf("expected empty token on failed login")
	}
	authMock.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	authMock := new(AuthRepositoryMock)
	userMock := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	hashed, err := utils.HashPassword("correctpass")
	assert.NoError(t, err)

	user := domain.User{
		Id:           uuid.New(),
		Email:        "user@example.com",
		PasswordHash: hashed,
		FullName:     "User",
		Role:         "user",
	}

	authMock.On("FindByEmail", mock.Anything, mock.Anything, user.Email).Return(user, nil)

	svc := service.NewAuthService(authMock, userMock, db, validate)
	token, err := svc.Login(context.Background(), web.AuthLoginRequest{Email: user.Email, Password: "wrongpass"})
	if err == nil {
		t.Fatalf("expected error for wrong password")
	}
	if token != "" {
		t.Fatalf("expected empty token on failed login")
	}

	userMock.AssertNotCalled(t, "UpdateLastLogin", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	authMock.AssertExpectations(t)
}
