package controller

import (
	"auth-api-jwt/helper"
	"auth-api-jwt/models/web"
	"auth-api-jwt/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserControllerImpl struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		userService: userService,
	}
}

func (controller *UserControllerImpl) Create(c *fiber.Ctx) error {
	userCreateRequest := web.UserCreateRequest{}
	if err := helper.ReadFromRequestBody(c, &userCreateRequest); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	if userCreateRequest.Email == "" || userCreateRequest.Password == "" || userCreateRequest.FullName == "" {
		return helper.BadRequest(c, "missing required fields")
	}

	user, err := controller.userService.Create(c.Context(), userCreateRequest)
	if err != nil {
		return c.Status(500).JSON(helper.ErrorResponse(err))
	}

	response := helper.ToUserResponse(user)

	return helper.ResponseSuccess(c, response)
}

func (controller *UserControllerImpl) Update(c *fiber.Ctx) error {
	targetUserId := c.Params("userId")

	if _, err := uuid.Parse(targetUserId); err != nil {
		return helper.BadRequest(c, "invalid UUID")
	}

	request := web.UserUpdateRequest{}
	if err := helper.ReadFromRequestBody(c, &request); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	request.Id = uuid.MustParse(targetUserId)

	user, err := controller.userService.Update(c.Context(), request)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	return helper.ResponseSuccess(c, helper.ToUserResponse(user))
}

func (controller *UserControllerImpl) UpdateMe(c *fiber.Ctx) error {
	authUserId := c.Locals("userId").(string)

	request := web.UserUpdateRequest{}
	if err := helper.ReadFromRequestBody(c, &request); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	request.Id = uuid.MustParse(authUserId)
	request.Role = ""

	user, err := controller.userService.UpdateMe(c.Context(), request)
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	return helper.ResponseSuccess(c, helper.ToUserResponse(user))
}

func (controller *UserControllerImpl) Delete(c *fiber.Ctx) error {
	userId := c.Params("userId")
	if _, err := uuid.Parse(userId); err != nil {
		return helper.BadRequest(c, "invalid UUID")
	}

	err := controller.userService.Delete(c.Context(), userId)
	if err != nil {
		return helper.BadRequest(c, "invalid UUID")
	}

	return c.JSON(fiber.Map{
		"message": "user deleted",
		"id":      userId,
	})
}

func (controller *UserControllerImpl) FindById(c *fiber.Ctx) error {
	userId := c.Params("userId")

	id, err := uuid.Parse(userId)
	if err != nil {
		return helper.BadRequest(c, "invalid UUID")
	}

	user, err := controller.userService.FindById(c.Context(), id.String())
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	response := helper.ToUserResponse(user)

	return helper.ResponseSuccess(c, response)
}

func (controller *UserControllerImpl) Me(c *fiber.Ctx) error {
	authUserId := c.Locals("userId").(string)

	uid, err := uuid.Parse(authUserId)
	if err != nil {
		return helper.BadRequest(c, "invalid user id")
	}

	user, err := controller.userService.FindById(c.Context(), uid.String())
	if err != nil {
		return helper.BadRequest(c, err.Error())
	}

	return helper.ResponseSuccess(c, helper.ToUserResponse(user))
}

func (controller *UserControllerImpl) FindAll(c *fiber.Ctx) error {
	users, err := controller.userService.FindAll(c.Context())
	if err != nil {
		return helper.BadRequest(c, "bad request")
	}

	response := helper.ToUserResponses(users)

	return c.JSON(fiber.Map{"data": response})
}
