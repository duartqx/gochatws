package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	e "github.com/duartqx/gochatws/application/errors"
	s "github.com/duartqx/gochatws/application/services"
	u "github.com/duartqx/gochatws/domains/entities/user"
)

type JwtAuthController struct {
	jwtAuthService *s.JwtAuthService
}

func NewJwtAuthController(jwtAuthService *s.JwtAuthService) *JwtAuthController {
	return &JwtAuthController{
		jwtAuthService: jwtAuthService,
	}
}

// private
func (jc JwtAuthController) getModel() *u.UserModel {
	return &u.UserModel{}
}

// private
func (jc JwtAuthController) authNotLoggedMiddleware(redirect bool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		claimsUser, _ := jc.jwtAuthService.ValidateAuth(
			c.Get("Authorization"),
			c.Cookies("jwt"),
		)
		if claimsUser != nil {
			if redirect {
				return c.Redirect("/")
			}
			return c.Status(http.StatusUnauthorized).JSON(e.LoggedInError)
		}
		return c.Next()
	}
}

// public
func (jc JwtAuthController) AuthNotLoggedMiddleware() func(c *fiber.Ctx) error {
	return jc.authNotLoggedMiddleware(false)
}

// public
func (jc JwtAuthController) AuthNotLoggedMiddlewareWithRedirect() func(c *fiber.Ctx) error {
	return jc.authNotLoggedMiddleware(true)
}

// private
func (jc JwtAuthController) authMiddleware(redirect bool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		claimsUser, err := jc.jwtAuthService.ValidateAuth(
			c.Get("Authorization"),
			c.Cookies("jwt"),
		)
		if err != nil {
			if redirect {
				return c.Redirect("/login")
			}
			return c.Status(http.StatusUnauthorized).JSON(e.InvalidTokenError)
		}

		c.Locals("user", claimsUser)

		return c.Next()
	}
}

// public
func (jc JwtAuthController) AuthMiddleware() func(c *fiber.Ctx) error {
	return jc.authMiddleware(false)
}

// public
func (jc JwtAuthController) AuthMiddlewareWithRedirect() func(c *fiber.Ctx) error {
	return jc.authMiddleware(true)
}

// public
func (jc JwtAuthController) Login(c *fiber.Ctx) error {

	user := jc.getModel()

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.BadRequestError)
	}

	response := jc.jwtAuthService.Login(user)

	if response.Cookie != nil {
		c.Cookie(&fiber.Cookie{
			Name:     response.Cookie.Name,
			Value:    response.Cookie.Value,
			Expires:  response.Cookie.Expires,
			Secure:   response.Cookie.Secure,
			HTTPOnly: response.Cookie.HTTPOnly,
		})
	}

	return c.Status(response.Status).JSON(response.Body)
}

// public
func (jc JwtAuthController) Logout(c *fiber.Ctx) error {

	response := jc.jwtAuthService.Logout(
		c.Get("Authorization"),
		c.Cookies("jwt"),
	)

	return c.Status(response.Status).JSON(response.Body)
}
