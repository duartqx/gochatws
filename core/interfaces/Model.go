package core

import (
	e "github.com/duartqx/gochatws/core/errors"
	"github.com/go-playground/validator/v10"
)

type ParserFunc func(out interface{}) error

type Model interface {
	ParseAndValidate(parser ParserFunc, v *validator.Validate) (
		Model, error, *[]e.ValidationErrorResponse,
	)
}
