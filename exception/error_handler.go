package exception

import (
	"auth-api-jwt/models/web"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewErrorHandler(c *fiber.Ctx, err error) error {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   validationError.Error(),
		})
	}

	if notFound, ok := err.(NotFoundError); ok {
		return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{
			Code:   fiber.StatusNotFound,
			Status: "NOT FOUND",
			Data:   notFound.Error(),
		})
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		code := fiberErr.Code
		if code == 0 {
			code = fiber.StatusInternalServerError
		}
		statusText := "ERROR"
		if code == fiber.StatusBadRequest {
			statusText = "BAD REQUEST"
		} else if code == fiber.StatusNotFound {
			statusText = "NOT FOUND"
		} else if code == fiber.StatusInternalServerError {
			statusText = "INTERNAL SERVICE ERROR"
		}
		return c.Status(code).JSON(web.WebResponse{
			Code:   code,
			Status: statusText,
			Data:   fiberErr.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
		Code:   fiber.StatusInternalServerError,
		Status: "INTERNAL SERVICE ERROR",
		Data:   err.Error(),
	})
}
