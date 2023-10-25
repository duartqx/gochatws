package contr

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	e "github.com/duartqx/gochatws/core/errors"
	i "github.com/duartqx/gochatws/core/interfaces"
	u "github.com/duartqx/gochatws/domains/users"
)

type UserController struct {
	userService *u.UserService
}

func NewUserController(us *u.UserService) *UserController {
	return &UserController{userService: us}
}

func (uc UserController) getUserFromLocals(localUser interface{}) (i.User, error) {
	if localUser == nil {
		return nil, fmt.Errorf("User not found on Locals\n")
	}
	userBytes, err := json.Marshal(localUser)
	if err != nil {
		return nil, err
	}

	userStruct := &u.UserModel{}
	err = json.Unmarshal(userBytes, userStruct)
	if err != nil {
		return nil, err
	}
	return userStruct, nil
}

func (uc UserController) All(c *fiber.Ctx) error {
	response := uc.userService.All()
	return c.Status(response.Status).JSON(response.Body)
}

func (uc UserController) Get(c *fiber.Ctx) error {
	user, err := uc.getUserFromLocals(c.Locals("user"))

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(e.UnauthorizedError)
	}
	response := uc.userService.Get(user.GetId())
	return c.Status(response.Status).JSON(response.Body)
}

func (uc UserController) Create(c *fiber.Ctx) error {
	user := &u.UserModel{}

	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.BadRequestError)
	}

	response := uc.userService.Create(user)
	return c.Status(response.Status).JSON(response.Body)
}

func (uc UserController) Update(c *fiber.Ctx) error {

	bodyUser := &u.UserModel{}
	if err := c.BodyParser(bodyUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.BadRequestError)
	}

	user, err := uc.getUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(e.UnauthorizedError)
	}

	bodyUser.SetId(user.GetId())

	response := uc.userService.Update(bodyUser)
	return c.Status(response.Status).JSON(response.Body)
}

func (uc UserController) Delete(c *fiber.Ctx) error {

	user, err := uc.getUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.UnauthorizedError)
	}

	response := uc.userService.Delete(user)
	return c.Status(response.Status).JSON(response.Body)
}
