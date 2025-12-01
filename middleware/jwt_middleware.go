package middleware

import (
	"auth-api-jwt/helper"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.Unauthorized(c, "missing authorized header")
		}

		fields := strings.Split(authHeader, " ")
		if len(fields) != 2 || fields[0] != "Bearer" {
			return helper.Unauthorized(c, "invalid authorization format")
		}

		tokenString := fields[1]
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token signature")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return helper.Unauthorized(c, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return helper.Unauthorized(c, "invalid token claims")
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			return helper.Unauthorized(c, "invalid user id in token")
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			return helper.Unauthorized(c, "invalid role in token")
		}

		c.Locals("userId", userId)
		c.Locals("role", role)

		return c.Next()
	}
}
