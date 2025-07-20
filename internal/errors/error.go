package errors

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Predefined errors
var (
	ErrNotFound         = New(http.StatusNotFound, "resource not found")
	ErrBadRequest       = New(http.StatusBadRequest, "bad request")
	ErrServerError      = New(http.StatusInternalServerError, "internal server error")
	ErrNotCreated       = New(http.StatusInternalServerError, "object not created")
	ErrSeatsUnavailable = New(http.StatusNoContent, "Seats are not available")
)
