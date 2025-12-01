package test

import (
	"auth-api-jwt/utils"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {

	password := "mypassword123"

	hashed, err := utils.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestCheckPassword_Success(t *testing.T) {
	password := "secret134"

	hashed, _ := utils.HashPassword(password)

	match := utils.CheckPassword(password, hashed)

	assert.True(t, match)
}

func TestCheckPassword_Failure(t *testing.T) {
	password := "secret321"
	wrongPassword := "incorrect"

	hashed, _ := utils.HashPassword(password)

	match := utils.CheckPassword(hashed, wrongPassword)

	assert.False(t, match, "password not match")

}

func TestJWTGeneration(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecretkey")

	tokenString, err := utils.GenerateJWT("12345", "admin")

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify its claims
	parseToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("testsecretkey"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parseToken.Valid)

	claims, ok := parseToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, "12345", claims["user_id"])
	assert.Equal(t, "admin", claims["role"])
	assert.NotNil(t, claims["exp"])
	assert.NotNil(t, claims["iat"])
}
