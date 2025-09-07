package repository

import (
	"context"

	"github.com/ramsesyok/mosscat/internal/domain/model"
)

// OssComponentFilter は OSS コンポーネント検索の条件を表す。
type OssComponentFilter struct {
	Name        string   // normalized_name への部分一致
	Layers      []string // OR 条件
	Tag         string   // タグ名の完全一致
	InScopeOnly bool
	Page        int
	Size        int
}

// OssComponentRepository は OSS コンポーネントの永続化処理を定義する。
type OssComponentRepository interface {
	Search(ctx context.Context, f OssComponentFilter) ([]model.OssComponent, int, error)
	Get(ctx context.Context, id string) (*model.OssComponent, error)
	Create(ctx context.Context, c *model.OssComponent) error
	// Update は既存 OSS コンポーネントを更新する。
	Update(ctx context.Context, c *model.OssComponent) error
}
