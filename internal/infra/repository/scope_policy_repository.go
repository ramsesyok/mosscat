package repository

import (
	"context"
	"database/sql"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// ScopePolicyRepository は domrepo.ScopePolicyRepository の実装。
type ScopePolicyRepository struct {
	DB *sql.DB
}

var _ domrepo.ScopePolicyRepository = (*ScopePolicyRepository)(nil)

// Get は現在のスコープポリシーを取得する。存在しない場合は sql.ErrNoRows を返す。
func (r *ScopePolicyRepository) Get(ctx context.Context) (*model.ScopePolicy, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1`)
	var p model.ScopePolicy
	if err := row.Scan(&p.ID, &p.RuntimeRequiredDefaultInScope, &p.ServerEnvIncluded, &p.AutoMarkForksInScope, &p.UpdatedAt, &p.UpdatedBy); err != nil {
		return nil, err
	}
	return &p, nil
}

// Update はスコープポリシーを登録または更新する。
func (r *ScopePolicyRepository) Update(ctx context.Context, p *model.ScopePolicy) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO scope_policies (id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET runtime_required_default_in_scope=excluded.runtime_required_default_in_scope, server_env_included=excluded.server_env_included, auto_mark_forks_in_scope=excluded.auto_mark_forks_in_scope, updated_at=excluded.updated_at, updated_by=excluded.updated_by`, p.ID, p.RuntimeRequiredDefaultInScope, p.ServerEnvIncluded, p.AutoMarkForksInScope, p.UpdatedAt, p.UpdatedBy)
	return err
}
