package errors

import "net/http"

type UnknownError struct {
	Message string
}

func (e UnknownError) Error() string {
	return e.Message
}

func (e UnknownError) StatusCode() int {
	return http.StatusInternalServerError
}

func NewUnknownError(message string) UnknownError {
	return UnknownError{message}
}
