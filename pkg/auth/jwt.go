package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ramsesyok/mosscat/internal/domain/model"
)

// GenerateToken generates a signed JWT for the given user.
// It returns the token string and expiresIn seconds.
func GenerateToken(u *model.User) (string, int, error) {
	expMin := 15
	if s := os.Getenv("JWT_EXPIRES_MIN"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			expMin = v
		}
	}
	now := time.Now()
	claims := Claims{
		Sub:      u.ID,
		Username: u.Username,
		Roles:    u.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expMin) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}
	return signed, expMin * 60, nil
}
