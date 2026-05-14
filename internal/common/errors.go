package common

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(status int, code string, message string, err error) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func BadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func Validation(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", message, err)
}

func NotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, "NOT_FOUND", message, nil)
}

func MethodNotAllowed(message string) *AppError {
	return NewAppError(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", message, nil)
}

func Internal(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", message, err)
}

func ToAppError(err error) *AppError {
	var appError *AppError
	if errors.As(err, &appError) {
		return appError
	}

	return Internal("internal server error", err)
}
