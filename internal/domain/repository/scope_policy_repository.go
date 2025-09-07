package repository

import (
	"context"

	"github.com/ramsesyok/mosscat/internal/domain/model"
)

// ScopePolicyRepository は ScopePolicy の永続化処理を定義する。
type ScopePolicyRepository interface {
	// Get は現在のポリシーを取得する。存在しない場合は nil を返す。
	Get(ctx context.Context) (*model.ScopePolicy, error)
	// Update はポリシーを登録または更新する。
	Update(ctx context.Context, p *model.ScopePolicy) error
}
