package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	internalErrors "github.com/gmerten/accounts_transactions/internal/error"
	"github.com/stretchr/testify/assert"
)

func TestHandleError_CustomError(t *testing.T) {
	customErr := internalErrors.ValidationError{
		Message: "Invalid request",
	}

	rr := httptest.NewRecorder()

	HandleError(rr, customErr)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleError_GenericError(t *testing.T) {
	genericErr := errors.New("something went wrong")
	rr := httptest.NewRecorder()

	HandleError(rr, genericErr)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
