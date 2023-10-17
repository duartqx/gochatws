package core

import "github.com/go-playground/validator/v10"

type parserFunc func(out interface{}) error

type Model interface {
	ParseAndValidate(parser parserFunc, v *validator.Validate) (
		*Model, error, *[]ValidationErrorResponse,
	)
}
