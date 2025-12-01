package middleware

import (
	"auth-api-jwt/helper"

	"github.com/gofiber/fiber/v2"
)

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return helper.Forbidden(c, "role not found in token")
		}

		if role != "admin" {
			return helper.Forbidden(c, "admin only endpoint")
		}

		return c.Next()
	}
}
