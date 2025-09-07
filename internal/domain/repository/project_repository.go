package repository

import (
	"context"

	"github.com/ramsesyok/mosscat/internal/domain/model"
)

// ProjectFilter はプロジェクト一覧取得の条件を表す。
type ProjectFilter struct {
	Code string
	Name string
	Page int
	Size int
}

// ProjectRepository はプロジェクトの永続化処理を定義する。
type ProjectRepository interface {
	Search(ctx context.Context, f ProjectFilter) ([]model.Project, int, error)
	Get(ctx context.Context, id string) (*model.Project, error)
	Create(ctx context.Context, p *model.Project) error
	Update(ctx context.Context, p *model.Project) error
	Delete(ctx context.Context, id string) error
}
