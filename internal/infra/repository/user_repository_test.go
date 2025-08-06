package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

func TestUserRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepository{DB: db}

	f := domrepo.UserFilter{Username: "adm", Page: 1, Size: 10}

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM users WHERE username LIKE ?")
	mock.ExpectQuery(countQuery).WithArgs("%adm%").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE username LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?")
	now := dbtime.DBTime{Time: time.Now()}
	rows := sqlmock.NewRows([]string{"id", "username", "display_name", "email", "password_hash", "roles", "active", "created_at", "updated_at"}).
		AddRow(uuid.NewString(), "admin", nil, nil, "hash", pq.StringArray{"ADMIN"}, true, now, now)
	mock.ExpectQuery(listQuery).WithArgs("%adm%", 10, 0).WillReturnRows(rows)

	res, total, err := repo.Search(context.Background(), f)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, res, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepository{DB: db}

	id := uuid.NewString()
	query := regexp.QuoteMeta("SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE id = ?")
	now := dbtime.DBTime{Time: time.Now()}
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "display_name", "email", "password_hash", "roles", "active", "created_at", "updated_at"}).
		AddRow(id, "admin", nil, nil, "hash", pq.StringArray{"ADMIN"}, true, now, now))

	u, err := repo.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, u.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepository{DB: db}

	username := "admin"
	query := regexp.QuoteMeta("SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE username = ?")
	now := dbtime.DBTime{Time: time.Now()}
	id := uuid.NewString()
	mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "display_name", "email", "password_hash", "roles", "active", "created_at", "updated_at"}).
			AddRow(id, username, nil, nil, "hash", pq.StringArray{"ADMIN"}, true, now, now),
	)

	u, err := repo.FindByUsername(context.Background(), username)
	require.NoError(t, err)
	require.Equal(t, username, u.Username)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepository{DB: db}

	u := &model.User{ID: uuid.NewString(), Username: "admin", Roles: []string{"ADMIN"}, Active: true, CreatedAt: dbtime.DBTime{Time: time.Now()}, UpdatedAt: dbtime.DBTime{Time: time.Now()}}

	query := regexp.QuoteMeta("INSERT INTO users (id, username, display_name, email, password_hash, roles, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(query).WithArgs(u.ID, u.Username, u.DisplayName, u.Email, u.PasswordHash, pq.Array(u.Roles), u.Active, u.CreatedAt, u.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), u)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepository{DB: db}

	u := &model.User{ID: uuid.NewString(), PasswordHash: "h", Roles: []string{"ADMIN"}, Active: true, UpdatedAt: dbtime.DBTime{Time: time.Now()}}

	query := regexp.QuoteMeta("UPDATE users SET display_name = ?, email = ?, password_hash = ?, roles = ?, active = ?, updated_at = ? WHERE id = ?")
	mock.ExpectExec(query).WithArgs(u.DisplayName, u.Email, u.PasswordHash, pq.Array(u.Roles), u.Active, u.UpdatedAt, u.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), u)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserRepository{DB: db}

	id := uuid.NewString()
	query := regexp.QuoteMeta("DELETE FROM users WHERE id = ?")
	mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
