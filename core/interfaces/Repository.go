package core

import (
	e "github.com/duartqx/gochatws/core/errors"
)

type M interface{}

type Repository interface {
	All() (*[]M, error)
	Create(m *M) error
	FindById(id int) (*M, error)
	Update(m M) error
	Validate(p ParserFunc) (*M, error, *[]e.ValidationErrorResponse)
}
