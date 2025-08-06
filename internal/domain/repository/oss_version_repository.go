package repository

import (
	"context"

	"github.com/ramsesyok/mosscat/internal/domain/model"
)

// OssVersionFilter は OSS バージョン検索の条件を表す。
type OssVersionFilter struct {
	OssID        string
	ReviewStatus string
	ScopeStatus  string
	Page         int
	Size         int
}

// OssVersionRepository は OSS バージョンの永続化処理を定義する。
type OssVersionRepository interface {
	Search(ctx context.Context, f OssVersionFilter) ([]model.OssVersion, int, error)
	Get(ctx context.Context, id string) (*model.OssVersion, error)
	Create(ctx context.Context, v *model.OssVersion) error
	Update(ctx context.Context, v *model.OssVersion) error
	Delete(ctx context.Context, id string) error
}
