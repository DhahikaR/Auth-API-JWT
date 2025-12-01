package test

import (
	"auth-api-jwt/controller"
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) Create(ctx context.Context, request web.UserCreateRequest) (domain.User, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserServiceMock) Update(ctx context.Context, request web.UserUpdateRequest) (domain.User, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserServiceMock) UpdateMe(ctx context.Context, request web.UserUpdateRequest) (domain.User, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserServiceMock) Delete(ctx context.Context, targetUserId string) error {
	args := m.Called(ctx, targetUserId)
	return args.Error(0)
}

func (m *UserServiceMock) FindById(ctx context.Context, targetUserId string) (domain.User, error) {
	args := m.Called(ctx, targetUserId)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserServiceMock) Me(ctx context.Context, targetUserId string) (domain.User, error) {
	args := m.Called(ctx, targetUserId)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserServiceMock) FindAll(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func TestUserController_Create_Success(t *testing.T) {
	mockService := new(UserServiceMock)

	requestBody := web.UserCreateRequest{
		Email:    "ucreate@example.com",
		Password: "pass",
		FullName: "Create User",
		Role:     "user",
	}

	expected := domain.User{
		Email:    requestBody.Email,
		FullName: requestBody.FullName,
		Role:     requestBody.Role,
	}
	mockService.On("Create", mock.Anything, mock.Anything).Return(expected, nil)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Post("/users", ctrl.Create)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	raw, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(raw), requestBody.Email)

	mockService.AssertExpectations(t)
}

