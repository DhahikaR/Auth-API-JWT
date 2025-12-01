package web

type UserCreateRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	FullName string `validate:"required"`
	Role     string `validate:"omitempty,oneof=user admin"`
}
