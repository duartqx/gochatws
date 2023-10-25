package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	e "github.com/duartqx/gochatws/core/errors"
	h "github.com/duartqx/gochatws/core/http"
	i "github.com/duartqx/gochatws/core/interfaces"
	s "github.com/duartqx/gochatws/core/sessions"
)

type ClaimsUser struct {
	Id       int
	Username string
	Name     string
}

type JwtAuthService struct {
	userRepository i.UserRepository
	secret         *[]byte
	sessionStore   i.SessionStore
}

func NewJwtAuthService(
	userRepository i.UserRepository,
	secret *[]byte,
	sessionStore i.SessionStore,
) *JwtAuthService {
	return &JwtAuthService{
		userRepository: userRepository,
		secret:         secret,
		sessionStore:   sessionStore,
	}
}

// private
func (jas JwtAuthService) keyFunc(t *jwt.Token) (interface{}, error) {
	return *jas.secret, nil
}

// private
func (jas JwtAuthService) generateToken(user *ClaimsUser, expiresAt time.Time) (
	string, *h.Cookie, error,
) {

	claims := jwt.MapClaims{
		"user": map[string]interface{}{
			"id":       user.Id,
			"username": user.Username,
			"name":     user.Name,
		},
		"exp": expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(*jas.secret)
	if err != nil {
		return "", &h.Cookie{}, fmt.Errorf("Bad secret key")
	}

	cookie := &h.Cookie{
		Name:     "jwt",
		Value:    tokenStr,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   true,
	}

	return tokenStr, cookie, nil
}

func (jas JwtAuthService) GetUnparsedToken(authorization, cookie string) string {
	if authorization != "" {
		token, found := strings.CutPrefix(authorization, "Bearer ")
		if found {
			return token
		}
	}
	return cookie
}

func (jas JwtAuthService) GetParsedToken(authorization, cookie string) *jwt.Token {

	unparsedToken := jas.GetUnparsedToken(authorization, cookie)
	if unparsedToken == "" {
		return nil
	}

	if _, err := jas.sessionStore.Get(unparsedToken); err != nil {
		return nil
	}

	parsedToken, err := jwt.Parse(unparsedToken, jas.keyFunc)
	if err != nil || !parsedToken.Valid {

		jas.sessionStore.Delete(unparsedToken)

		return nil
	}

	return parsedToken
}

func (jas JwtAuthService) Login(user i.User) (*h.HttpResponse, error) {

	if user.GetUsername() == "" || user.GetPassword() == "" {
		resp := &h.HttpResponse{
			Status: http.StatusBadRequest,
			Body:   e.BadRequestError,
		}
		return resp, fmt.Errorf("Missing required fields!")
	}

	dbUser, err := jas.userRepository.FindByUsername(user.GetUsername())
	if err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusUnauthorized,
			Body:   e.WrongUsernameOrPasswordError,
		}
		return resp, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(dbUser.GetPassword()), []byte(user.GetPassword()),
	); err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusUnauthorized,
			Body:   e.WrongUsernameOrPasswordError,
		}
		return resp, err
	}

	createdAt := time.Now()
	expiresAt := createdAt.Add(time.Hour * 12)

	tokenStr, cookie, err := jas.generateToken(
		&ClaimsUser{
			Id:       dbUser.GetId(),
			Username: dbUser.GetUsername(),
			Name:     dbUser.GetName(),
		},
		expiresAt,
	)
	if err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
		return resp, err
	}

	if err := jas.sessionStore.Set(tokenStr, &s.SessionModel{
		Token: tokenStr, CreationAt: createdAt, UserId: dbUser.GetId(),
	}); err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
		return resp, err
	}

	resp := &h.HttpResponse{
		Status: http.StatusOK,
		Body: h.LoginResponse{
			Token:     tokenStr,
			CreatedAt: createdAt,
			ExpiresAt: expiresAt,
			Status:    "Valid",
		},
		Cookie: cookie,
	}
	return resp, nil
}

func (jas *JwtAuthService) DeleteFromStore(token string) error {
	return jas.sessionStore.Delete(token)
}
