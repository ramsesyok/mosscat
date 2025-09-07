package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	infrarepo "github.com/ramsesyok/mosscat/internal/infra/repository"
)

func TestGetScopePolicy_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.ScopePolicyRepository{DB: db}
	h := &Handler{ScopePolicyRepo: repo}
	e := setupEcho(h)

	query := regexp.QuoteMeta("SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1")
	mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/scope/policy", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetScopePolicy_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.ScopePolicyRepository{DB: db}
	h := &Handler{ScopePolicyRepo: repo}
	e := setupEcho(h)

	now := time.Now()
	pid := uuid.NewString()
	query := regexp.QuoteMeta("SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1")
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "runtime_required_default_in_scope", "server_env_included", "auto_mark_forks_in_scope", "updated_at", "updated_by"}).AddRow(pid, true, false, true, now, "u"))

	req := httptest.NewRequest(http.MethodGet, "/scope/policy", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var res gen.ScopePolicy
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, pid, res.Id.String())
	require.True(t, *res.RuntimeRequiredDefaultInScope)
	require.False(t, *res.ServerEnvIncluded)
	require.True(t, *res.AutoMarkForksInScope)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateScopePolicy_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.ScopePolicyRepository{DB: db}
	h := &Handler{ScopePolicyRepo: repo}
	e := setupEcho(h)

	getQuery := regexp.QuoteMeta("SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1")
	mock.ExpectQuery(getQuery).WillReturnError(sql.ErrNoRows)

	updateQuery := regexp.QuoteMeta("INSERT INTO scope_policies (id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET runtime_required_default_in_scope=excluded.runtime_required_default_in_scope, server_env_included=excluded.server_env_included, auto_mark_forks_in_scope=excluded.auto_mark_forks_in_scope, updated_at=excluded.updated_at, updated_by=excluded.updated_by")
	mock.ExpectExec(updateQuery).WithArgs(sqlmock.AnyArg(), true, false, false, sqlmock.AnyArg(), "api-user").WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"runtimeRequiredDefaultInScope":true}`
	req := httptest.NewRequest(http.MethodPatch, "/scope/policy", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateScopePolicy_UpdateExisting(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.ScopePolicyRepository{DB: db}
	h := &Handler{ScopePolicyRepo: repo}
	e := setupEcho(h)

	now := time.Now()
	pid := uuid.NewString()
	getQuery := regexp.QuoteMeta("SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1")
	mock.ExpectQuery(getQuery).WillReturnRows(sqlmock.NewRows([]string{"id", "runtime_required_default_in_scope", "server_env_included", "auto_mark_forks_in_scope", "updated_at", "updated_by"}).AddRow(pid, false, false, false, now, "u"))

	updateQuery := regexp.QuoteMeta("INSERT INTO scope_policies (id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET runtime_required_default_in_scope=excluded.runtime_required_default_in_scope, server_env_included=excluded.server_env_included, auto_mark_forks_in_scope=excluded.auto_mark_forks_in_scope, updated_at=excluded.updated_at, updated_by=excluded.updated_by")
	mock.ExpectExec(updateQuery).WithArgs(pid, false, true, false, sqlmock.AnyArg(), "api-user").WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"serverEnvIncluded":true}`
	req := httptest.NewRequest(http.MethodPatch, "/scope/policy", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateScopePolicy_InvalidBody(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.ScopePolicyRepository{DB: db}
	h := &Handler{ScopePolicyRepo: repo}
	e := setupEcho(h)

	req := httptest.NewRequest(http.MethodPatch, "/scope/policy", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}
func TestGetScopePolicy_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.ScopePolicyRepository{DB: db}
	h := &Handler{ScopePolicyRepo: repo}
	e := setupEcho(h)

	query := regexp.QuoteMeta("SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1")
	mock.ExpectQuery(query).WillReturnError(errors.New("fail"))

	req := httptest.NewRequest(http.MethodGet, "/scope/policy", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateScopePolicy_GetError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.ScopePolicyRepository{DB: db}
	h := &Handler{ScopePolicyRepo: repo}
	e := setupEcho(h)

	getQuery := regexp.QuoteMeta("SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1")
	mock.ExpectQuery(getQuery).WillReturnError(errors.New("fail"))

	body := `{"serverEnvIncluded":true}`
	req := httptest.NewRequest(http.MethodPatch, "/scope/policy", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateScopePolicy_UpdateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.ScopePolicyRepository{DB: db}
	h := &Handler{ScopePolicyRepo: repo}
	e := setupEcho(h)

	getQuery := regexp.QuoteMeta("SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1")
	mock.ExpectQuery(getQuery).WillReturnError(sql.ErrNoRows)

	updateQuery := regexp.QuoteMeta("INSERT INTO scope_policies (id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO UPDATE SET runtime_required_default_in_scope=excluded.runtime_required_default_in_scope, server_env_included=excluded.server_env_included, auto_mark_forks_in_scope=excluded.auto_mark_forks_in_scope, updated_at=excluded.updated_at, updated_by=excluded.updated_by")
	mock.ExpectExec(updateQuery).WillReturnError(errors.New("fail"))

	body := `{"serverEnvIncluded":true}`
	req := httptest.NewRequest(http.MethodPatch, "/scope/policy", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}
