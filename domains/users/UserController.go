package users

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type UserController struct {
	userRepository *UserRepository
}

func NewUserController(userRepository *UserRepository) *UserController {
	return &UserController{userRepository: userRepository}
}

func (uc UserController) All(c *fiber.Ctx) error {
	users, err := uc.userRepository.All()
	if err != nil {
		slog.Error(err.Error())
		json_res := map[string]string{"error": "internal"}
		return c.Status(http.StatusInternalServerError).JSON(json_res)
	}
	return c.JSON(users)
}

/*
*
* @Params -> :id<int>
* @Returns -> error
*
 */
func (uc UserController) Get(c *fiber.Ctx) error {
	dbUser, err := uc.findUserByIdParam(c.Params("id"))
	if err != nil {
		slog.Error(err.Error())
		json_res := map[string]string{"error": "User Not Found"}
		return c.Status(http.StatusNotFound).JSON(json_res)
	}
	return c.Status(http.StatusOK).JSON(dbUser)
}

/*
*
* @Body -> UserModel
* @Returns -> error
*
 */
func (uc UserController) Create(c *fiber.Ctx) error {

	bodyUser, err := uc.parseAndValidate(c)
	if err != nil {
		return err
	}

	if uc.userRepository.ExistsByUsername(bodyUser.Username) == true {
		return c.Status(400).JSON(map[string]string{"error": "Invalid username"})
	}

	uc.userRepository.Create(bodyUser)

	return c.Status(http.StatusCreated).JSON(&bodyUser)
}

/*
*
* @Params -> :id<int>
* @Returns -> error
*
 */
func (uc UserController) Update(c *fiber.Ctx) error {
	// TODO: Remove id Param and get the id of the authenticated user

	dbUser, err := uc.findUserByIdParam(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
	}

	bodyUser, err := uc.parseAndValidate(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
	}

	dbUser.UpdateFromAnother(bodyUser)

	uc.userRepository.Update(dbUser)

	return c.Status(http.StatusOK).JSON(dbUser)
}

func (uc UserController) Delete(c *fiber.Ctx) error {
	return nil
}

// Util methods

func (uc UserController) parseAndValidate(c *fiber.Ctx) (*UserModel, error) {
	bodyUser := &UserModel{}
	if err := c.BodyParser(bodyUser); err != nil {
		return nil, err
	}
	if err := uc.userRepository.Validate(bodyUser); err != nil {
		return nil, err
	}
	return bodyUser, nil
}

func (uc UserController) findUserByIdParam(id string) (*UserModel, error) {
	var u *UserModel

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	u, err = uc.userRepository.FindById(idInt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
