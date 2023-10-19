package core

import "github.com/go-playground/validator/v10"

var (
	// map[string]string Errors
	InternalError           = map[string]string{"error": "Internal"}
	InvalidUsernameError    = map[string]string{"error": "Invalid username"}
	InvalidTokenError       = map[string]string{"error": "Invalid or missing token"}
	WrongUsernameOrPassword = map[string]string{"error": "Wrong username or password"}
	NotFoundError           = map[string]string{"error": "Not Found"}
	PasswordTooLongError    = map[string]string{"error": "Unfortunately your password is too long"}
	SerializerError         = map[string]string{"error": "Error deserializing JSON"}

	// func Errors
	CustomMessageError = func(err interface{}) map[string]interface{} {
		return map[string]interface{}{"error": err}
	}
	ValidationError = func(errs *[]ValidationErrorResponse) map[string]interface{} {
		return map[string]interface{}{
			"error":            "Validation Error",
			"validationErrors": errs,
		}
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
