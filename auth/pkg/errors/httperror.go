package errors

import (
	"fmt"
	"net/http"
	"strings"
)

// HTTPError represents an httpserver error.
type HTTPError struct {
	Code     int   `json:"code"`
	Message  any   `json:"message"`
	Internal error `json:"-"`
}

var (
	// ErrHTTPInternalServerError represents an internal http server error.
	ErrHTTPInternalServerError = NewHTTPError(http.StatusInternalServerError,
		strings.ToLower(http.StatusText(http.StatusInternalServerError)))
)

// NewHTTPError creates a new UserError instance.
func NewHTTPError(code int, message ...any) *HTTPError {
	e := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		e.Message = message[0]
	}
	return e
}

// Error makes it compatible with `error` interface.
func (e *HTTPError) Error() string {
	if e.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v", e.Code, e.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", e.Code, e.Message, e.Internal)
}

// Unwrap satisfies the Go 1.13 error wrapper interface.
func (e *HTTPError) Unwrap() error {
	return e.Internal
}

// SetInternal sets error to HTTPError.Internal.
func (e *HTTPError) SetInternal(err error) *HTTPError {
	e.Internal = err
	return e
}
