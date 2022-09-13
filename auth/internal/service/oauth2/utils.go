package oauth2

import (
	"net/http"

	"github.com/udholdenhed/unotes/auth/pkg/errors"
)

var (
	ErrInvalidOrExpiredToken = errors.NewHTTPError(http.StatusBadRequest, "invalid or expired token")
)
