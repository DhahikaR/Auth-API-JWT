package web

type AuthLoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
