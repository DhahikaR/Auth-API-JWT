package test

import (
	"auth-api-jwt/helper"
	"auth-api-jwt/middleware"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/assert"
)

func TestAdminOnly_Success(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "admin")
		return c.Next()
	})

	app.Use(middleware.AdminOnly())
	app.Get("/admin", func(c *fiber.Ctx) error {
		return helper.ResponseSuccess(c, "OK")
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminOnly_NoRole(t *testing.T) {
	app := fiber.New()

	app.Use(middleware.AdminOnly())
	app.Get("/admin", func(c *fiber.Ctx) error {
		return helper.ResponseSuccess(c, "OK")
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

func TestAdminOnly_Forbidden(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "user")
		return c.Next()
	})

	app.Use(middleware.AdminOnly())
	app.Get("/admin", func(c *fiber.Ctx) error {
		return helper.ResponseSuccess(c, "OK")
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}
