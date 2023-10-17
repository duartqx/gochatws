package core

import "github.com/go-playground/validator/v10"

type ValidationErrorResponse struct {
	Error bool
	Field string
	Tag   string
	Value interface{}
}

func BuildErrorResponse(errs error) *[]ValidationErrorResponse {
	validationErrors := []ValidationErrorResponse{}
	for _, err := range errs.(validator.ValidationErrors) {
		elem := ValidationErrorResponse{
			Tag:   err.Tag(), // Export struct tag
			Field: err.Field(),
			Value: err.Value(), // Export field value
			Error: true,
		}
		validationErrors = append(validationErrors, elem)
	}
	return &validationErrors
}
