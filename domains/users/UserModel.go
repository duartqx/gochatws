package users

import (
	"gochatws/core"

	"github.com/go-playground/validator/v10"
)

type parserFunc func(out interface{}) error

type UserModel struct {
	Id       int    `db:"id" json:"user_id"`
	Username string `db:"username" json:"username" validate:"email,required"`
	Name     string `db:"name" json:"name" validate:"required,min=3,max=50"`
	Password string `db:"password" json:"-"`
}

func (u *UserModel) UpdateFromAnother(other *UserModel) {
	if other.Name != "" {
		u.Name = other.Name
	}

	if other.Username != "" {
		u.Username = other.Username
	}
}

func (u UserModel) ParseAndValidate(parser parserFunc, v *validator.Validate) (
	*UserModel, error, *[]core.ValidationErrorResponse,
) {
	parsedUser := &UserModel{}

	if err := parser(parsedUser); err != nil {
		return nil, err, nil
	}

	if err := v.Struct(parsedUser); err != nil {
		return nil, err, core.BuildErrorResponse(err)
	}

	return parsedUser, nil, nil
}
