package test

import (
	"auth-api-jwt/helper"
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPanicIfError(t *testing.T) {
	helper.PanicIfError(nil)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic")
		}

		if fiberErr, ok := r.(*fiber.Error); ok {
			assert.Equal(t, fiber.StatusBadRequest, fiberErr.Code)
			assert.Equal(t, "error", fiberErr.Message)
			return
		}
		t.Fatalf("unexpected panic type: %T", r)
	}()

	helper.PanicIfError(errors.New("error"))
}

func TestUserResponseHelpers(t *testing.T) {
	user := domain.User{Id: uuid.New(), Email: "test@example.com", FullName: "Test", Role: "user"}

	resp := helper.ToUserResponse(user)

	assert.Equal(t, resp.Id, user.Id)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "Test", resp.FullName)
	assert.Equal(t, "user", resp.Role)

	users := []domain.User{user, {Id: uuid.New(), Email: "test2@example.com", FullName: "Test 2", Role: "user"}}

	respList := helper.ToUserResponses(users)
	assert.Len(t, respList, 2)
	assert.Equal(t, respList[1].Email, "test2@example.com")
}

func TestResponseAndReadFromRequestBody(t *testing.T) {
	app := fiber.New()

	app.Post("/parse", func(c *fiber.Ctx) error {
		var req web.UserCreateRequest
		if err := helper.ReadFromRequestBody(c, &req); err != nil {
			return helper.BadRequest(c, err.Error())
		}
		return helper.ResponseSuccess(c, req)
	})

	// valid JSON
	body, _ := json.Marshal(web.UserCreateRequest{Email: "test@example.com", Password: "test", FullName: "Test"})
	req := httptest.NewRequest(http.MethodPost, "/parse", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// invalid JSON
	req2 := httptest.NewRequest(http.MethodPost, "/parse", bytes.NewReader([]byte(`{invalid}`)))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2, -1)

	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)
}

func TestCommitOrRollback(t *testing.T) {
	db := setupTestDB(t)

	// commit case
	tx := db.Begin()
	tx.Create(&domain.User{
		Id:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "test"})
	helper.CommitOrRollback(tx)
	var count int64
	db.Model(&domain.User{}).Count(&count)
	assert.EqualValues(t, int64(2), count)

	// rollback case: simulate panic inside function
	func() {
		tx2 := db.Begin()
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		defer helper.CommitOrRollback(tx2)
		tx2.Create(&domain.User{
			Id:           uuid.New(),
			Email:        "test",
			PasswordHash: "t12"})
		panic("error")
	}()

	db.Model(&domain.User{}).Count(&count)

	assert.EqualValues(t, int64(2), count)
}

func TestValidateEmail(t *testing.T) {

	// valid email
	tests := []string{
		"user@example.com",
		"abcdef@gmail.com",
		"usertag@gmail.com",
		"u@x.com",
	}

	for _, email := range tests {
		err := helper.ValidateEmail(email)

		assert.NoError(t, err, "expected valid email")
	}

	// invalid email
	invalidTests := []string{
		"userexample.com",
		"abcdefgmail.com",
		"usertag@gmail",
		"u@.com",
		"@example.com",
	}

	for _, email := range invalidTests {
		err := helper.ValidateEmail(email)

		assert.Error(t, err, "expected invalid email")
	}
}
