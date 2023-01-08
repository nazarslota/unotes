package errors

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
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
			if gotCode := tt.err.Code; gotCode != tt.wantCode {
				t.Errorf("HTTPError.Code = %v, want %v", gotCode, tt.wantCode)
			}

			if gotMessage := tt.err.Message; gotMessage != tt.wantMessage {
				t.Errorf("HTTPError.Message = %v, want %v", gotMessage, tt.wantMessage)
			}

			if gotErr := tt.err.Internal; !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("HTTPError.Internal = %v, want %v", gotErr, tt.wantErr)
			}

			if gotString := tt.err.Error(); gotString != tt.wantString {
				t.Errorf("HTTPError.Error() = %v, want %v", gotString, tt.wantString)
			}

			if gotUnwrap := tt.err.Unwrap(); !reflect.DeepEqual(gotUnwrap, tt.wantUnwrap) {
				t.Errorf("HTTPError.Unwrap() = %v, want %v", gotUnwrap, tt.wantUnwrap)
			}
		})
	}
}
