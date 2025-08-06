package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// OssComponentRepository は domrepo.OssComponentRepository の実装。
type OssComponentRepository struct {
	DB *sql.DB
}

// Search はフィルタに合致する OSS コンポーネント一覧を返す。
func (r *OssComponentRepository) Search(ctx context.Context, f domrepo.OssComponentFilter) ([]model.OssComponent, int, error) {
	var args []any
	var wheres []string
	if f.Name != "" {
		args = append(args, "%"+strings.ToLower(f.Name)+"%")
		wheres = append(wheres, "normalized_name LIKE ?")
	}
	if len(f.Layers) > 0 {
		placeholders := make([]string, len(f.Layers))
		for i, l := range f.Layers {
			placeholders[i] = "?"
			args = append(args, l)
		}
		wheres = append(wheres, fmt.Sprintf("EXISTS (SELECT 1 FROM oss_component_layers l WHERE l.oss_id = oc.id AND l.layer IN (%s))", strings.Join(placeholders, ",")))
	}
	if f.Tag != "" {
		wheres = append(wheres, "EXISTS (SELECT 1 FROM oss_component_tags t JOIN tags tg ON t.tag_id = tg.id WHERE t.oss_id = oc.id AND tg.name = ?)")
		args = append(args, f.Tag)
	}
	whereSQL := whereClause(wheres)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM oss_components oc %s", whereSQL)
	row := r.DB.QueryRowContext(ctx, countQuery, args...)
	var total int
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Size
	query := fmt.Sprintf(`SELECT oc.id, oc.name, oc.normalized_name, oc.homepage_url, oc.repository_url, oc.description, oc.primary_language, oc.default_usage_role, oc.deprecated, oc.created_at, oc.updated_at FROM oss_components oc %s ORDER BY oc.created_at DESC LIMIT ? OFFSET ?`, whereSQL)
	argsWithLimit := append(args, f.Size, offset)

	rows, err := r.DB.QueryContext(ctx, query, argsWithLimit...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comps []model.OssComponent
	for rows.Next() {
		var c model.OssComponent
		if err := rows.Scan(&c.ID, &c.Name, &c.NormalizedName, &c.HomepageURL, &c.RepositoryURL, &c.Description, &c.PrimaryLanguage, &c.DefaultUsageRole, &c.Deprecated, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		comps = append(comps, c)
	}
	return comps, total, rows.Err()
}

// Create は新しい OSS コンポーネントを登録する。
func (r *OssComponentRepository) Create(ctx context.Context, c *model.OssComponent) error {
	query := `INSERT INTO oss_components (id, name, normalized_name, homepage_url, repository_url, description, primary_language, default_usage_role, deprecated, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, c.ID, c.Name, c.NormalizedName, c.HomepageURL, c.RepositoryURL, c.Description, c.PrimaryLanguage, c.DefaultUsageRole, c.Deprecated, c.CreatedAt, c.UpdatedAt)
	return err
}

// Update は既存 OSS コンポーネントを更新する。
func (r *OssComponentRepository) Update(ctx context.Context, c *model.OssComponent) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE oss_components SET name = ?, normalized_name = ?, homepage_url = ?, repository_url = ?, description = ?, primary_language = ?, default_usage_role = ?, deprecated = ?, updated_at = ? WHERE id = ?`,
		c.Name, c.NormalizedName, c.HomepageURL, c.RepositoryURL, c.Description, c.PrimaryLanguage, c.DefaultUsageRole, c.Deprecated, c.UpdatedAt, c.ID,
	)
	return err
}

// Get は指定された ID の OSS コンポーネントを返す。
func (r *OssComponentRepository) Get(ctx context.Context, id string) (*model.OssComponent, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT id, name, normalized_name, homepage_url, repository_url, description, primary_language, default_usage_role, deprecated, created_at, updated_at FROM oss_components WHERE id = ?`, id)
	var c model.OssComponent
	if err := row.Scan(&c.ID, &c.Name, &c.NormalizedName, &c.HomepageURL, &c.RepositoryURL, &c.Description, &c.PrimaryLanguage, &c.DefaultUsageRole, &c.Deprecated, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}
