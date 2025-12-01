package helper

import (
	"auth-api-jwt/models/domain"
	"auth-api-jwt/models/web"
)

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		Id:       user.Id,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	}
}

func ToUserResponses(users []domain.User) []web.UserResponse {
	var userResponses []web.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ToUserResponse(user))
	}

	return userResponses
}
