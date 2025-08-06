package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// AuditLogRepository は domrepo.AuditLogRepository の実装。
type AuditLogRepository struct {
	DB *sql.DB
}

var _ domrepo.AuditLogRepository = (*AuditLogRepository)(nil)

// Search は条件に合致する監査ログを作成日時の降順で取得する。
func (r *AuditLogRepository) Search(ctx context.Context, f domrepo.AuditLogFilter) ([]model.AuditLog, error) {
	var args []any
	var wheres []string
	if f.EntityType != nil && *f.EntityType != "" {
		wheres = append(wheres, "entity_type = ?")
		args = append(args, *f.EntityType)
	}
	if f.EntityID != nil && *f.EntityID != "" {
		wheres = append(wheres, "entity_id = ?")
		args = append(args, *f.EntityID)
	}
	if f.From != nil {
		wheres = append(wheres, "created_at >= ?")
		args = append(args, *f.From)
	}
	if f.To != nil {
		wheres = append(wheres, "created_at <= ?")
		args = append(args, *f.To)
	}
	whereSQL := whereClause(wheres)
	query := fmt.Sprintf("SELECT id, entity_type, entity_id, action, user_name, summary, created_at FROM audit_logs %s ORDER BY created_at DESC", whereSQL)
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	for rows.Next() {
		var l model.AuditLog
		var summary sql.NullString
		if err := rows.Scan(&l.ID, &l.EntityType, &l.EntityID, &l.Action, &l.UserName, &summary, &l.CreatedAt); err != nil {
			return nil, err
		}
		l.Summary = strPtr(summary)
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// Create は新しい監査ログを登録する。
func (r *AuditLogRepository) Create(ctx context.Context, l *model.AuditLog) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO audit_logs (id, entity_type, entity_id, action, user_name, summary, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		l.ID, l.EntityType, l.EntityID, l.Action, l.UserName, l.Summary, l.CreatedAt,
	)
	return err
}
