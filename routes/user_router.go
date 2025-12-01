package routes

import (
	"auth-api-jwt/controller"
	"auth-api-jwt/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewUserRouter(app *fiber.App, userController controller.UserController) {
	user := app.Group("/users", middleware.JWTMiddleware())

	user.Put("/me", userController.UpdateMe)
	user.Get("/me", userController.Me)

	admin := user.Group("/", middleware.AdminOnly())

	admin.Get("/", userController.FindAll)
	admin.Get("/:userId", userController.FindById)
	admin.Post("/", userController.Create)
	admin.Put("/:userId", userController.Update)
	admin.Delete("/:userId", userController.Delete)
}
