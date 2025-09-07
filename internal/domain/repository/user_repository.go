package repository

import (
	"context"

	"github.com/ramsesyok/mosscat/internal/domain/model"
)

// UserFilter はユーザー一覧取得の条件を表す。
type UserFilter struct {
	Username string
	Role     string
	Page     int
	Size     int
}

// UserRepository はユーザー永続化処理を定義する。
type UserRepository interface {
	Search(ctx context.Context, f UserFilter) ([]model.User, int, error)
	Get(ctx context.Context, id string) (*model.User, error)
	// FindByUsername returns a user for exact username match.
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, u *model.User) error
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id string) error
}
