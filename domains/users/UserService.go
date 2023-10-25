package users

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	e "github.com/duartqx/gochatws/core/errors"
	h "github.com/duartqx/gochatws/core/http"
	i "github.com/duartqx/gochatws/core/interfaces"
)

type UserService struct {
	userRepository *UserRepository
	validator      *validator.Validate
}

func NewUserService(userRespository *UserRepository, v *validator.Validate) *UserService {
	return &UserService{
		userRepository: userRespository,
		validator:      v,
	}
}

func (us UserService) All() (*h.HttpResponse, error) {
	users, err := us.userRepository.All()
	if err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
		return resp, fmt.Errorf("Could not read All users")
	}
	resp := &h.HttpResponse{Status: http.StatusOK, Body: users}
	return resp, nil
}

func (us UserService) Get(userId int) (*h.HttpResponse, error) {
	dbUser, err := us.userRepository.FindById(userId)
	if err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusNotFound,
			Body:   e.NotFoundError,
		}
		return resp, fmt.Errorf("User not found")
	}
	resp := &h.HttpResponse{Status: http.StatusNotFound, Body: dbUser.Clean()}
	return resp, nil
}

func (us UserService) Create(user i.User) (*h.HttpResponse, error) {

	if err := us.validator.Struct(user); err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusBadRequest,
			Body:   e.ValidationError(e.BuildErrorResponse(err)),
		}
		return resp, fmt.Errorf("Validation Error")
	}

	if us.userRepository.ExistsByUsername(user.GetUsername()) {
		resp := &h.HttpResponse{
			Status: http.StatusBadRequest,
			Body:   e.InvalidUsernameError,
		}
		return resp, fmt.Errorf("Username not unique")
	}

	hashedPassword, err :=
		bcrypt.GenerateFromPassword([]byte(user.GetPassword()), 10)
	if err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusBadRequest,
			Body:   e.PasswordTooLongError,
		}
		return resp, fmt.Errorf("Password too long error")
	}
	user.SetPassword(string(hashedPassword))

	us.userRepository.Create(user)

	resp := &h.HttpResponse{Status: http.StatusCreated, Body: user.Clean()}
	return resp, nil
}

func (us UserService) Update(bodyUser i.User) (*h.HttpResponse, error) {

	if err := us.validator.Struct(bodyUser); err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusBadRequest,
			Body:   e.ValidationError(e.BuildErrorResponse(err)),
		}
		return resp, fmt.Errorf("Validation Error")
	}

	dbUser, err := us.userRepository.FindById(bodyUser.GetId())
	if err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusNotFound,
			Body:   e.NotFoundError,
		}
		return resp, fmt.Errorf("User not found")
	}

	dbUser.UpdateFromAnother(bodyUser)
	resp := &h.HttpResponse{Status: http.StatusOK, Body: dbUser.Clean()}
	return resp, nil
}

func (us UserService) Delete(user i.User) (*h.HttpResponse, error) {
	err := us.userRepository.Delete(user)
	if err != nil {
		resp := &h.HttpResponse{
			Status: http.StatusInternalServerError,
			Body:   e.InternalError,
		}
		return resp, fmt.Errorf("Could not delete user")
	}
	msg := fmt.Sprintf("Successfully deleted user with id: %d", user.GetId())
	resp := &h.HttpResponse{
		Status: http.StatusOK,
		Body:   map[string]string{"user": msg},
	}
	return resp, nil
}
