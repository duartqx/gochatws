package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"

	e "github.com/duartqx/gochatws/core/errors"
	i "github.com/duartqx/gochatws/core/interfaces"
)

type BasicAuthController struct {
	userRepository i.UserRepository
	sessionStore   *session.Store
}

func NewBasicAuthController(ur i.UserRepository, st *session.Store) *BasicAuthController {
	return &BasicAuthController{
		userRepository: ur,
		sessionStore:   st,
	}
}

func (lm BasicAuthController) basicAuth(auth string) (string, string, error) {

	splitAuth := strings.Split(auth, " ")

	if len(splitAuth) != 2 {
		return "", "", fmt.Errorf("%s", "Bad Authorization Header")
	}

	authType := splitAuth[0]
	authBase64Str := splitAuth[1]

	if authType != "Basic" {
		return "", "", fmt.Errorf("%s", "Bad Authorization Type")
	}

	authStr, err := base64.StdEncoding.DecodeString(authBase64Str)
	if err != nil {
		return "", "", err
	}

	authParts := strings.Split(string(authStr), ":")
	if len(authParts) != 2 {
		return "", "", fmt.Errorf("%s", "Bad Auth Parts")
	}

	username := authParts[0]
	password := authParts[1]

	return username, password, nil
}

func (lm BasicAuthController) Login(c *fiber.Ctx) error {

	username, password, err := lm.basicAuth(c.Get("Authorization"))
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	user, err := lm.userRepository.FindByUsername(username)
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.WrongUsernameOrPasswordError)
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.GetPassword()), []byte(password),
	); err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.WrongUsernameOrPasswordError)
	}

	se, err := lm.sessionStore.Get(c)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	se.Set("user", user)
	se.SetExpiry(time.Hour * 12)
	se.Save()
	return c.Status(http.StatusOK).JSON(map[string]string{"message": "Logged in"})
}

func (lm BasicAuthController) AuthMiddleware(c *fiber.Ctx) error {
	se, err := lm.sessionStore.Get(c)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	userInterface := se.Get("user")
	if userInterface == nil {
		return c.SendStatus(http.StatusUnauthorized)
	}
	user, ok := userInterface.(i.User)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}

	// Gets username and password from Authorization: Basic header
	username, password, err := lm.basicAuth(c.Get("Authorization"))
	if err != nil || username != user.GetUsername() {
		return c.SendStatus(http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.GetPassword()), []byte(password),
	); err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	return c.Next()
}
