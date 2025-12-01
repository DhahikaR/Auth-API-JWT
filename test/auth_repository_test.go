package test

import (
	"auth-api-jwt/models/domain"
	"auth-api-jwt/repository"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepository_CreateAndFindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuthRepository(db)
	ctx := context.Background()

	u := domain.User{
		Id:           uuid.New(),
		Email:        "auth-repo@example.com",
		PasswordHash: "secret",
		FullName:     "Auth Repo",
		Role:         "user",
	}

	created, err := repo.Create(ctx, db, u)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, created.Email)

	found, err := repo.FindByEmail(ctx, db, u.Email)
	assert.NoError(t, err)
	assert.Equal(t, created.Email, found.Email)

	// lookup missing email
	_, err = repo.FindByEmail(ctx, db, "no-such@example.com")
	if err == nil {
		t.Fatalf("expected error when finding non-existent email")
	}
}
