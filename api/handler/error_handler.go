package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

type CustomError interface {
	StatusCode() int
}

func HandleError(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json")

	var customErr CustomError
	if errors.As(err, &customErr) {
		w.WriteHeader(customErr.StatusCode())
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: customErr.StatusCode(),
			Message:    err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	})

}
