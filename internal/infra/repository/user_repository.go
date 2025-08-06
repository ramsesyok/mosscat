package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// UserRepository は domrepo.UserRepository の実装。
type UserRepository struct {
	DB *sql.DB
}

var _ domrepo.UserRepository = (*UserRepository)(nil)

// Search は条件に合致するユーザー一覧を返す。
func (r *UserRepository) Search(ctx context.Context, f domrepo.UserFilter) ([]model.User, int, error) {
	var args []any
	var wheres []string
	if f.Username != "" {
		wheres = append(wheres, "username LIKE ?")
		args = append(args, "%"+f.Username+"%")
	}
	if f.Role != "" {
		wheres = append(wheres, "? = ANY(roles)")
		args = append(args, f.Role)
	}
	whereSQL := whereClause(wheres)

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereSQL)
	row := r.DB.QueryRowContext(ctx, countQuery, args...)
	var total int
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Size
	listQuery := fmt.Sprintf(`SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereSQL)
	argsWithLimit := append(args, f.Size, offset)
	rows, err := r.DB.QueryContext(ctx, listQuery, argsWithLimit...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var display, email sql.NullString
		var roles pq.StringArray
		if err := rows.Scan(&u.ID, &u.Username, &display, &email, &u.PasswordHash, &roles, &u.Active, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		u.DisplayName = strPtr(display)
		u.Email = strPtr(email)
		u.Roles = []string(roles)
		users = append(users, u)
	}
	return users, total, rows.Err()
}

// Get は ID 指定でユーザーを取得する。
func (r *UserRepository) Get(ctx context.Context, id string) (*model.User, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE id = ?`, id)
	var u model.User
	var display, email sql.NullString
	var roles pq.StringArray
	if err := row.Scan(&u.ID, &u.Username, &display, &email, &u.PasswordHash, &roles, &u.Active, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	u.DisplayName = strPtr(display)
	u.Email = strPtr(email)
	u.Roles = []string(roles)
	return &u, nil
}

// FindByUsername returns a user by username.
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE username = ?`, username)
	var u model.User
	var display, email sql.NullString
	var roles pq.StringArray
	if err := row.Scan(&u.ID, &u.Username, &display, &email, &u.PasswordHash, &roles, &u.Active, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	u.DisplayName = strPtr(display)
	u.Email = strPtr(email)
	u.Roles = []string(roles)
	return &u, nil
}

// Create は新しいユーザーを登録する。
func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO users (id, username, display_name, email, password_hash, roles, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.DisplayName, u.Email, u.PasswordHash, pq.Array(u.Roles), u.Active, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

// Update は既存ユーザーを更新する。
func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE users SET display_name = ?, email = ?, password_hash = ?, roles = ?, active = ?, updated_at = ? WHERE id = ?`,
		u.DisplayName, u.Email, u.PasswordHash, pq.Array(u.Roles), u.Active, u.UpdatedAt, u.ID,
	)
	return err
}

// Delete は ID 指定でユーザーを削除する。
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	return err
}
