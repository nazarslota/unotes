package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/udholdenhed/unotes/auth/pkg/errors"
)

func NewErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		res := echo.Map{}
		switch err := err.(type) {
		case *echo.HTTPError:
			res = echo.Map{"code": err.Code, "message": err.Message}
			if s, ok := err.Message.(string); ok {
				res["message"] = strings.ToLower(s)
				if e.Debug && err.Internal != nil {
					res["debug"] = strings.ToLower(err.Internal.Error())
				}
			}
		case *errors.HTTPError:
			res = echo.Map{"code": err.Code, "message": err.Message}
			if s, ok := err.Message.(string); ok {
				res["message"] = strings.ToLower(s)
				if e.Debug && err.Internal != nil {
					res["debug"] = strings.ToLower(err.Internal.Error())
				}
			}
		default:
			res = echo.Map{
				"code":    http.StatusInternalServerError,
				"message": strings.ToLower(http.StatusText(http.StatusInternalServerError)),
			}
			if e.Debug {
				res["debug"] = err.Error()
			}
		}

		switch c.Request().Method {
		case http.MethodHead:
			err = c.NoContent(res["code"].(int))
		default:
			err = c.JSON(res["code"].(int), echo.Map{
				"errors": []echo.Map{res},
			})
		}

		if err != nil {
			e.Logger.Error(err)
		}
	}
}
