package auth

import (
	"github.com/labstack/echo/v4"
	problem "github.com/ramsesyok/mosscat/pkg/response"
)

// RolesRequired ensures the user has at least one of the allowed roles.
func RolesRequired(allowed ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := GetClaims(c)
			if claims == nil {
				return problem.Forbidden(c, "RBAC_DENY", "missing claims")
			}
			for _, r := range claims.Roles {
				for _, a := range allowed {
					if r == a {
						return next(c)
					}
				}
			}
			return problem.Forbidden(c, "RBAC_DENY", "role not allowed")
		}
	}
}
