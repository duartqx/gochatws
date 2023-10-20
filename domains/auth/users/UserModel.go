package users

import (
	e "github.com/duartqx/gochatws/core/errors"
	c "github.com/duartqx/gochatws/core/interfaces"
	i "github.com/duartqx/gochatws/core/interfaces"

	"github.com/go-playground/validator/v10"
)

type UserModel struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username" validate:"email,required"`
	Name     string `db:"name" json:"name" validate:"required,min=3,max=50"`
	Password string `db:"password" json:"password"`
}

func (u UserModel) GetId() int {
	return u.Id
}

func (u UserModel) GetName() string {
	return u.Name
}

func (u UserModel) GetUsername() string {
	return u.Username
}

func (u *UserModel) UpdateFromAnother(other i.User) {
	if other.GetName() != "" {
		u.Name = other.GetName()
	}

	if other.GetUsername() != "" {
		u.Username = other.GetUsername()
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

func (u UserModel) Clean() i.User {
	return &UserClean{Id: u.Id, Name: u.Name, Username: u.Username}
}

// UserClean is returned from UserModel.Clean method to make sure the Password
// field is not leaked even if it's hashed. fiber.Ctx{}.BodyParser was not able
// to parse the password if I set the tag to json:"-" when creating users, so I
// decided for this approach of returning a new struct via a method
type UserClean struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Name     string `db:"name" json:"name"`
}

func (u *UserClean) UpdateFromAnother(other i.User) {}

func (u UserClean) Clean() i.User {
	return &u
}

func (u UserClean) GetId() int {
	return u.Id
}

func (u UserClean) GetName() string {
	return u.Name
}

func (u UserClean) GetUsername() string {
	return u.Username
}
