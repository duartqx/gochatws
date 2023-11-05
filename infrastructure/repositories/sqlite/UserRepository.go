package repositories

import (
	"strconv"

	"github.com/jmoiron/sqlx"

	u "github.com/duartqx/gochatws/domains/entities/user"
	r "github.com/duartqx/gochatws/domains/repositories"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) r.UserRepository {
	return &UserRepository{db: db}
}

func (ur UserRepository) All() (*[]u.User, error) {
	users := []u.User{}
	rows, err := ur.db.Query("SELECT id, name, username FROM User")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := ur.GetCleanModel()

		var (
			id       int
			name     string
			username string
		)

		if err := rows.Scan(&id, &name, &username); err != nil {
			return nil, err
		}

		user.SetId(id).SetName(name).SetUsername(username)

		var iUser u.User = user

		users = append(users, iUser)
	}
	return &users, nil
}

func (ur UserRepository) Create(u u.User) error {
	result, err := ur.db.Exec(
		"INSERT INTO User (name, username, password) VALUES ($1, $2, $3)",
		u.GetName(),
		u.GetUsername(),
		u.GetPassword(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.SetId(int(id))

	return nil
}

func (ur UserRepository) Delete(u u.User) error {
	_, err := ur.db.Exec(
		"DELETE FROM User WHERE id = $1", u.GetId(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (ur UserRepository) Update(u u.User) error {
	_, err := ur.db.Exec(
		"UPDATE User SET name = $1, username = $2 WHERE id = $3",
		u.GetName(),
		u.GetUsername(),
		u.GetId(),
	)
	return err
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

func (ur UserRepository) ExistsById(id int) bool {
	var exists int
	_ = ur.db.Get(
		&exists,
		"SELECT EXISTS(SELECT 1 FROM User WHERE id = $1)",
		id,
	)
	if exists > 0 {
		return true
	}
	return false
}

func (ur UserRepository) FindById(id int) (u.User, error) {
	user := ur.GetModel()
	err := ur.db.Get(user, "SELECT * FROM User WHERE ID = $1 LIMIT 1", id)
	return user, err
}

func (ur UserRepository) FindByUsername(username string) (u.User, error) {
	user := ur.GetModel()
	err := ur.db.Get(user, "SELECT * FROM User WHERE username = $1", username)
	return user, err
}

func (ur UserRepository) GetModel() *u.UserModel {
	return &u.UserModel{}
}

func (ur UserRepository) GetCleanModel() *u.UserDTO {
	return &u.UserDTO{}
}

func (ur UserRepository) FindByIdParam(id string) (u.User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return ur.FindById(idInt)
}
