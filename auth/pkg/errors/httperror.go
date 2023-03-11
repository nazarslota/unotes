// Package errors implements a custom error type HTTPError.
package errors

import (
	"fmt"
	"net/http"
)

// HTTPError is a custom error type that holds information about an HTTP error.
// It contains the HTTP status code, a message, and an internal error if there is one.
type HTTPError struct {
	Code     int   `json:"code"`    // HTTP status code.
	Message  any   `json:"message"` // Error message.
	Internal error `json:"-"`       // Internal error.
}

// ErrHTTPInternalServerError is a predefined HTTPError with status code 500 and message "Internal Server Error".
var ErrHTTPInternalServerError = NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

// NewHTTPError returns a new HTTPError with the given status code and message.
// If message is not provided, it will use the default message for the status code.
func NewHTTPError(code int, message ...any) *HTTPError {
	e := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		e.Message = message[0]
	}
	return e
}

// Error returns a string representation of the HTTPError.
// It includes the HTTP status code, message, and internal error if there is one.
func (e *HTTPError) Error() string {
	if e.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v", e.Code, e.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", e.Code, e.Message, e.Internal)
}

// Unwrap returns the internal error of the HTTPError.
func (e *HTTPError) Unwrap() error {
	return e.Internal
}

// SetInternal sets the internal error of the HTTPError.
func (e *HTTPError) SetInternal(err error) *HTTPError {
	e.Internal = err
	return e
}
