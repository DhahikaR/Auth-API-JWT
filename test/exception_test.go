package test

import (
	"auth-api-jwt/exception"
	"auth-api-jwt/models/web"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	return fiber.New(fiber.Config{ErrorHandler: exception.NewErrorHandler})
}

func decodeResponse(t *testing.T, response *http.Response) web.WebResponse {
	var webResponse web.WebResponse
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&webResponse); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return webResponse
}

func TestErrorHandler_ValidationError(t *testing.T) {
	app := setupApp()

	app.Get("/validate", func(c *fiber.Ctx) error {
		v := validator.New()
		type R struct {
			Name string `validate:"required"`
		}
		err := v.Struct(R{})
		return err
	})

	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	writer := decodeResponse(t, resp)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "BAD REQUEST", writer.Status)
}

func TestErrorHandler_NotFoundError(t *testing.T) {
	app := setupApp()

	app.Get("/notfound", func(c *fiber.Ctx) error {
		return exception.NotFoundError{Message: "user not found"}
	})

	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	writer := decodeResponse(t, resp)
	assert.Equal(t, http.StatusNotFound, writer.Code)
	assert.Equal(t, "NOT FOUND", writer.Status)

	if message, ok := writer.Data.(string); ok {
		assert.Equal(t, "user not found", message)
	} else {
		t.Fatalf("expected Data to be string, got %T", writer.Data)
	}
}

func TestErrorHandler_ValidationMultipleFields(t *testing.T) {
	app := setupApp()

	app.Get("/vm", func(c *fiber.Ctx) error {
		v := validator.New()
		type R struct {
			Name     string `validate:"required"`
			Password string `validate:"required,min=6"`
		}
		err := v.Struct(R{})
		return err
	})

	req := httptest.NewRequest(http.MethodGet, "/vm", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	writer := decodeResponse(t, resp)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "BAD REQUEST", writer.Status)

	if s, ok := writer.Data.(string); ok {
		assert.Contains(t, s, "Name")
		assert.Contains(t, s, "Password")
	} else {
		t.Fatalf("expected Data to be string, got %T", writer.Data)
	}
}

func TestErrorHandler_FiberError(t *testing.T) {
	app := setupApp()

	app.Get("/fe", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	})

	req := httptest.NewRequest(http.MethodGet, "/fe", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	writer := decodeResponse(t, resp)
	assert.Equal(t, http.StatusNotFound, writer.Code)
	assert.Equal(t, "NOT FOUND", writer.Status)
}

func TestErrorHandler_InternalServerError(t *testing.T) {
	app := setupApp()

	app.Get("/ise", func(c *fiber.Ctx) error {
		return errors.New("error")
	})

	req := httptest.NewRequest(http.MethodGet, "/ise", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	writer := decodeResponse(t, resp)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "INTERNAL SERVICE ERROR", writer.Status)
}
