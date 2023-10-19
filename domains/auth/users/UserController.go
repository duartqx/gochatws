package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	e "github.com/duartqx/gochatws/core/errors"
)

type UserController struct {
	userRepository *UserRepository
}

func NewUserController(ur *UserRepository) *UserController {
	return &UserController{userRepository: ur}
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
	dbUser, err := uc.userRepository.FindByIdParam(c.Params("id"))
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

	if uc.userRepository.ExistsByUsername(bodyUser.Username) == true {
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

	dbUser, err := uc.userRepository.FindByIdParam(c.Params("id"))
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(e.CustomMessageError(err.Error()))
	}

	bodyUser, err, validationErrs := uc.userRepository.ParseAndValidate(c.BodyParser)

	if validationErrs != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(e.ValidationError(validationErrs))
	}

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(e.SerializerError)
	}

	dbUser.UpdateFromAnother(bodyUser)

	uc.userRepository.Update(dbUser)

	return c.Status(http.StatusOK).JSON(dbUser.Clean())
}

func (uc UserController) Delete(c *fiber.Ctx) error {
	return nil
}
