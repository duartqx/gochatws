package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	e "github.com/duartqx/gochatws/core/errors"
	s "github.com/duartqx/gochatws/domains/auth/sessions"
	u "github.com/duartqx/gochatws/domains/auth/users"
)

type sessionStore map[string]*s.SessionModel

type ClaimsUser struct {
	Id       int
	Username string
	Name     string
}

type LoginResponse struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Status    string    `json:"status"`
}

type JwtAuthController struct {
	userRepository *u.UserRepository
	secret         *[]byte
	sessionStore   *sessionStore
}

func NewJwtAuthController(ur *u.UserRepository, se *[]byte) *JwtAuthController {
	return &JwtAuthController{
		userRepository: ur,
		secret:         se,
		sessionStore:   &sessionStore{},
	}
}

// private
func (jc JwtAuthController) getUnparsedToken(c *fiber.Ctx) string {
	var (
		token string
		found bool
	)

	token = c.Get("Authorization")

	if token != "" {
		token, found = strings.CutPrefix(token, "Bearer ")
		if found {
			return token
		}
	}
	return c.Cookies("jwt")
}

// private
func (jc JwtAuthController) getParsedToken(c *fiber.Ctx) *jwt.Token {

	unparsedToken := jc.getUnparsedToken(c)
	if unparsedToken == "" {
		return nil
	}

	if _, sessionOk := (*jc.sessionStore)[unparsedToken]; !sessionOk {
		return nil
	}

	parsedToken, err := jwt.Parse(unparsedToken, jc.keyFunc)
	if err != nil || !parsedToken.Valid {

		delete(*jc.sessionStore, unparsedToken)

		return nil
	}

	return parsedToken
}

// private
func (jc JwtAuthController) generateToken(user *ClaimsUser, expiresAt time.Time) (
	string, *fiber.Cookie, error,
) {

	claims := jwt.MapClaims{
		"user": fiber.Map{
			"id":       user.Id,
			"username": user.Username,
			"name":     user.Name,
		},
		"exp": expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(*jc.secret)
	if err != nil {
		return "", &fiber.Cookie{}, fmt.Errorf("Bad secret key")
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenStr,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   true,
	}

	return tokenStr, &cookie, nil
}

// private
func (jc JwtAuthController) keyFunc(t *jwt.Token) (interface{}, error) {
	return *jc.secret, nil
}

// public
func (jc JwtAuthController) AuthNotLoggedMiddleware(c *fiber.Ctx) error {
	if jc.getParsedToken(c) != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.LoggedInError)
	}
	return c.Next()
}

// public
func (jc JwtAuthController) AuthMiddleware(c *fiber.Ctx) error {
	token := jc.getParsedToken(c)
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
			JSON(e.WrongUsernameOrPasswordError)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(dbUser.GetPassword()), []byte(bodyUser.Password),
	); err != nil {
		return c.
			Status(http.StatusUnauthorized).
			JSON(e.WrongUsernameOrPasswordError)
	}

	createdAt := time.Now()
	expiresAt := createdAt.Add(time.Hour * 12)

	tokenStr, cookie, err := jc.generateToken(
		&ClaimsUser{
			Id:       dbUser.GetId(),
			Username: dbUser.GetUsername(),
			Name:     dbUser.GetName(),
		},
		expiresAt,
	)
	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(e.InternalError)
	}

	(*jc.sessionStore)[tokenStr] = &s.SessionModel{
		Token: tokenStr, CreationAt: createdAt, UserId: dbUser.GetId(),
	}

	c.Cookie(cookie)

	return c.
		Status(http.StatusOK).
		JSON(
			LoginResponse{
				Token:     tokenStr,
				CreatedAt: createdAt,
				ExpiresAt: expiresAt,
				Status:    "Logged In",
			},
		)
}

// public
func (jc JwtAuthController) Logout(c *fiber.Ctx) error {

	sessionToken := jc.getUnparsedToken(c)

	delete(*jc.sessionStore, sessionToken)

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{"status": "Logged out"})
}
