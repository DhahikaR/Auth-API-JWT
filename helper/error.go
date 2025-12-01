package helper

import "github.com/gofiber/fiber/v2"

func PanicIfError(err error) {
	if err != nil {
		panic(fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
}

func ErrorResponse(err error) fiber.Map {
	return fiber.Map{
		"error": err.Error(),
	}
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status":  "UNAUTHORIZED",
		"message": message,
	})
}
