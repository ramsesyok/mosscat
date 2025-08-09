package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// ProjectUsageRepository は domrepo.ProjectUsageRepository の実装。
type ProjectUsageRepository struct {
	DB *sql.DB
}

var _ domrepo.ProjectUsageRepository = (*ProjectUsageRepository)(nil)

// Search は条件に合致する ProjectUsage を取得する。
func (r *ProjectUsageRepository) Search(ctx context.Context, f domrepo.ProjectUsageFilter) ([]model.ProjectUsage, int, error) {
	var args []any
	wheres := []string{"project_id = ?"}
	args = append(args, f.ProjectID)
	if f.ScopeStatus != "" {
		wheres = append(wheres, "scope_status = ?")
		args = append(args, f.ScopeStatus)
	}
	if f.UsageRole != "" {
		wheres = append(wheres, "usage_role = ?")
		args = append(args, f.UsageRole)
	}
	if f.Direct != nil {
		wheres = append(wheres, "direct_dependency = ?")
		args = append(args, *f.Direct)
	}
	whereSQL := whereClause(wheres)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM project_usages %s", whereSQL)
	row := r.DB.QueryRowContext(ctx, countQuery, args...)
	var total int
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Size
	listQuery := fmt.Sprintf(`SELECT id, project_id, oss_id, oss_version_id, usage_role, scope_status, inclusion_note, direct_dependency, added_at, evaluated_at, evaluated_by, modified, modification_description, supplier_type, fork_origin_url FROM project_usages %s ORDER BY added_at DESC LIMIT ? OFFSET ?`, whereSQL)
	argsWithLimit := append(args, f.Size, offset)
	rows, err := r.DB.QueryContext(ctx, listQuery, argsWithLimit...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var usages []model.ProjectUsage
	for rows.Next() {
		var u model.ProjectUsage
		var note, evalBy, modDesc, supplier, fork sql.NullString
		var evalAt sql.NullTime
		if err := rows.Scan(&u.ID, &u.ProjectID, &u.OssID, &u.OssVersionID, &u.UsageRole, &u.ScopeStatus, &note, &u.DirectDependency, &u.AddedAt, &evalAt, &evalBy, &u.Modified, &modDesc, &supplier, &fork); err != nil {
			return nil, 0, err
		}
		u.InclusionNote = strPtr(note)
		u.EvaluatedAt = timePtr(evalAt)
		u.EvaluatedBy = strPtr(evalBy)
		u.ModificationDescription = strPtr(modDesc)
		u.SupplierType = strPtr(supplier)
		u.ForkOriginURL = strPtr(fork)
		usages = append(usages, u)
	}
	return usages, total, rows.Err()
}

// Create は新しい利用情報を登録する。
func (r *ProjectUsageRepository) Create(ctx context.Context, u *model.ProjectUsage) error {
	query := `INSERT INTO project_usages (id, project_id, oss_id, oss_version_id, usage_role, scope_status, inclusion_note, direct_dependency, added_at, evaluated_at, evaluated_by, modified, modification_description, supplier_type, fork_origin_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, u.ID, u.ProjectID, u.OssID, u.OssVersionID, u.UsageRole, u.ScopeStatus, u.InclusionNote, u.DirectDependency, u.AddedAt, u.EvaluatedAt, u.EvaluatedBy, u.Modified, u.ModificationDescription, u.SupplierType, u.ForkOriginURL)
	return err
}

// Update は既存の利用情報を更新する。
func (r *ProjectUsageRepository) Update(ctx context.Context, u *model.ProjectUsage) error {
	query := `UPDATE project_usages SET oss_version_id = ?, usage_role = ?, direct_dependency = ?, inclusion_note = ?, scope_status = ?, evaluated_at = ?, evaluated_by = ?, modified = ?, modification_description = ?, supplier_type = ?, fork_origin_url = ? WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, u.OssVersionID, u.UsageRole, u.DirectDependency, u.InclusionNote, u.ScopeStatus, u.EvaluatedAt, u.EvaluatedBy, u.Modified, u.ModificationDescription, u.SupplierType, u.ForkOriginURL, u.ID)
	return err
}

// Delete は ID を指定して利用情報を削除する。
func (r *ProjectUsageRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM project_usages WHERE id = ?`, id)
	return err
}

// UpdateScope はスコープ関連のフィールドのみ更新する。
func (r *ProjectUsageRepository) UpdateScope(ctx context.Context, id string, scopeStatus string, inclusionNote *string, evaluatedAt dbtime.DBTime, evaluatedBy *string) error {
	query := `UPDATE project_usages SET scope_status = ?, inclusion_note = ?, evaluated_at = ?, evaluated_by = ? WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, scopeStatus, inclusionNote, evaluatedAt, evaluatedBy, id)
	return err
}
