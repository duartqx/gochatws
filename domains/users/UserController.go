package users

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
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
	dbUser, err := uc.userRepository.FindUserByIdParam(c.Params("id"))
	if err != nil {
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

	bodyUser, err, validationErrs := uc.userRepository.Validate(c.BodyParser)
	if validationErrs != nil {
		json_res := map[string]interface{}{
			"error":            "Validation Error",
			"validationErrors": validationErrs,
		}
		return c.Status(http.StatusBadRequest).JSON(json_res)
	} else if err != nil {
		json_res := map[string]string{"error": "Error deserializing JSON"}
		return c.Status(http.StatusBadRequest).JSON(json_res)
	}

	if uc.userRepository.ExistsByUsername(bodyUser.Username) == true {
		json_res := map[string]string{"error": "Invalid username"}
		return c.Status(http.StatusBadRequest).JSON(json_res)
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

	dbUser, err := uc.userRepository.FindUserByIdParam(c.Params("id"))
	if err != nil {
		json_res := map[string]string{"error": err.Error()}
		return c.Status(http.StatusBadRequest).JSON(json_res)
	}

	bodyUser, err, validationErrs := uc.userRepository.Validate(c.BodyParser)
	if validationErrs != nil {
		json_res := map[string]interface{}{
			"error":            "Validation Error",
			"validationErrors": validationErrs,
		}
		return c.Status(400).JSON(json_res)

	} else if err != nil {
		json_res := map[string]string{"error": "Error deserializing JSON"}
		return c.Status(400).JSON(json_res)
	}

	dbUser.UpdateFromAnother(bodyUser)

	uc.userRepository.Update(dbUser)

	return c.Status(http.StatusOK).JSON(dbUser)
}

func (uc UserController) Delete(c *fiber.Ctx) error {
	return nil
}
