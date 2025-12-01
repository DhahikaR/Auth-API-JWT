package main

import (
	"auth-api-jwt/config"
	"auth-api-jwt/controller"
	"auth-api-jwt/exception"

	_ "auth-api-jwt/docs"
	"auth-api-jwt/repository"
	"auth-api-jwt/routes"
	"auth-api-jwt/service"

	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: exception.NewErrorHandler,
	})

	db := config.NewDB()
	validate := validator.New()

	userRepository := repository.NewUserRepository(db)
	authRepository := repository.NewAuthRepository(db)

	userService := service.NewUserService(userRepository, db, validate)
	authService := service.NewAuthService(authRepository, userRepository, db, validate)

	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)

	routes.NewUserRouter(app, userController)
	routes.NewAuthRoutes(app, authController)

	app.Listen(":3000")

}
