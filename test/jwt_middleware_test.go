package test

import (
	"auth-api-jwt/middleware"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(secret, userId, role string) string {
	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))

	return signed
}

func TestJWTMiddleware_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")

	app := fiber.New()
	app.Use(middleware.JWTMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	token := generateTestToken("testsecret", "123456", "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestJWTMiddleware_MissingHeader(t *testing.T) {

	app := fiber.New()
	app.Use(middleware.JWTMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/protected", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTMiddleware_InvalidFormat(t *testing.T) {

	app := fiber.New()
	app.Use(middleware.JWTMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Token abc1234")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTMiddleware_InvalidSignature(t *testing.T) {
	os.Setenv("JWT_SECRET", "true-secret")

	app := fiber.New()
	app.Use(middleware.JWTMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	token := generateTestToken("wrong-secret", "123456", "admin")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTMiddleware_MissingUserId(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")

	claims := jwt.MapClaims{
		"role": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("testsecret"))

	app := fiber.New()
	app.Use(middleware.JWTMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestJWTMiddleware_MissingRole(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")

	claims := jwt.MapClaims{
		"user_id": "12345",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("testsecret"))

	app := fiber.New()
	app.Use(middleware.JWTMiddleware())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
