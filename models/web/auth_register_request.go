package web

type AuthRegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	FullName string `validate:"required"`
}
