package test

import (
	"auth-api-jwt/controller"
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Register(ctx context.Context, request web.AuthRegisterRequest) (domain.User, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *AuthServiceMock) Login(ctx context.Context, request web.AuthLoginRequest) (string, error) {
	args := m.Called(ctx, request)
	return args.String(0), args.Error(1)
}

func TestAuthController_Register_Success(t *testing.T) {
	mockService := new(AuthServiceMock)

	requestBody := web.AuthRegisterRequest{
		Email:    "test@example.com",
		Password: "secret",
		FullName: "Test Name",
	}

	expectedUser := domain.User{
		Email:    requestBody.Email,
		FullName: requestBody.FullName,
	}

	mockService.On("Register", mock.Anything, mock.Anything).Return(expectedUser, nil)

	app := fiber.New()
	ctrl := controller.NewAuthController(mockService)
	app.Post("/register", ctrl.Register)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	raw, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(raw), requestBody.Email)
	mockService.AssertExpectations(t)
}

func TestUserController_Register_InvalidJSON(t *testing.T) {
	mockService := new(UserServiceMock)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Post("/register", ctrl.Create)

	body := []byte(`{ invalid json }`)

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthController_Register_ServiceError(t *testing.T) {
	mockService := new(AuthServiceMock)

	requestBody := web.AuthRegisterRequest{
		Email:    "err@example.com",
		Password: "secret",
		FullName: "Err",
	}

	mockService.On("Register", mock.Anything, mock.Anything).Return(domain.User{}, assert.AnError)

	app := fiber.New()
	ctrl := controller.NewAuthController(mockService)
	app.Post("/register", ctrl.Register)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthController_Login_Success(t *testing.T) {
	mockService := new(AuthServiceMock)

	requestBody := web.AuthLoginRequest{
		Email:    "login@example.com",
		Password: "secret",
	}

	token := "jwt.token.value"
	mockService.On("Login", mock.Anything, mock.Anything).Return(token, nil)

	app := fiber.New()
	ctrl := controller.NewAuthController(mockService)
	app.Post("/login", ctrl.Login)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUserController_Login_InvalidJSON(t *testing.T) {
	mockService := new(UserServiceMock)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Post("/login", ctrl.Create)

	body := []byte(`{ invalid json }`)

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthController_Login_ServiceError(t *testing.T) {
	mockService := new(AuthServiceMock)

	requestBody := web.AuthLoginRequest{
		Email:    "bad@example.com",
		Password: "secret",
	}

	mockService.On("Login", mock.Anything, mock.Anything).Return("", assert.AnError)

	app := fiber.New()
	ctrl := controller.NewAuthController(mockService)
	app.Post("/login", ctrl.Login)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}
