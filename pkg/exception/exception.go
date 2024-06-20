package exception

import "net/http"

type Exception interface {
	Status() int
	Message() string
	Error() string
}

type ExceptionErr struct {
	StatusErr  int    `json:"status"`
	MessageErr string `json:"message"`
	ErrorErr   string `json:"error"`
}

// Error implements Exception.
func (e *ExceptionErr) Error() string {
	return e.ErrorErr
}

// Message implements Exception.
func (e *ExceptionErr) Message() string {
	return e.MessageErr
}

// Status implements Exception.
func (e *ExceptionErr) Status() int {
	return e.StatusErr
}

func NewInternalServerError(message string) Exception {
	return &ExceptionErr{
		StatusErr:  http.StatusInternalServerError,
		MessageErr: message,
		ErrorErr:   "INTERNAL_SERVER_ERROR",
	}
}

func NewNotFoundError(message string) Exception {
	return &ExceptionErr{
		StatusErr:  http.StatusNotFound,
		MessageErr: message,
		ErrorErr:   "NOT_FOUND",
	}
}

func NewUnauthenticationError(message string) Exception {
	return &ExceptionErr{
		StatusErr:  http.StatusUnauthorized,
		MessageErr: message,
		ErrorErr:   "UNAUTHENTICATED",
	}
}

func NewBadRequestError(message string) Exception {
	return &ExceptionErr{
		StatusErr:  http.StatusBadRequest,
		MessageErr: message,
		ErrorErr:   "BAD_REQUEST",
	}
}

func NewUnprocessableEntityError(message string) Exception {
	return &ExceptionErr{
		StatusErr:  http.StatusUnprocessableEntity,
		MessageErr: message,
		ErrorErr:   "UNPROCESSABLE_ENTITY",
	}
}

func NewUnauthorizedError(message string) Exception {
	return &ExceptionErr{
		StatusErr:  http.StatusForbidden,
		MessageErr: message,
		ErrorErr:   "UNAUTHORIZED",
	}
}

func NewConflictError(message string) Exception {
	return &ExceptionErr{
		StatusErr:  http.StatusConflict,
		MessageErr: message,
		ErrorErr:   "CONFLICT",
	}
}
