package services

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	e "github.com/duartqx/gochatws/application/errors"
	h "github.com/duartqx/gochatws/application/http"

	u "github.com/duartqx/gochatws/domains/entities/user"
	r "github.com/duartqx/gochatws/domains/repositories"
)

type UserService struct {
	userRepository r.UserRepository
	validator      *validator.Validate
}

func NewUserService(userRespository r.UserRepository, v *validator.Validate) *UserService {
	return &UserService{
		userRepository: userRespository,
		validator:      v,
	}
}

func (us UserService) All() *h.HttpResponse {
	users, err := us.userRepository.All()
	if err != nil {
		return &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
	}
	return &h.HttpResponse{Status: http.StatusOK, Body: users}
}

func (us UserService) Get(userId int) *h.HttpResponse {
	dbUser, err := us.userRepository.FindById(userId)
	if err != nil {
		return &h.HttpResponse{Status: http.StatusNotFound, Body: e.NotFoundError}
	}
	return &h.HttpResponse{Status: http.StatusOK, Body: dbUser.Clean()}
}

func (us UserService) Create(user *u.UserModel) *h.HttpResponse {

	if err := us.validator.Struct(user); err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusBadRequest,
			Body:   e.ValidationError(e.BuildErrorResponse(err)),
		}
		return resp
	}

	if us.userRepository.ExistsByUsername(user.GetUsername()) {
		return &h.HttpResponse{
			Status: http.StatusBadRequest,
			Body:   e.InvalidUsernameError,
		}
	}

	hashedPassword, err :=
		bcrypt.GenerateFromPassword([]byte(user.GetPassword()), 10)
	if err != nil {
		return &h.HttpResponse{
			Status: http.StatusBadRequest,
			Body:   e.PasswordTooLongError,
		}
	}
	user.SetPassword(string(hashedPassword))

	if err := us.userRepository.Create(user); err != nil {
		return &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
	}

	return &h.HttpResponse{Status: http.StatusCreated, Body: user.Clean()}
}

func (us UserService) Update(bodyUser u.User) *h.HttpResponse {

	dbUser, err := us.userRepository.FindById(bodyUser.GetId())
	if err != nil {
		return &h.HttpResponse{Status: http.StatusNotFound, Body: e.NotFoundError}
	}

	dbUser.UpdateFromAnother(bodyUser)
	if err := us.userRepository.Update(dbUser); err != nil {
		return &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
	}
	return &h.HttpResponse{Status: http.StatusOK, Body: dbUser.Clean()}
}

func (us UserService) Delete(user u.User) *h.HttpResponse {
	err := us.userRepository.Delete(user)
	if err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
		return resp
	}
	msg := fmt.Sprintf("Successfully deleted user with id: %d", user.GetId())
	resp := &h.HttpResponse{
		Status: http.StatusOK,
		Body:   map[string]string{"user": msg},
	}
	return resp
}
