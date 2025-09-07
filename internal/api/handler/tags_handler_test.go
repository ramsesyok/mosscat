package handler

import (
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
	"github.com/stretchr/testify/require"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	infrarepo "github.com/ramsesyok/mosscat/internal/infra/repository"
)

func TestListTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.TagRepository{DB: db}
	h := &Handler{TagRepo: repo}
	e := setupEcho(h)

	now := time.Now()
	query := regexp.QuoteMeta(`SELECT id, name, created_at FROM tags ORDER BY created_at DESC`)
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(uuid.NewString(), "db", now))

	req := httptest.NewRequest(http.MethodGet, "/tags", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res []gen.Tag
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Len(t, res, 1)
}

func TestListTags_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.TagRepository{DB: db}
	h := &Handler{TagRepo: repo}
	e := setupEcho(h)

	query := regexp.QuoteMeta(`SELECT id, name, created_at FROM tags ORDER BY created_at DESC`)
	mock.ExpectQuery(query).WillReturnError(sqlmock.ErrCancelled)

	req := httptest.NewRequest(http.MethodGet, "/tags", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.TagRepository{DB: db}
	h := &Handler{TagRepo: repo}
	e := setupEcho(h)

	body := `{"name":"db"}`
	query := regexp.QuoteMeta(`INSERT INTO tags (id, name, created_at) VALUES (?, ?, ?)`)
	mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), "db", sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodPost, "/tags", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTag_InvalidBody(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.TagRepository{DB: db}
	h := &Handler{TagRepo: repo}
	e := setupEcho(h)

	req := httptest.NewRequest(http.MethodPost, "/tags", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTag_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.TagRepository{DB: db}
	h := &Handler{TagRepo: repo}
	e := setupEcho(h)

	body := `{"name":"db"}`
	query := regexp.QuoteMeta(`INSERT INTO tags (id, name, created_at) VALUES (?, ?, ?)`)
	mock.ExpectExec(query).WillReturnError(sqlmock.ErrCancelled)

	req := httptest.NewRequest(http.MethodPost, "/tags", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.TagRepository{DB: db}
	h := &Handler{TagRepo: repo}
	e := setupEcho(h)

	id := uuid.NewString()
	query := regexp.QuoteMeta(`DELETE FROM tags WHERE id = ?`)
	mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodDelete, "/tags/"+id, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTag_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.TagRepository{DB: db}
	h := &Handler{TagRepo: repo}
	e := setupEcho(h)

	id := uuid.NewString()
	query := regexp.QuoteMeta(`DELETE FROM tags WHERE id = ?`)
	mock.ExpectExec(query).WithArgs(id).WillReturnError(sqlmock.ErrCancelled)

	req := httptest.NewRequest(http.MethodDelete, "/tags/"+id, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}
