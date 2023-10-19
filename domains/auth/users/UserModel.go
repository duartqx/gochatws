package auth

import (
	e "github.com/duartqx/gochatws/core/errors"
	c "github.com/duartqx/gochatws/core/interfaces"

	"github.com/go-playground/validator/v10"
)

type UserClean struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Name     string `db:"name" json:"name"`
}

type UserModel struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username" validate:"email,required"`
	Name     string `db:"name" json:"name" validate:"required,min=3,max=50"`
	Password string `db:"password" json:"password"`
}

func (u *UserModel) UpdateFromAnother(other *UserModel) {
	if other.Name != "" {
		u.Name = other.Name
	}

	if other.Username != "" {
		u.Username = other.Username
	}
}

func (u UserModel) ParseAndValidate(parser c.ParserFunc, v *validator.Validate) (
	*UserModel, error, *[]e.ValidationErrorResponse,
) {
	parsedUser := &UserModel{}

	if err := parser(parsedUser); err != nil {
		return nil, err, nil
	}

	if err := v.Struct(parsedUser); err != nil {
		return nil, err, e.BuildErrorResponse(err)
	}

	return parsedUser, nil, nil
}

func (u UserModel) Clean() *UserClean {
	return &UserClean{Id: u.Id, Name: u.Name, Username: u.Username}
}
