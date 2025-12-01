package controller

import "github.com/gofiber/fiber/v2"

type UserController interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	UpdateMe(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	FindById(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
}
