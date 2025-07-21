package errors

import "net/http"

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func New(code int, message string) *ApiError {
	return &ApiError{Code: code, Message: message}
}

// Predefined errors
var (
	// 4xx Client Errors
	ErrNotFound              = New(http.StatusNotFound, "resource not found")
	ErrBadRequest            = New(http.StatusBadRequest, "bad request")
	ErrAlreadyRefunded       = New(http.StatusBadRequest, "Ticket is already refunded")
	ErrSeatsUnavailable      = New(http.StatusConflict, "requested seats are no longer available")
	ErrDuplicateSocialCode   = New(http.StatusBadRequest, "ducpicate social code")
	ErrSeatLessThanPassenger = New(http.StatusBadRequest, "seats are less than the number of passangers")
	// 5xx Client Errors
	ErrNotCreated   = New(http.StatusInternalServerError, "object not created")
	ErrServerError  = New(http.StatusInternalServerError, "internal server error")
	ErrCreateTicket = New(http.StatusInternalServerError, "ticket not created")
)
