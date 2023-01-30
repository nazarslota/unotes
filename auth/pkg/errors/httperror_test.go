package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPError(t *testing.T) {
	tests := []struct {
		name        string
		err         *HTTPError
		wantCode    int
		wantMessage string
		wantErr     error
		wantString  string
		wantUnwrap  error
	}{
		{
			name:        "basic error",
			err:         NewHTTPError(http.StatusBadRequest, "invalid request"),
			wantCode:    http.StatusBadRequest,
			wantMessage: "invalid request",
			wantErr:     nil,
			wantString:  "code=400, message=invalid request",
			wantUnwrap:  nil,
		},
		{
			name:        "error with underlying cause",
			err:         NewHTTPError(http.StatusInternalServerError).SetInternal(fmt.Errorf("database error")),
			wantCode:    http.StatusInternalServerError,
			wantMessage: "Internal Server Error",
			wantErr:     fmt.Errorf("database error"),
			wantString:  "code=500, message=Internal Server Error, internal=database error",
			wantUnwrap:  fmt.Errorf("database error"),
		},
		{
			name:        "predefined internal server error",
			err:         ErrHTTPInternalServerError,
			wantCode:    http.StatusInternalServerError,
			wantMessage: "Internal Server Error",
			wantErr:     nil,
			wantString:  "code=500, message=Internal Server Error",
			wantUnwrap:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantCode, tt.err.Code, "HTTPError.Code")
			assert.Equal(t, tt.wantMessage, tt.err.Message, "HTTPError.Message")
			assert.Equal(t, tt.wantErr, tt.err.Internal, "HTTPError.Internal")
			assert.Equal(t, tt.wantString, tt.err.Error(), "HTTPError.Error()")
			assert.Equal(t, tt.wantUnwrap, tt.err.Unwrap(), "HTTPError.Unwrap()")
		})
	}
}
