package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// ProjectRepository は domrepo.ProjectRepository の実装。
type ProjectRepository struct {
	DB *sql.DB
}

var _ domrepo.ProjectRepository = (*ProjectRepository)(nil)

// Search は条件に合致するプロジェクト一覧を利用数付きで返す。
func (r *ProjectRepository) Search(ctx context.Context, f domrepo.ProjectFilter) ([]model.Project, int, error) {
	var args []any
	var wheres []string
	if f.Code != "" {
		wheres = append(wheres, "project_code LIKE ?")
		args = append(args, "%"+f.Code+"%")
	}
	if f.Name != "" {
		wheres = append(wheres, "name LIKE ?")
		args = append(args, "%"+f.Name+"%")
	}
	whereSQL := whereClause(wheres)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM projects %s", whereSQL)
	row := r.DB.QueryRowContext(ctx, countQuery, args...)
	var total int
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Size
	listQuery := fmt.Sprintf(`SELECT id, project_code, name, department, manager, delivery_date, description, created_at, updated_at, (SELECT COUNT(*) FROM project_usages u WHERE u.project_id = projects.id) FROM projects %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereSQL)
	argsWithLimit := append(args, f.Size, offset)
	rows, err := r.DB.QueryContext(ctx, listQuery, argsWithLimit...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		var dept, mgr, desc sql.NullString
		var delivery sql.NullTime
		var usageCount int
		if err := rows.Scan(&p.ID, &p.ProjectCode, &p.Name, &dept, &mgr, &delivery, &desc, &p.CreatedAt, &p.UpdatedAt, &usageCount); err != nil {
			return nil, 0, err
		}
		p.Department = strPtr(dept)
		p.Manager = strPtr(mgr)
		p.DeliveryDate = timePtr(delivery)
		p.Description = strPtr(desc)
		p.OssUsageCount = usageCount
		projects = append(projects, p)
	}
	return projects, total, rows.Err()
}

// Get は ID を指定してプロジェクトを取得する。
func (r *ProjectRepository) Get(ctx context.Context, id string) (*model.Project, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT id, project_code, name, department, manager, delivery_date, description, created_at, updated_at, (SELECT COUNT(*) FROM project_usages u WHERE u.project_id = projects.id) FROM projects WHERE id = ?`, id)
	var p model.Project
	var dept, mgr, desc sql.NullString
	var delivery sql.NullTime
	var usageCount int
	if err := row.Scan(&p.ID, &p.ProjectCode, &p.Name, &dept, &mgr, &delivery, &desc, &p.CreatedAt, &p.UpdatedAt, &usageCount); err != nil {
		return nil, err
	}
	p.Department = strPtr(dept)
	p.Manager = strPtr(mgr)
	p.DeliveryDate = timePtr(delivery)
	p.Description = strPtr(desc)
	p.OssUsageCount = usageCount
	return &p, nil
}

// Create は新しいプロジェクトを登録する。
func (r *ProjectRepository) Create(ctx context.Context, p *model.Project) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO projects (id, project_code, name, department, manager, delivery_date, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, p.ID, p.ProjectCode, p.Name, p.Department, p.Manager, p.DeliveryDate, p.Description, p.CreatedAt, p.UpdatedAt)
	return err
}

// Update は既存プロジェクトを更新する。
func (r *ProjectRepository) Update(ctx context.Context, p *model.Project) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE projects SET name = ?, department = ?, manager = ?, delivery_date = ?, description = ?, updated_at = ? WHERE id = ?`, p.Name, p.Department, p.Manager, p.DeliveryDate, p.Description, p.UpdatedAt, p.ID)
	return err
}

// Delete は ID 指定でプロジェクトを削除する。
func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM projects WHERE id = ?`, id)
	return err
}
