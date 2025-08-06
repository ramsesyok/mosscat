package repository

import "context"

// OssComponentLayerRepository は oss_component_layers テーブル操作を定義する。
type OssComponentLayerRepository interface {
	// ListByOssID は指定コンポーネントに紐づくレイヤーを取得する。
	ListByOssID(ctx context.Context, ossID string) ([]string, error)
	// Replace はレイヤーを置き換える。
	Replace(ctx context.Context, ossID string, layers []string) error
}
