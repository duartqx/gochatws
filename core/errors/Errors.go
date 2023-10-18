package core

import "github.com/go-playground/validator/v10"

var (
	// map[string]string Errors
	NotFoundError        = map[string]string{"error": "Not Found"}
	InternalError        = map[string]string{"error": "Internal"}
	InvalidUsernameError = map[string]string{"error": "Invalid username"}
	SerializerError      = map[string]string{"error": "Error deserializing JSON"}

	// func Errors
	ValidationError = func(errs *[]ValidationErrorResponse) map[string]interface{} {
		return map[string]interface{}{
			"error":            "Validation Error",
			"validationErrors": errs,
		}
	}
	CustomMessageError = func(err interface{}) map[string]interface{} {
		return map[string]interface{}{"error": err}
	}
)

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
