package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	e "github.com/duartqx/gochatws/core/errors"
	i "github.com/duartqx/gochatws/core/interfaces"
)

type UserController struct {
	userRepository *UserRepository
}

func NewUserController(ur *UserRepository) *UserController {
	return &UserController{userRepository: ur}
}

func (uc UserController) getUserFromLocals(localUser interface{}) (i.User, error) {
	if localUser == nil {
		return nil, fmt.Errorf("User not found on Locals\n")
	}
	userBytes, err := json.Marshal(localUser)
	if err != nil {
		return nil, err
	}

	userStruct := uc.userRepository.GetModel()
	err = json.Unmarshal(userBytes, userStruct)
	if err != nil {
		return nil, err
	}
	return userStruct, nil
}

func (uc UserController) All(c *fiber.Ctx) error {

	users, err := uc.userRepository.All()
	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(e.InternalError)
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
	user, err := uc.getUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.UnauthorizedError)
	}

	dbUser, err := uc.userRepository.FindById(user.GetId())
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(e.NotFoundError)
	}
	return c.Status(http.StatusOK).JSON(dbUser.Clean())
}

/*
*
* @Body -> UserModel
* @Returns -> error
*
 */
func (uc UserController) Create(c *fiber.Ctx) error {

	bodyUser, err, validationErrs := uc.userRepository.ParseAndValidate(c.BodyParser)

	if validationErrs != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(e.ValidationError(validationErrs))
	}

	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(e.SerializerError)
	}

	if uc.userRepository.ExistsByUsername(bodyUser.Username) {
		return c.
			Status(http.StatusBadRequest).
			JSON(e.InvalidUsernameError)
	}

	hashedPassword, err :=
		bcrypt.GenerateFromPassword([]byte(bodyUser.Password), 10)

	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(e.PasswordTooLongError)
	}
	bodyUser.Password = string(hashedPassword)

	uc.userRepository.Create(bodyUser)

	return c.Status(http.StatusCreated).JSON(bodyUser.Clean())
}

/*
*
* @Params -> :id<int>
* @Returns -> error
*
 */
func (uc UserController) Update(c *fiber.Ctx) error {
	// TODO: Remove id Param and get the id of the authenticated user

	user, err := uc.getUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.UnauthorizedError)
	}

	dbUser, err := uc.userRepository.FindById(user.GetId())
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(e.CustomMessageError(err.Error()))
	}

	bodyUser, err := uc.userRepository.Parse(c.BodyParser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.SerializerError)
	}

	dbUser.UpdateFromAnother(bodyUser)

	uc.userRepository.Update(dbUser)

	return c.Status(http.StatusOK).JSON(dbUser.Clean())
}

func (uc UserController) Delete(c *fiber.Ctx) error {

	user, err := uc.getUserFromLocals(c.Locals("user"))
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.UnauthorizedError)
	}

	err = uc.userRepository.Delete(user)
	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(e.InternalError)
	}

	resp := fmt.Sprintf("Successfully deleted user with id: %d", user.GetId())

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{"user": resp})
}
