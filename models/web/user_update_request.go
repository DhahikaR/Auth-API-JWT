package web

import "github.com/google/uuid"

type UserUpdateRequest struct {
	Id           uuid.UUID
	Email        string `validate:"omitempty,email"`
	PasswordHash string `validate:"omitempty,min=6"`
	FullName     string `validate:"omitempty"`
	Role         string `validate:"omitempty,oneof=user admin"`
}
