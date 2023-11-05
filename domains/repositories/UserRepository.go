package repositories

import (
	u "github.com/duartqx/gochatws/domains/entities/user"
)

type UserRepository interface {
	All() (*[]u.User, error)
	Create(user u.User) error
	Delete(user u.User) error
	ExistsById(id int) bool
	ExistsByUsername(username string) bool
	FindById(id int) (u.User, error)
	FindByIdParam(id string) (u.User, error)
	FindByUsername(username string) (u.User, error)
	Update(user u.User) error
}
