package auth

import (
	cerr "gochatws/core/errors"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
			JSON(cerr.InternalError)
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
		return c.Status(http.StatusNotFound).JSON(cerr.NotFoundError)
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

	bodyUser, err, validationErrs := uc.userRepository.Validate(c.BodyParser)

	if validationErrs != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(cerr.ValidationError(validationErrs))
	}

	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(cerr.SerializerError)
	}

	if uc.userRepository.ExistsByUsername(bodyUser.Username) == true {
		return c.
			Status(http.StatusBadRequest).
			JSON(cerr.InvalidUsernameError)
	}

	hashedPassword, err :=
		bcrypt.GenerateFromPassword([]byte(bodyUser.Password), 10)

	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(cerr.PasswordTooLongError)
	}
	bodyUser.Password = string(hashedPassword)

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

	dbUser, err := uc.userRepository.FindByIdParam(c.Params("id"))
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(cerr.CustomMessageError(err.Error()))
	}

	bodyUser, err, validationErrs := uc.userRepository.Validate(c.BodyParser)

	if validationErrs != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(cerr.ValidationError(validationErrs))
	}

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(cerr.SerializerError)
	}

	dbUser.UpdateFromAnother(bodyUser)

	uc.userRepository.Update(dbUser)

	return c.Status(http.StatusOK).JSON(dbUser)
}

func (uc UserController) Delete(c *fiber.Ctx) error {
	return nil
}
