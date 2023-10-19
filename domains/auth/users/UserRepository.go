package auth

import (
	"encoding/gob"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"

	e "gochatws/core/errors"
	c "gochatws/core/interfaces"
)

type UserRepository struct {
	db *sqlx.DB
	v  *validator.Validate
}

func NewUserRepo(db *sqlx.DB, v *validator.Validate) *UserRepository {

	gob.Register(UserModel{})

	return &UserRepository{db: db, v: v}
}

func (ur UserRepository) getModel() *UserModel {
	return &UserModel{}
}

func (ur UserRepository) FindById(id int) (*UserModel, error) {
	user := ur.getModel()
	err := ur.db.Get(user, "SELECT * FROM User WHERE ID = $1", id)
	return user, err
}

func (ur UserRepository) FindByIdParam(id string) (*UserModel, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	foundUser, err := ur.FindById(idInt)
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

func (ur UserRepository) FindByUsername(username string) (*UserModel, error) {
	user := ur.getModel()
	err := ur.db.Get(user, "SELECT * FROM User WHERE username = $1", username)
	return user, err
}

func (ur UserRepository) ParseAndValidate(parser c.ParserFunc) (
	*UserModel, error, *[]e.ValidationErrorResponse,
) {
	return ur.getModel().ParseAndValidate(parser, ur.v)
}

func (ur UserRepository) Parse(parser c.ParserFunc) (*UserModel, error) {
	parsedUser := ur.getModel()

	if err := parser(parsedUser); err != nil {
		return nil, err
	}

	return parsedUser, nil
}

func (ur UserRepository) ExistsByUsername(username string) bool {
	var exists int
	_ = ur.db.Get(
		&exists,
		"SELECT EXISTS(SELECT 1 FROM User WHERE username = $1)",
		username,
	)
	if exists > 0 {
		return true
	}
	return false
}

func (ur UserRepository) All() (*[]UserModel, error) {
	users := []UserModel{}
	err := ur.db.Select(&users, "SELECT * FROM User")
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (ur UserRepository) Update(u *UserModel) error {
	_, err := ur.db.Exec(
		"UPDATE User SET name = $1, username = $2 WHERE id = $3",
		u.Name,
		u.Username,
		u.Id,
	)
	return err
}

func (ur UserRepository) Create(u *UserModel) error {
	_, err := ur.db.Exec(
		"INSERT INTO User (name, username, password) VALUES ($1, $2, $3)",
		u.Name,
		u.Username,
		u.Password,
	)
	if err != nil {
		return err
	}
	return ur.db.QueryRow("SELECT last_insert_rowid()").Scan(&u.Id)
}