func TestUserController_Create_InvalidJSON(t *testing.T) {
	mockService := new(UserServiceMock)

	requestBody := []byte(`{ invalid json }`)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Post("/users", ctrl.Create)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUserController_Create_ValidationError(t *testing.T) {
	mockService := new(UserServiceMock)

	requestBody := web.UserCreateRequest{
		Email:    "",
		Password: "secret",
		FullName: "Test Name",
	}

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Post("/users", func(c *fiber.Ctx) error {
		return ctrl.Create(c)
	})

	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertNotCalled(t, "Create")
}

func TestUserController_Create_ServiceError(t *testing.T) {
	mockService := new(UserServiceMock)

	requestBody := web.UserCreateRequest{
		Email:    "err@example.com",
		Password: "password",
		FullName: "Test Name",
	}
	mockService.On("Create", mock.Anything, mock.Anything).Return(domain.User{}, assert.AnError)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Post("/users", ctrl.Create)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUserController_Update_Success(t *testing.T) {
	mockService := new(UserServiceMock)

	targetId := uuid.New()

	requestBody := web.UserUpdateRequest{
		Email:    "updated@example.com",
		FullName: "Updated Name",
	}

	requestBody.Id = targetId

	updated := domain.User{
		Id:       targetId,
		Email:    requestBody.Email,
		FullName: requestBody.FullName,
	}

	mockService.On("Update", mock.Anything, mock.AnythingOfType("web.UserUpdateRequest")).Return(updated, nil)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Put("/users/:userId", func(c *fiber.Ctx) error {
		c.Locals("userId", "admin-id")
		c.Locals("role", "admin")
		return ctrl.Update(c)
	})

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/users/"+targetId.String(), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	raw, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(raw), requestBody.FullName)

	mockService.AssertExpectations(t)
}

func TestUserController_Update_InvalidUUID(t *testing.T) {
	mockService := new(UserServiceMock)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)

	app.Put("/users/:userId", func(c *fiber.Ctx) error {
		c.Locals("role", "admin")
		return ctrl.Update(c)
	})

	req := httptest.NewRequest("PUT", "/users/not-a-uuid", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertNotCalled(t, "Update")
}

func TestUserController_Update_InvalidJSON(t *testing.T) {
	mockService := new(UserServiceMock)
	targetId := uuid.New().String()

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)

	app.Put("/users/:userId", func(c *fiber.Ctx) error {
		c.Locals("userId", targetId)
		c.Locals("role", "user")
		return ctrl.Update(c)
	})

	requestBody := []byte(`{ invalid json }`)
	req := httptest.NewRequest("PUT", "/users/"+targetId, bytes.NewReader(requestBody))

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUserController_Update_ServiceError(t *testing.T) {
	mockService := new(UserServiceMock)

	targetId := uuid.New().String()
	requestBody := web.UserUpdateRequest{
		Email:    "test@example.com",
		FullName: "Updated Name",
	}

	mockService.
		On("Update", mock.Anything, mock.AnythingOfType("web.UserUpdateRequest")).
		Return(domain.User{}, errors.New("database error"))

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)

	app.Put("/users/:userId", func(c *fiber.Ctx) error {
		c.Locals("userId", targetId)
		c.Locals("role", "admin")
		return ctrl.Update(c)
	})

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/users/"+targetId, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUserController_UpdateMe_Success(t *testing.T) {
	mockService := new(UserServiceMock)

	userId := uuid.New().String()

	requestBody := web.UserUpdateRequest{
		Email:    "meupdated@mail.com",
		FullName: "Updated Me",
	}

	expectedUser := domain.User{
		Id:       uuid.MustParse(userId),
		Email:    requestBody.Email,
		FullName: requestBody.FullName,
		Role:     "user",
	}

	mockService.On("UpdateMe", mock.Anything, mock.AnythingOfType("web.UserUpdateRequest")).Return(expectedUser, nil)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Put("/users/me", func(c *fiber.Ctx) error {
		c.Locals("userId", userId)
		c.Locals("role", "user")
		return ctrl.UpdateMe(c)
	})

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/users/me", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	raw, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(raw), requestBody.Email)

	mockService.AssertExpectations(t)
}

func TestUserController_UpdateMe_InvalidJSON(t *testing.T) {
	mockService := new(UserServiceMock)
	targetId := uuid.New().String()

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)

	app.Put("/users/me", func(c *fiber.Ctx) error {
		c.Locals("userId", targetId)
		c.Locals("role", "user")
		return ctrl.UpdateMe(c)
	})

	body := []byte(`{ invalid json }`)
	req := httptest.NewRequest("PUT", "/users/me", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertNotCalled(t, "UpdateMe")
}

func TestUserController_UpdateMe_ServiceError(t *testing.T) {
	mockService := new(UserServiceMock)

	targetId := uuid.New().String()
	requestBody := web.UserUpdateRequest{
		Email:    "test@example.com",
		FullName: "Updated Name",
	}

	mockService.
		On("UpdateMe", mock.Anything, mock.AnythingOfType("web.UserUpdateRequest")).
		Return(domain.User{}, errors.New("database error"))

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)

	app.Put("/users/me", func(c *fiber.Ctx) error {
		c.Locals("userId", targetId)
		c.Locals("role", "user")
		return ctrl.UpdateMe(c)
	})

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/users/me", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUserController_FindById_Success(t *testing.T) {
	mockService := new(UserServiceMock)
	targetId := uuid.New().String()

	expectedUser := domain.User{
		Id:       uuid.MustParse(targetId),
		Email:    "test@example.com",
		FullName: "Test User",
		Role:     "user",
	}

	mockService.On("FindById", mock.Anything, targetId).Return(expectedUser, nil)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Get("/users/:userId", func(c *fiber.Ctx) error {
		c.Locals("userId", targetId)
		c.Locals("role", "user")
		return ctrl.FindById(c)
	})

	req := httptest.NewRequest("GET", "/users/"+targetId, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]any
	json.NewDecoder(resp.Body).Decode(&response)
	data := response["Data"].(map[string]any)
	assert.Equal(t, expectedUser.Email, data["email"])
	assert.Equal(t, expectedUser.FullName, data["full_name"])

	mockService.AssertExpectations(t)
}

