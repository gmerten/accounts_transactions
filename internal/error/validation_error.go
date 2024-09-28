package errors

import "net/http"

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e ValidationError) StatusCode() int {
	return http.StatusBadRequest
}

func NewValidationError(message string) ValidationError {
	return ValidationError{message}
}
