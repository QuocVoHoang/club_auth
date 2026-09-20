package http_error

import "net/http"

type HTTPError struct {
	Code    int
	Message string
}

func (e HTTPError) Error() string {
	return e.Message
}

func (e HTTPError) StatusCode() int {
	return e.Code
}

func BadRequestError(message string) error {
	return HTTPError{Code: http.StatusBadRequest, Message: message}
}

func UnauthorizedError(message string) error {
	return HTTPError{Code: http.StatusUnauthorized, Message: message}
}

func ConflictError(message string) error {
	return HTTPError{Code: http.StatusConflict, Message: message}
}

func NotFoundError(message string) error {
	return HTTPError{Code: http.StatusNotFound, Message: message}
}
