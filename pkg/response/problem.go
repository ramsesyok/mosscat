package problem

import (
	"net/http"

	"github.com/labstack/echo/v4"
	gen "github.com/ramsesyok/mosscat/internal/api/gen"
)

func respond(c echo.Context, status int, title, code, detail string) error {
	p := gen.Problem{Title: title, Status: status}
	if code != "" {
		p.Code = &code
	}
	if detail != "" {
		p.Detail = &detail
	}
	return c.JSON(status, p)
}

// Unauthorized returns 401 Problem JSON.
func Unauthorized(c echo.Context, code, detail string) error {
	return respond(c, http.StatusUnauthorized, "UNAUTHORIZED", code, detail)
}

// Forbidden returns 403 Problem JSON.
func Forbidden(c echo.Context, code, detail string) error {
	return respond(c, http.StatusForbidden, "FORBIDDEN", code, detail)
}
