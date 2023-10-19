package core

import (
	cerr "gochatws/core/errors"
)

type M interface{}

type Repository interface {
	All() (*[]M, error)
	Create(m *M) error
	FindById(id int) (*M, error)
	Update(m M) error
	Validate(p ParserFunc) (*M, error, *[]cerr.ValidationErrorResponse)
}
