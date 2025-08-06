package repository

import (
	"context"

	"github.com/ramsesyok/mosscat/internal/domain/model"
)

// TagRepository はタグの永続化処理を定義する。
type TagRepository interface {
	// List は全てのタグを作成日時降順で返す。
	List(ctx context.Context) ([]model.Tag, error)
	// Create は新しいタグを登録する。
	Create(ctx context.Context, t *model.Tag) error
	// Delete は指定 ID のタグを削除する。
	Delete(ctx context.Context, id string) error
}
