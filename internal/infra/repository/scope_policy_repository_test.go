package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ramsesyok/mosscat/internal/domain/model"
)

func TestScopePolicyRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ScopePolicyRepository{DB: db}

	query := regexp.QuoteMeta(`SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1`)
	now := dbtime.DBTime{Time: time.Now()}
	rows := sqlmock.NewRows([]string{"id", "runtime_required_default_in_scope", "server_env_included", "auto_mark_forks_in_scope", "updated_at", "updated_by"}).
		AddRow(uuid.NewString(), true, false, true, now, "user")
	mock.ExpectQuery(query).WillReturnRows(rows)

	p, err := repo.Get(context.Background())
	require.NoError(t, err)
	require.NotNil(t, p)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestScopePolicyRepository_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ScopePolicyRepository{DB: db}

	query := regexp.QuoteMeta(`SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1`)
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, err = repo.Get(context.Background())
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestScopePolicyRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ScopePolicyRepository{DB: db}

	p := &model.ScopePolicy{
		ID:                            uuid.NewString(),
		RuntimeRequiredDefaultInScope: true,
		ServerEnvIncluded:             false,
		AutoMarkForksInScope:          true,
		UpdatedAt:                     dbtime.DBTime{Time: time.Now()},
		UpdatedBy:                     "user",
	}

	query := regexp.QuoteMeta(`INSERT INTO scope_policies (id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET runtime_required_default_in_scope=excluded.runtime_required_default_in_scope, server_env_included=excluded.server_env_included, auto_mark_forks_in_scope=excluded.auto_mark_forks_in_scope, updated_at=excluded.updated_at, updated_by=excluded.updated_by`)
	mock.ExpectExec(query).WithArgs(p.ID, p.RuntimeRequiredDefaultInScope, p.ServerEnvIncluded, p.AutoMarkForksInScope, p.UpdatedAt, p.UpdatedBy).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), p)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
