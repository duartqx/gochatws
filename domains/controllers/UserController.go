package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	e "github.com/duartqx/gochatws/core/errors"
	m "github.com/duartqx/gochatws/domains/models"
	s "github.com/duartqx/gochatws/domains/services"
	"github.com/duartqx/gochatws/domains/utils"
)

type UserController struct {
	userService *s.UserService
}

func NewUserController(us *s.UserService) *UserController {
	return &UserController{userService: us}
}

func (uc UserController) All(c *fiber.Ctx) error {
	response := uc.userService.All()
	return c.Status(response.Status).JSON(response.Body)
}

func (uc UserController) Get(c *fiber.Ctx) error {
	user, err := utils.GetUserFromLocals(c.Locals("user"))

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(e.UnauthorizedError)
	}
	response := uc.userService.Get(user.GetId())
	return c.Status(response.Status).JSON(response.Body)
}

func (uc UserController) Create(c *fiber.Ctx) error {
	user := &m.UserModel{}

	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.BadRequestError)
	}

	response := uc.userService.Create(user)
	return c.Status(response.Status).JSON(response.Body)
}

func (uc UserController) Update(c *fiber.Ctx) error {

	bodyUser := &m.UserModel{}
	if err := c.BodyParser(bodyUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.BadRequestError)
	}

	user, err := utils.GetUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(e.UnauthorizedError)
	}

	bodyUser.SetId(user.GetId())

	response := uc.userService.Update(bodyUser)
	return c.Status(response.Status).JSON(response.Body)
}

func (uc UserController) Delete(c *fiber.Ctx) error {

	user, err := utils.GetUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.UnauthorizedError)
	}

	response := uc.userService.Delete(user)
	return c.Status(response.Status).JSON(response.Body)
}