func TestUserController_FindById_InvalidUUID(t *testing.T) {
	mockService := new(UserServiceMock)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Get("/users/:userId", func(c *fiber.Ctx) error {
		c.Locals("role", "admin")
		return ctrl.FindById(c)
	})

	req := httptest.NewRequest("GET", "/users/not-uuid", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUserController_FindById_ServiceError(t *testing.T) {
	mockService := new(UserServiceMock)
	targetId := uuid.New().String()

	mockService.On("FindById", mock.Anything, targetId).Return(domain.User{}, errors.New("database error"))

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Get("/users/:userId", func(c *fiber.Ctx) error {
		c.Locals("role", "admin")
		return ctrl.FindById(c)
	})

	req := httptest.NewRequest("GET", "/users/"+targetId, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUserController_Me_Success(t *testing.T) {
	mockService := new(UserServiceMock)
	targetId := uuid.New().String()

	expectedUser := domain.User{
		Id:       uuid.MustParse(targetId),
		Email:    "test@example.com",
		FullName: "Test User",
		Role:     "user",
	}

	mockService.On("FindById", mock.Anything, targetId).Return(expectedUser, nil)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Get("/users/me", func(c *fiber.Ctx) error {
		c.Locals("userId", targetId)
		return ctrl.Me(c)
	})

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response map[string]any
	json.NewDecoder(resp.Body).Decode(&response)
	data := response["Data"].(map[string]any)
	assert.Equal(t, expectedUser.Email, data["email"])
	assert.Equal(t, expectedUser.FullName, data["full_name"])

	mockService.AssertExpectations(t)
}

func TestUserController_Me_InvalidUUID(t *testing.T) {
	mockService := new(UserServiceMock)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Get("/users/me", func(c *fiber.Ctx) error {
		c.Locals("userId", "not-uuid")
		return ctrl.Me(c)
	})

	req := httptest.NewRequest("GET", "/users/me", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUserController_Me_ServiceError(t *testing.T) {
	mockService := new(UserServiceMock)
	targetId := uuid.New().String()

	mockService.On("FindById", mock.Anything, targetId).Return(domain.User{}, errors.New("database error"))

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Get("/users/me", func(c *fiber.Ctx) error {
		c.Locals("userId", targetId)
		return ctrl.Me(c)
	})

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUserController_Delete_Success(t *testing.T) {
	mockService := new(UserServiceMock)
	id := uuid.New().String()
	mockService.On("Delete", mock.Anything, id).Return(nil)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Delete("/users/:userId", ctrl.Delete)

	req := httptest.NewRequest("DELETE", "/users/"+id, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var respMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respMap)
	assert.NoError(t, err)
	assert.Equal(t, "user deleted", respMap["message"])
	assert.Equal(t, id, respMap["id"])

	mockService.AssertExpectations(t)
}

func TestUserController_Delete_InvalidUUID(t *testing.T) {
	mockService := new(UserServiceMock)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Delete("/users/:userId", func(c *fiber.Ctx) error {
		return ctrl.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/users/not-uuid", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertNotCalled(t, "Delete")
}

func TestUserController_Delete_ServiceError(t *testing.T) {
	mockService := new(UserServiceMock)
	id := uuid.New().String()
	mockService.On("Delete", mock.Anything, id).Return(errors.New("failed to delete"))

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Delete("/users/:userId", ctrl.Delete)

	req := httptest.NewRequest("DELETE", "/users/"+id, nil)

	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUserController_FindAll_Success(t *testing.T) {
	mockService := new(UserServiceMock)
	users := []domain.User{{Email: "a@a.com"}}
	mockService.On("FindAll", mock.Anything).Return(users, nil)

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Get("/users", ctrl.FindAll)

	req := httptest.NewRequest("GET", "/users", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var respMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respMap)
	assert.NoError(t, err)
	assert.NotNil(t, respMap["data"])

	mockService.AssertExpectations(t)
}

func TestUserController_FindAll_ServiceError(t *testing.T) {
	mockService := new(UserServiceMock)

	mockService.On("FindAll", mock.Anything).Return([]domain.User{}, errors.New("failed"))

	app := fiber.New()
	ctrl := controller.NewUserController(mockService)
	app.Get("/users", func(c *fiber.Ctx) error {
		return ctrl.FindAll(c)
	})

	req := httptest.NewRequest("GET", "/users", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}
