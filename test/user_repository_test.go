package test

import (
	"auth-api-jwt/models/domain"
	"auth-api-jwt/repository"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type testUser struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"type:varchar(255);unique;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	FullName     string    `gorm:"type:varchar(100);not null"`
	IsVerified   bool      `gorm:"default:false"`
	Role         string    `gorm:"type:varchar(50);default:'user'"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (testUser) TableName() string { return "users" }

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	if err := db.AutoMigrate(&testUser{}); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}

	return db
}

func TestUserRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	// Create
	u := domain.User{
		Id:           uuid.New(),
		Email:        "repo-test@example.com",
		PasswordHash: "hash",
		FullName:     "Repo Test",
		Role:         "user",
	}

	saved, err := repo.Save(ctx, db, u)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, saved.Email)
	assert.NotEqual(t, uuid.UUID{}, saved.Id)

	// FindById
	found, err := repo.FindById(ctx, db, saved.Id.String())
	assert.NoError(t, err)
	assert.Equal(t, saved.Email, found.Email)

	// FindAll
	all, err := repo.FindAll(ctx, db)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(all), 1)

	// UpdateLastLogin
	now := time.Now()
	err = repo.UpdateLastLogin(ctx, db, saved.Id.String(), now)
	assert.NoError(t, err)

	updated, err := repo.FindById(ctx, db, saved.Id.String())
	assert.NoError(t, err)
	if assert.NotNil(t, updated.LastLoginAt) {
		assert.WithinDuration(t, now, *updated.LastLoginAt, time.Second)
	}

	// Update
	updated.FullName = "Repo Test Updated"
	up, err := repo.Update(ctx, db, updated)
	assert.NoError(t, err)
	assert.Equal(t, "Repo Test Updated", up.FullName)

	// Delete (soft delete)
	err = repo.Delete(ctx, db, saved.Id.String())
	assert.NoError(t, err)

	// After delete, FindById should return error (record not found)
	_, err = repo.FindById(ctx, db, saved.Id.String())
	if err == nil {
		t.Fatalf("expected record to be not found after delete, but it was found")
	}
}
