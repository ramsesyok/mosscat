package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
	"github.com/ramsesyok/mosscat/pkg/auth"
)

type stubUserRepo struct {
	user *model.User
	err  error
}

func (s *stubUserRepo) Search(ctx context.Context, f domrepo.UserFilter) ([]model.User, int, error) {
	return nil, 0, nil
}
func (s *stubUserRepo) Get(ctx context.Context, id string) (*model.User, error) { return nil, nil }
func (s *stubUserRepo) Create(ctx context.Context, u *model.User) error         { return nil }
func (s *stubUserRepo) Update(ctx context.Context, u *model.User) error         { return nil }
func (s *stubUserRepo) Delete(ctx context.Context, id string) error             { return nil }
func (s *stubUserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.user, s.err
}

func TestAuthenticate_Success(t *testing.T) {
	hash, _ := auth.Hash("pass")
	repo := &stubUserRepo{user: &model.User{ID: uuid.NewString(), Username: "a", PasswordHash: hash, Active: true}}
	svc := AuthService{UserRepo: repo}
	u, err := svc.Authenticate(context.Background(), "a", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == nil || u.Username != "a" {
		t.Fatalf("unexpected user: %#v", u)
	}
}

func TestAuthenticate_UserDisabled(t *testing.T) {
	hash, _ := auth.Hash("pass")
	repo := &stubUserRepo{user: &model.User{ID: uuid.NewString(), Username: "a", PasswordHash: hash, Active: false}}
	svc := AuthService{UserRepo: repo}
	_, err := svc.Authenticate(context.Background(), "a", "pass")
	if err != ErrUserDisabled {
		t.Fatalf("expected ErrUserDisabled, got %v", err)
	}
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	hash, _ := auth.Hash("pass")
	repo := &stubUserRepo{user: &model.User{ID: uuid.NewString(), Username: "a", PasswordHash: hash, Active: true}}
	svc := AuthService{UserRepo: repo}
	_, err := svc.Authenticate(context.Background(), "a", "wrong")
	if err != ErrInvalidCredential {
		t.Fatalf("expected ErrInvalidCredential, got %v", err)
	}
}

func TestAuthenticate_NotFound(t *testing.T) {
	repo := &stubUserRepo{err: sql.ErrNoRows}
	svc := AuthService{UserRepo: repo}
	_, err := svc.Authenticate(context.Background(), "a", "p")
	if err != ErrInvalidCredential {
		t.Fatalf("expected ErrInvalidCredential, got %v", err)
	}
}
