package users

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db       *sqlx.DB
	validate *validator.Validate
}

func NewUserRepo(db *sqlx.DB, v *validator.Validate) *UserRepository {
	return &UserRepository{db: db, validate: v}
}

func (ur UserRepository) getModel() *UserModel {
	return &UserModel{}
}

func (ur UserRepository) All() (*[]UserModel, error) {
	var users []UserModel
	err := ur.db.Select(&users, "SELECT * FROM User")
	return &users, err
}

func (ur UserRepository) FindById(id int) (*UserModel, error) {
	user := ur.getModel()
	err := ur.db.Get(user, "SELECT * FROM User WHERE ID = $1", id)
	return user, err
}

func (ur UserRepository) ExistsByUsername(username string) bool {
	var exists int
	_ = ur.db.Get(
		&exists,
		"SELECT COUNT(*) FROM User WHERE username = $1",
		username,
	)
	if exists > 0 {
		return true
	}
	return false
}

func (ur UserRepository) Validate(u *UserModel) error {
	return ur.validate.Struct(*u)
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
		u.Id,
	)
	if err != nil {
		return err
	}
	return ur.db.QueryRow("SELECT last_insert_rowid()").Scan(&u.Id)
}
