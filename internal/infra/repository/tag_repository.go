package repository

import (
	"context"
	"database/sql"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// TagRepository は domrepo.TagRepository の実装。
type TagRepository struct {
	DB *sql.DB
}

var _ domrepo.TagRepository = (*TagRepository)(nil)

// List はすべてのタグを作成日時降順で取得する。
func (r *TagRepository) List(ctx context.Context) ([]model.Tag, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, created_at FROM tags ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var t model.Tag
		var created dbtime.DBTime
		if err := rows.Scan(&t.ID, &t.Name, &created); err != nil {
			return nil, err
		}
		t.CreatedAt = &created
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

// Create は新しいタグを登録する。
func (r *TagRepository) Create(ctx context.Context, t *model.Tag) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO tags (id, name, created_at) VALUES (?, ?, ?)`,
		t.ID, t.Name, t.CreatedAt,
	)
	return err
}

// Delete は ID 指定でタグを削除する。
func (r *TagRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM tags WHERE id = ?`, id)
	return err
}
