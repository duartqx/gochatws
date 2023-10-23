package errors

import "github.com/go-playground/validator/v10"

type ErrorResponse map[string]string

var (
	// map[string]string Errors
	BadRequestError              = ErrorResponse{"error": "Bad Request"}
	InternalError                = ErrorResponse{"error": "Internal"}
	InvalidUsernameError         = ErrorResponse{"error": "Invalid username"}
	InvalidTokenError            = ErrorResponse{"error": "Invalid or missing token"}
	WrongUsernameOrPasswordError = ErrorResponse{"error": "Wrong username or password"}
	LoggedInError                = ErrorResponse{"error": "You are logged in"}
	NotFoundError                = ErrorResponse{"error": "Not Found"}
	PasswordTooLongError         = ErrorResponse{"error": "Unfortunately your password is too long"}
	SerializerError              = ErrorResponse{"error": "Error deserializing JSON"}
	UnauthorizedError            = ErrorResponse{"error": "Unauthorized, please check your token or contact the support"}

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
