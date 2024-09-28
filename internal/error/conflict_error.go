package errors

import (
	"errors"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return e.Message
}

func (e ConflictError) StatusCode() int {
	return http.StatusConflict
}

func NewConflictError(message string) ConflictError {
	return ConflictError{message}
}

func IsDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) || err.Error() == "UNIQUE constraint failed" {
		return true
	}

	return false
}
