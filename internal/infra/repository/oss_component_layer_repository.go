package repository

import (
	"context"
	"database/sql"

	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// OssComponentLayerRepository は domrepo.OssComponentLayerRepository の実装。
type OssComponentLayerRepository struct {
	DB *sql.DB
}

var _ domrepo.OssComponentLayerRepository = (*OssComponentLayerRepository)(nil)

// ListByOssID は指定されたコンポーネントのレイヤーを名前順で取得する。
func (r *OssComponentLayerRepository) ListByOssID(ctx context.Context, ossID string) ([]string, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT layer FROM oss_component_layers WHERE oss_id = ? ORDER BY layer`,
		ossID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var layers []string
	for rows.Next() {
		var layer string
		if err := rows.Scan(&layer); err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}
	return layers, rows.Err()
}

// Replace は指定コンポーネントのレイヤーを与えられた値で置き換える。
func (r *OssComponentLayerRepository) Replace(ctx context.Context, ossID string, layers []string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM oss_component_layers WHERE oss_id = ?`, ossID); err != nil {
		tx.Rollback()
		return err
	}
	for _, l := range layers {
		if _, err := tx.ExecContext(ctx, `INSERT INTO oss_component_layers (oss_id, layer) VALUES (?, ?)`, ossID, l); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
