package oauth2

import (
	"net/http"

	"github.com/udholdenhed/unotes/auth-service/pkg/errors"
)

var (
	ErrInvalidOrExpiredToken = errors.NewHTTPError(http.StatusBadRequest, "invalid or expired token")
)
