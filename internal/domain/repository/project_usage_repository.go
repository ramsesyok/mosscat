package repository

import (
	"context"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/ramsesyok/mosscat/internal/domain/model"
)

// ProjectUsageFilter はプロジェクト利用状況検索の条件を表す。
type ProjectUsageFilter struct {
	ProjectID   string
	ScopeStatus string
	UsageRole   string
	Direct      *bool
	Page        int
	Size        int
}

// ProjectUsageRepository は ProjectUsage の永続化処理を定義する。
type ProjectUsageRepository interface {
	Search(ctx context.Context, f ProjectUsageFilter) ([]model.ProjectUsage, int, error)
	Create(ctx context.Context, u *model.ProjectUsage) error
	Update(ctx context.Context, u *model.ProjectUsage) error
	Delete(ctx context.Context, id string) error
	UpdateScope(ctx context.Context, id string, scopeStatus string, inclusionNote *string, evaluatedAt dbtime.DBTime, evaluatedBy *string) error
}
