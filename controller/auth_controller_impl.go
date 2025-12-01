package controller

import (
	"auth-api-jwt/helper"
	"auth-api-jwt/models/web"
	"auth-api-jwt/service"

	"github.com/gofiber/fiber/v2"
)

type AuthControllerImpl struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &AuthControllerImpl{
		authService: authService,
	}
}

func (controller *AuthControllerImpl) Register(c *fiber.Ctx) error {
	authRegisterRequest := web.AuthRegisterRequest{}
	if err := helper.ReadFromRequestBody(c, &authRegisterRequest); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	user, err := controller.authService.Register(c.Context(), authRegisterRequest)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	response := helper.ToUserResponse(user)

	return helper.ResponseSuccess(c, response)
}

func (controller *AuthControllerImpl) Login(c *fiber.Ctx) error {
	authLoginRequest := web.AuthLoginRequest{}
	if err := helper.ReadFromRequestBody(c, &authLoginRequest); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	token, err := controller.authService.Login(c.Context(), authLoginRequest)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	return helper.ResponseSuccess(c, fiber.Map{
		"token": token,
	})
}
