package helper

import (
	"auth-api-jwt/models/web"

	"github.com/gofiber/fiber/v2"
)

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
		Code:   fiber.StatusBadRequest,
		Status: "BAD REQUEST",
		Data:   message,
	})
}

func ResponseSuccess(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data,
	})
}

func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"code":   fiber.StatusForbidden,
		"status": "FORBIDDEN",
		"data":   message,
	})
}
