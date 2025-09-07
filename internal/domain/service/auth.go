package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
	"github.com/ramsesyok/mosscat/pkg/auth"
)

var (
	ErrInvalidCredential = errors.New("invalid credential")
	ErrUserDisabled      = errors.New("user disabled")
)

// AuthService handles authentication logic.
type AuthService struct {
	UserRepo domrepo.UserRepository
}

// Authenticate verifies username and password and returns user if valid.
func (s *AuthService) Authenticate(ctx context.Context, username, password string) (*model.User, error) {
	u, err := s.UserRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredential
		}
		return nil, err
	}
	if !u.Active {
		return nil, ErrUserDisabled
	}
	if err := auth.Compare(u.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredential
	}
	return u, nil
}
