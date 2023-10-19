package core

import (
	"github.com/go-playground/validator/v10"
	cerr "gochatws/core/errors"
)

type ParserFunc func(out interface{}) error

type Model interface {
	ParseAndValidate(parser ParserFunc, v *validator.Validate) (
		Model, error, *[]cerr.ValidationErrorResponse,
	)
}
