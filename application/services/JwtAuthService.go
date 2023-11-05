package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	e "github.com/duartqx/gochatws/application/errors"
	h "github.com/duartqx/gochatws/application/http"
	u "github.com/duartqx/gochatws/domains/entities/user"
	r "github.com/duartqx/gochatws/domains/repositories"
	s "github.com/duartqx/gochatws/infrastructure/sessions"
)

type ClaimsUser struct {
	Id       int
	Username string
	Name     string
}

type JwtAuthService struct {
	userRepository r.UserRepository
	secret         *[]byte
	sessionStore   *s.SessionStore
}

func NewJwtAuthService(
	userRepository r.UserRepository,
	secret *[]byte,
	sessionStore *s.SessionStore,
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

// private
func (jas JwtAuthService) getUnparsedToken(authorization, cookie string) string {
	if authorization != "" {
		token, found := strings.CutPrefix(authorization, "Bearer ")
		if found {
			return token
		}
	}
	return cookie
}

func (jas JwtAuthService) ValidateAuth(authorization, cookie string) (interface{}, error) {

	unparsedToken := jas.getUnparsedToken(authorization, cookie)
	if unparsedToken == "" {
		return nil, fmt.Errorf("Missing Token")
	}

	if _, err := jas.sessionStore.Get(unparsedToken); err != nil {
		return nil, fmt.Errorf("Missing session")
	}

	parsedToken, err := jwt.Parse(unparsedToken, jas.keyFunc)
	if err != nil || !parsedToken.Valid {

		jas.sessionStore.Delete(unparsedToken)

		return nil, fmt.Errorf("Expired session")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Could not parse claims")
	}

	return claims["user"], nil
}

func (jas JwtAuthService) Login(user u.User) *h.HttpResponse {

	invalidCookie := &h.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(time.Hour * -3),
		HTTPOnly: true,
		Secure:   true,
	}

	if user.GetUsername() == "" || user.GetPassword() == "" {
		return &h.HttpResponse{
			Status: http.StatusBadRequest,
			Body:   e.BadRequestError,
			Cookie: invalidCookie,
		}
	}

	dbUser, err := jas.userRepository.FindByUsername(user.GetUsername())
	if err != nil {
		return &h.HttpResponse{
			Status: http.StatusUnauthorized,
			Body:   e.WrongUsernameOrPasswordError,
			Cookie: invalidCookie,
		}
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(dbUser.GetPassword()), []byte(user.GetPassword()),
	); err != nil {
		return &h.HttpResponse{
			Status: http.StatusUnauthorized,
			Body:   e.WrongUsernameOrPasswordError,
			Cookie: invalidCookie,
		}
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
		return &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
			Cookie: invalidCookie,
		}
	}

	if err := jas.sessionStore.Set(tokenStr, &s.SessionModel{
		Token: tokenStr, CreationAt: createdAt, UserId: dbUser.GetId(),
	}); err != nil {
		return &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
			Cookie: invalidCookie,
		}
	}

	return &h.HttpResponse{
		Status: http.StatusOK,
		Body: h.LoginResponse{
			Token:     tokenStr,
			CreatedAt: createdAt,
			ExpiresAt: expiresAt,
			Status:    "Valid",
		},
		Cookie: cookie,
	}
}

func (jas *JwtAuthService) Logout(authorization, cookie string) *h.HttpResponse {
	unparsedToken := jas.getUnparsedToken(authorization, cookie)
	if unparsedToken == "" {
		return &h.HttpResponse{
			Status: http.StatusUnauthorized,
			Body:   e.UnauthorizedError,
		}
	}

	if err := jas.sessionStore.Delete(unparsedToken); err != nil {
		return &h.HttpResponse{
			Status: http.StatusUnauthorized,
			Body:   e.UnauthorizedError,
		}
	}

	return &h.HttpResponse{
		Status: http.StatusOK,
		Body:   map[string]string{"status": "Logged out"},
	}
}
