package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	infrarepo "github.com/ramsesyok/mosscat/internal/infra/repository"
)

func TestListUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &Handler{UserRepo: repo}
	e := setupEcho(h)

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM users WHERE username LIKE ?")
	mock.ExpectQuery(countQuery).WithArgs("%adm%").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE username LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?")
	now := time.Now()
	mock.ExpectQuery(listQuery).WithArgs("%adm%", 50, 0).WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "display_name", "email", "password_hash", "roles", "active", "created_at", "updated_at"}).
			AddRow(uuid.NewString(), "admin", nil, nil, "h", pq.StringArray{"ADMIN"}, true, now, now),
	)

	req := httptest.NewRequest(http.MethodGet, "/users?username=adm", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.PagedResultUser
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotNil(t, res.Items)
	require.Len(t, *res.Items, 1)
}

func TestGetUser_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &Handler{UserRepo: repo}
	e := setupEcho(h)

	id := uuid.New()
	query := regexp.QuoteMeta("SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(id.String()).WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &Handler{UserRepo: repo}
	e := setupEcho(h)

	id := uuid.New()
	now := time.Now()
	query := regexp.QuoteMeta("SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(id.String()).WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "display_name", "email", "password_hash", "roles", "active", "created_at", "updated_at"}).
			AddRow(id.String(), "admin", nil, nil, "h", pq.StringArray{"ADMIN"}, true, now, now),
	)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.User
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, id, res.Id)
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &Handler{UserRepo: repo}
	e := setupEcho(h)

	query := regexp.QuoteMeta("INSERT INTO users (id, username, display_name, email, password_hash, roles, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"username":"adm","roles":["ADMIN"]}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateUser_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	id := uuid.NewString()
	req := httptest.NewRequest(http.MethodPatch, "/users/"+id, strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &Handler{UserRepo: repo}
	e := setupEcho(h)

	id := uuid.NewString()
	query := regexp.QuoteMeta("DELETE FROM users WHERE id = ?")
	mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodDelete, "/users/"+id, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &Handler{UserRepo: repo}
	e := setupEcho(h)

	id := uuid.NewString()
	query := regexp.QuoteMeta("DELETE FROM users WHERE id = ?")
	mock.ExpectExec(query).WithArgs(id).WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+id, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}
