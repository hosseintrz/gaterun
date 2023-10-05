package util

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Message string
	Status  int
	BaseErr error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("error : %d - %s \n %v", e.Status, e.Message, e.BaseErr)
}

func NewHTTPError(status int, message string, baseErr error) *HTTPError {
	return &HTTPError{
		Status:  status,
		Message: message,
		BaseErr: baseErr,
	}
}

func NotFoundError(message string, baseErr error) *HTTPError {
	return NewHTTPError(http.StatusNotFound, message, baseErr)
}

func InternalError(message string, baseErr error) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, message, baseErr)
}

func UnAuthorizedError(message string, baseErr error) *HTTPError {
	return NewHTTPError(http.StatusUnauthorized, message, baseErr)
}

func ForbiddenError(message string, baseErr error) *HTTPError {
	return NewHTTPError(http.StatusForbidden, message, baseErr)
}
