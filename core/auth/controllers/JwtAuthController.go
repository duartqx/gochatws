package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	as "github.com/duartqx/gochatws/core/auth/service"
	e "github.com/duartqx/gochatws/core/errors"
)

type JwtAuthController struct {
	jwtAuthService *as.JwtAuthService
}

func NewJwtAuthController(jwtAuthService *as.JwtAuthService) *JwtAuthController {
	return &JwtAuthController{
		jwtAuthService: jwtAuthService,
	}
}

// public
func (jc JwtAuthController) AuthNotLoggedMiddleware(c *fiber.Ctx) error {
	var token *jwt.Token = jc.jwtAuthService.GetParsedToken(
		c.Get("Authorization"),
		c.Cookies("jwt"),
	)
	if token != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.LoggedInError)
	}
	return c.Next()
}

// public
func (jc JwtAuthController) AuthMiddleware(c *fiber.Ctx) error {
	var token *jwt.Token = jc.jwtAuthService.GetParsedToken(
		c.Get("Authorization"),
		c.Cookies("jwt"),
	)
	if token == nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.InvalidTokenError)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.InvalidTokenError)
	}

	c.Locals("user", claims["user"])

	return c.Next()
}

// public
func (jc JwtAuthController) Login(c *fiber.Ctx) error {

	response, err := jc.jwtAuthService.Login(c.BodyParser)

	if err != nil {
		return c.Status(response.Status).JSON(response.Body)
	}

	c.Cookie(&fiber.Cookie{
		Name:     response.Cookie.Name,
		Value:    response.Cookie.Value,
		Expires:  response.Cookie.Expires,
		Secure:   response.Cookie.Secure,
		HTTPOnly: response.Cookie.HTTPOnly,
	})

	return c.Status(response.Status).JSON(response.Body)
}

// public
func (jc JwtAuthController) Logout(c *fiber.Ctx) error {

	sessionToken := jc.jwtAuthService.GetUnparsedToken(
		c.Get("Authorization"),
		c.Cookies("jwt"),
	)

	if err := jc.jwtAuthService.DeleteFromStore(sessionToken); err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.UnauthorizedError)
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{"status": "Logged out"})
}
