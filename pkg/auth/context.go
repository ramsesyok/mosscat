package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetClaims retrieves JWT claims from echo context.
func GetClaims(c echo.Context) *Claims {
	v := c.Get("authUser")
	if token, ok := v.(*jwt.Token); ok {
		if cl, ok := token.Claims.(*Claims); ok {
			return cl
		}
	}
	return nil
}
