package helper

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateEmail(email string) error {
	err := validate.Var(email, "required,email")
	if err != nil {
		return errors.New("invalid email format")
	}
	return nil
}
