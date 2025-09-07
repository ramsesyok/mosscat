package repository

import (
	"context"

	"github.com/ramsesyok/mosscat/internal/domain/model"
)

// OssComponentTagRepository は oss_component_tags テーブル操作を定義する。
type OssComponentTagRepository interface {
	// ListByOssID は指定コンポーネントに紐づくタグを作成日時降順で取得する。
	ListByOssID(ctx context.Context, ossID string) ([]model.Tag, error)
	// Replace はタグを指定 ID 群で置き換える。
	Replace(ctx context.Context, ossID string, tagIDs []string) error
}
