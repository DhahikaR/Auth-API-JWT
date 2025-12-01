package test

import (
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
	"auth-api-jwt/service"
	"context"
	"errors"
	"testing"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Save(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepositoryMock) Delete(ctx context.Context, tx *gorm.DB, userId string) error {
	args := m.Called(ctx, tx, userId)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindById(ctx context.Context, tx *gorm.DB, userId string) (domain.User, error) {
	args := m.Called(ctx, tx, userId)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepositoryMock) FindAll(ctx context.Context, tx *gorm.DB) ([]domain.User, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *UserRepositoryMock) UpdateLastLogin(ctx context.Context, tx *gorm.DB, userId string, loginAt time.Time) error {
	args := m.Called(ctx, tx, userId, loginAt)
	return args.Error(0)
}

func TestUserService_Create_Success(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	request := web.UserCreateRequest{
		Email:    "test@example.com",
		Password: "password",
		FullName: "Test User",
		Role:     "user",
	}

	expected := domain.User{
		Id:           uuid.New(),
		Email:        request.Email,
		PasswordHash: request.Password,
		FullName:     request.FullName,
		Role:         request.Role,
	}

	mockRepo.On("Save", mock.Anything, mock.Anything, mock.AnythingOfType("domain.User")).Return(expected, nil)

	svc := service.NewUserService(mockRepo, db, validate)
	got, err := svc.Create(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, expected.Email, got.Email)
}

func TestUserService_Create_Failed(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	request := web.UserCreateRequest{
		Email:    "",
		Password: "",
		FullName: "",
		Role:     "",
	}

	svc := service.NewUserService(mockRepo, db, validate)

	assert.Panics(t, func() {
		svc.Create(context.Background(), request)
	})

	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserService_Create_RepoSaveError(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	request := web.UserCreateRequest{
		Email:    "test@example.com",
		Password: "password",
		FullName: "Test User",
		Role:     "user",
	}

	mockRepo.On("Save", mock.Anything, mock.Anything, mock.AnythingOfType("domain.User")).Return(domain.User{}, errors.New("db error"))

	svc := service.NewUserService(mockRepo, db, validate)
	_, err := svc.Create(context.Background(), request)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New()

	existing := domain.User{
		Id:           id,
		Email:        "old@example.com",
		PasswordHash: "oldhash",
		FullName:     "Old",
		Role:         "user",
	}

	request := web.UserUpdateRequest{
		Id:           id,
		Email:        "updated@example.com",
		FullName:     "Updated",
		PasswordHash: "",
		Role:         "",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, id.String()).Return(existing, nil)

	updated := existing
	updated.Email = request.Email
	updated.FullName = request.FullName
	updated.Role = request.Role

	mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("domain.User")).Return(updated, nil)

	svc := service.NewUserService(mockRepo, db, validate)
	got, err := svc.Update(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", got.FullName)
}

func TestUserService_Update_FindByIdError(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New()

	request := web.UserUpdateRequest{
		Id:       id,
		FullName: "Test Name",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, id.String()).Return(domain.User{}, assert.AnError)

	svc := service.NewUserService(mockRepo, db, validate)

	result, err := svc.Update(context.Background(), request)
	assert.Error(t, err)
	assert.Empty(t, result)

	mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserService_Update_RepoUpdateError(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New()

	request := web.UserUpdateRequest{
		Id:       id,
		FullName: "Test Name",
	}

	existing := domain.User{
		Id:    id,
		Email: "test@example.com",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, id.String()).Return(existing, nil)

	mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("domain.User")).Return(domain.User{}, errors.New("db update error"))

	svc := service.NewUserService(mockRepo, db, validate)
	_, err := svc.Update(context.Background(), request)
	assert.Error(t, err)
	assert.Equal(t, "db update error", err.Error())
}

func TestUserService_UpdateMe_Success(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New()

	existing := domain.User{
		Id:           id,
		Email:        "old@example.com",
		PasswordHash: "oldhash",
		FullName:     "Old",
		Role:         "user",
	}

	request := web.UserUpdateRequest{
		Id:       id,
		Email:    "updated@example.com",
		FullName: "Updated",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, id.String()).Return(existing, nil)

	updated := existing
	updated.Email = request.Email
	updated.FullName = request.FullName
	updated.Role = request.Role

	mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("domain.User")).Return(updated, nil)

	svc := service.NewUserService(mockRepo, db, validate)

	got, err := svc.UpdateMe(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", got.FullName)
	assert.Equal(t, "updated@example.com", got.Email)
}

func TestUserService_UpdateMe_FindByIdError(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New()

	request := web.UserUpdateRequest{
		Id:       id,
		FullName: "Test Name",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, id.String()).Return(domain.User{}, assert.AnError)

	svc := service.NewUserService(mockRepo, db, validate)

	result, err := svc.UpdateMe(context.Background(), request)
	assert.Error(t, err)
	assert.Empty(t, result)

	mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserService_UpdateMe_UpdateError(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New()

	request := web.UserUpdateRequest{
		Id:       id,
		FullName: "Test Name",
	}

	existing := domain.User{
		Id:    id,
		Email: "test@example.com",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, id.String()).Return(existing, nil)

	mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("domain.User")).Return(domain.User{}, errors.New("db update error"))

	svc := service.NewUserService(mockRepo, db, validate)
	_, err := svc.UpdateMe(context.Background(), request)
	assert.Error(t, err)
	assert.Equal(t, "db update error", err.Error())
}

func TestUserService_Delete_Success(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New()

	existing := domain.User{
		Id:           id,
		Email:        "todelete@example.com",
		PasswordHash: "hash",
		FullName:     "To Delete",
		Role:         "user",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, existing.Id.String()).Return(existing, nil)

	mockRepo.On("Delete", mock.Anything, mock.Anything, existing.Id.String()).Return(nil)

	svc := service.NewUserService(mockRepo, db, validate)
	err := svc.Delete(context.Background(), existing.Id.String())
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_NotFound(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New().String()

	mockRepo.On("FindById", mock.Anything, mock.Anything, id).Return(domain.User{}, assert.AnError)

	svc := service.NewUserService(mockRepo, db, validate)
	err := svc.Delete(context.Background(), id)
	if err == nil {
		t.Fatalf("expected error when finding non-existent user")
	}

	mockRepo.AssertExpectations(t)
}
func TestUserService_Delete_RepoDeleteError(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New().String()

	existing := domain.User{Id: uuid.MustParse(id)}

	mockRepo.On("FindById", mock.Anything, mock.Anything, id).Return(existing, nil)

	mockRepo.On("Delete", mock.Anything, mock.Anything, id).Return(errors.New("delete failed"))

	svc := service.NewUserService(mockRepo, db, validate)

	err := svc.Delete(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, "delete failed", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindById_Success(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	existing := domain.User{
		Id:       uuid.New(),
		Email:    "test@example.com",
		FullName: "To Delete",
		Role:     "user",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, existing.Id.String()).Return(existing, nil)

	svc := service.NewUserService(mockRepo, db, validate)
	result, err := svc.FindById(context.Background(), existing.Id.String())
	assert.NoError(t, err)
	assert.Equal(t, existing.Id.String(), result.Id.String())
	assert.Equal(t, existing.Email, result.Email)
	assert.Equal(t, existing.FullName, result.FullName)
	assert.Equal(t, existing.Role, result.Role)

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindById_NotFound(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New().String()

	mockRepo.On("FindById", mock.Anything, mock.Anything, id).Return(domain.User{}, assert.AnError)

	svc := service.NewUserService(mockRepo, db, validate)
	_, err := svc.FindById(context.Background(), id)
	if err == nil {
		t.Fatalf("expected error when finding non-existent user")
	}

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindById_RepoFindByIdError(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	id := uuid.New().String()

	mockRepo.On("FindById", mock.Anything, mock.Anything, id).Return(domain.User{}, errors.New("database error"))

	svc := service.NewUserService(mockRepo, db, validate)

	_, err := svc.FindById(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindAll_Success(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	existing := []domain.User{}

	mockRepo.On("FindAll", mock.Anything, mock.Anything).Return(existing, nil)

	svc := service.NewUserService(mockRepo, db, validate)
	result, err := svc.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, existing, result)

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindAll_RepoFindAllError(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	db := setupTestDB(t)
	validate := validator.New()

	existing := []domain.User{}

	mockRepo.On("FindAll", mock.Anything, mock.Anything).Return(existing, assert.AnError)

	svc := service.NewUserService(mockRepo, db, validate)
	_, err := svc.FindAll(context.Background())
	if err == nil {
		t.Fatalf("expected error when finding non-existent user")
	}

	mockRepo.AssertExpectations(t)
}
