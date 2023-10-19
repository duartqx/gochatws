package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	e "gochatws/core/errors"
	u "gochatws/domains/auth/users"
)

type JwtAuthController struct {
	userRepository *u.UserRepository
	secret         *[]byte
}

func NewJwtAuthController(ur *u.UserRepository, se *[]byte) *JwtAuthController {
	return &JwtAuthController{
		userRepository: ur,
		secret:         se,
	}
}

func (jc JwtAuthController) getToken(c *fiber.Ctx) string {
	var (
		token string
		found bool
	)
	token = c.Get("Authorization")
	if token != "" {
		token, found = strings.CutPrefix(token, "Bearer ")
		if !found {
			return ""
		}
		return token
	}
	return c.Cookies("jwt")
}

func (jc JwtAuthController) keyFunc(t *jwt.Token) (interface{}, error) {
	return *jc.secret, nil
}

func (jc JwtAuthController) AuthenticationMiddleware(c *fiber.Ctx) error {
	unparsedToken := jc.getToken(c)

	if unparsedToken == "" {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.InvalidTokenError)
	}

	parsedToken, err := jwt.Parse(unparsedToken, jc.keyFunc)
	if err != nil || !parsedToken.Valid {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.InvalidTokenError)
	}

	return c.Next()
}

func (jc JwtAuthController) Login(c *fiber.Ctx) error {
	bodyUser, err := jc.userRepository.Parse(c.BodyParser)

	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(e.SerializerError)
	}

	dbUser, err := jc.userRepository.FindByUsername(bodyUser.Username)
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.WrongUsernameOrPassword)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(dbUser.Password), []byte(bodyUser.Password),
	); err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.WrongUsernameOrPassword)
	}

	claims := jwt.MapClaims{
		"user": fiber.Map{
			"id":       dbUser.Id,
			"username": dbUser.Username,
			"name":     dbUser.Name,
		},
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(*jc.secret)
	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(e.InternalError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tokenStr,
		HTTPOnly: true,
	})

	return c.
		Status(http.StatusOK).
		JSON(map[string]string{"token": tokenStr})
}
