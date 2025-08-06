package auth

import "github.com/golang-jwt/jwt/v5"

// Claims defines JWT claims for access tokens.
type Claims struct {
	Sub      string   `json:"sub"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}
