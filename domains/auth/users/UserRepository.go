package users

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"

	e "github.com/duartqx/gochatws/core/errors"
	i "github.com/duartqx/gochatws/core/interfaces"
)

type UserRepository struct {
	db *sqlx.DB
	v  *validator.Validate
}

func NewUserRepository(db *sqlx.DB, v *validator.Validate) *UserRepository {
	return &UserRepository{db: db, v: v}
}

func (ur UserRepository) GetModel() *UserModel {
	return &UserModel{}
}

func (ur UserRepository) GetCleanModel() *UserClean {
	return &UserClean{}
}

func (ur UserRepository) FindById(id int) (i.User, error) {
	user := ur.GetModel()
	err := ur.db.Get(user, "SELECT * FROM User WHERE ID = $1 LIMIT 1", id)
	return user, err
}

func (ur UserRepository) FindByIdParam(id string) (i.User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return ur.FindById(idInt)
}

func (ur UserRepository) FindByUsername(username string) (i.User, error) {
	user := ur.GetModel()
	err := ur.db.Get(user, "SELECT * FROM User WHERE username = $1", username)
	return user, err
}

func (ur UserRepository) ParseAndValidate(parser i.ParserFunc) (
	*UserModel, error, *[]e.ValidationErrorResponse,
) {
	return ur.GetModel().ParseAndValidate(parser, ur.v)
}

func (ur UserRepository) Parse(parser i.ParserFunc) (*UserModel, error) {
	parsedUser := ur.GetModel()

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

func (ur UserRepository) All() (*[]i.User, error) {
	users := []i.User{}
	rows, err := ur.db.Query("SELECT id, name, username FROM User")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := ur.GetCleanModel()

		if err := rows.Scan(&user.Id, &user.Name, &user.Username); err != nil {
			return nil, err
		}

		var iUser i.User = user

		users = append(users, iUser)
	}
	return &users, nil
}

func (ur UserRepository) Update(u i.User) error {
	_, err := ur.db.Exec(
		"UPDATE User SET name = $1, username = $2 WHERE id = $3",
		u.GetName(),
		u.GetUsername(),
		u.GetId(),
	)
	return err
}

func (ur UserRepository) Create(u i.User) error {
	_, err := ur.db.Exec(
		"INSERT INTO User (name, username, password) VALUES ($1, $2, $3)",
		u.GetName(),
		u.GetUsername(),
		u.GetPassword(),
	)
	if err != nil {
		return err
	}

	var id int

	err = ur.db.QueryRow("SELECT last_insert_rowid()").Scan(&id)
	if err != nil {
		return err
	}

	u.SetId(id)

	return nil
}

func (ur UserRepository) Delete(u i.User) error {
	_, err := ur.db.Exec(
		"DELETE FROM User WHERE id = $1", u.GetId(),
	)
	if err != nil {
		return err
	}
	return nil
}
