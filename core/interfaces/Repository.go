package core

import (
	cerr "gochatws/core/errors"
)

type Repository interface {
	FindById(id int) (*Model, error)
	Update(m *Model) error
	Create(m *Model) error
	All() (*[]Model, error)
	Validate(p parserFunc) (*Model, error, *[]cerr.ValidationErrorResponse)
}
