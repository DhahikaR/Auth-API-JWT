package routes

import (
	"auth-api-jwt/controller"

	"github.com/gofiber/fiber/v2"
)

func NewAuthRoutes(app *fiber.App, authController controller.AuthController) {
	auth := app.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
}
