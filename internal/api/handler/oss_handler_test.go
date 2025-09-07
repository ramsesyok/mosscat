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

func setupEcho(h *Handler) *echo.Echo {
	e := echo.New()
	gen.RegisterHandlers(e, h)
	return e
}

func TestListOssComponents(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	compRepo := &infrarepo.OssComponentRepository{DB: db}
	layerRepo := &infrarepo.OssComponentLayerRepository{DB: db}
	tagRepo := &infrarepo.OssComponentTagRepository{DB: db}

	h := &Handler{OssComponentRepo: compRepo, OssComponentLayerRepo: layerRepo, OssComponentTagRepo: tagRepo}
	e := setupEcho(h)

	id := uuid.NewString()
	now := time.Now()
	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM oss_components oc WHERE normalized_name LIKE ?")
	mock.ExpectQuery(countQuery).WithArgs("%redis%").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT oc.id, oc.name, oc.normalized_name, oc.homepage_url, oc.repository_url, oc.description, oc.primary_language, oc.default_usage_role, oc.deprecated, oc.created_at, oc.updated_at FROM oss_components oc WHERE normalized_name LIKE ? ORDER BY oc.created_at DESC LIMIT ? OFFSET ?")
	mock.ExpectQuery(listQuery).WithArgs("%redis%", 50, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "normalized_name", "homepage_url", "repository_url", "description", "primary_language", "default_usage_role", "deprecated", "created_at", "updated_at"}).
			AddRow(id, "Redis", "redis", nil, nil, nil, nil, nil, false, now, now))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT layer FROM oss_component_layers WHERE oss_id = ? ORDER BY layer")).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"layer"}))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT tg.id, tg.name, tg.created_at FROM tags tg JOIN oss_component_tags ct ON ct.tag_id = tg.id WHERE ct.oss_id = ? ORDER BY tg.created_at DESC")).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}))

	req := httptest.NewRequest(http.MethodGet, "/oss?name=redis", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.PagedResultOssComponent
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotNil(t, res.Items)
	require.Len(t, *res.Items, 1)
}

func TestGetOssComponent_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	compRepo := &infrarepo.OssComponentRepository{DB: db}
	h := &Handler{OssComponentRepo: compRepo}
	e := setupEcho(h)

	oid := uuid.New()
	query := regexp.QuoteMeta("SELECT id, name, normalized_name, homepage_url, repository_url, description, primary_language, default_usage_role, deprecated, created_at, updated_at FROM oss_components WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(oid.String()).WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/oss/"+oid.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOssComponent_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	compRepo := &infrarepo.OssComponentRepository{DB: db}
	layerRepo := &infrarepo.OssComponentLayerRepository{DB: db}
	tagRepo := &infrarepo.OssComponentTagRepository{DB: db}
	h := &Handler{OssComponentRepo: compRepo, OssComponentLayerRepo: layerRepo, OssComponentTagRepo: tagRepo}
	e := setupEcho(h)

	oid := uuid.NewString()
	now := time.Now()
	query := regexp.QuoteMeta("SELECT id, name, normalized_name, homepage_url, repository_url, description, primary_language, default_usage_role, deprecated, created_at, updated_at FROM oss_components WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(oid).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "normalized_name", "homepage_url", "repository_url", "description", "primary_language", "default_usage_role", "deprecated", "created_at", "updated_at"}).AddRow(oid, "Redis", "redis", nil, nil, nil, nil, nil, false, now, now))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT layer FROM oss_component_layers WHERE oss_id = ? ORDER BY layer")).WithArgs(oid).WillReturnRows(sqlmock.NewRows([]string{"layer"}))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT tg.id, tg.name, tg.created_at FROM tags tg JOIN oss_component_tags ct ON ct.tag_id = tg.id WHERE ct.oss_id = ? ORDER BY tg.created_at DESC")).WithArgs(oid).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}))

	req := httptest.NewRequest(http.MethodGet, "/oss/"+oid, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.OssComponent
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, oid, res.Id.String())
}

func TestGetOssVersion_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	vid := uuid.New()
	query := regexp.QuoteMeta("SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, review_status, last_reviewed_at, scope_status, created_at, updated_at FROM oss_versions WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(vid.String()).WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/oss/"+uuid.New().String()+"/versions/"+vid.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOssVersion_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	vid := uuid.New()
	oid := uuid.New()
	now := time.Now()
	query := regexp.QuoteMeta("SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, review_status, last_reviewed_at, scope_status, created_at, updated_at FROM oss_versions WHERE id = ?")
	mockRows := sqlmock.NewRows([]string{"id", "oss_id", "version", "release_date", "license_expression_raw", "license_concluded", "purl", "cpe_list", "hash_sha256", "review_status", "last_reviewed_at", "scope_status", "created_at", "updated_at"}).
		AddRow(vid.String(), oid.String(), "1.0.0", now, nil, nil, nil, pq.StringArray{}, nil, "draft", nil, "IN_SCOPE", now, now)
	mock.ExpectQuery(query).WithArgs(vid.String()).WillReturnRows(mockRows)

	req := httptest.NewRequest(http.MethodGet, "/oss/"+oid.String()+"/versions/"+vid.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.OssVersion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, vid, res.Id)
}

func TestListOssComponents_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	compRepo := &infrarepo.OssComponentRepository{DB: db}
	h := &Handler{OssComponentRepo: compRepo}
	e := setupEcho(h)

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM oss_components oc")
	mock.ExpectQuery(countQuery).WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodGet, "/oss", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOssComponent_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodPost, "/oss", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateOssComponent_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+uuid.NewString(), strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListOssVersions_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)

	ossID := uuid.NewString()
	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM oss_versions WHERE oss_id = ?")
	mock.ExpectQuery(countQuery).WithArgs(ossID).WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodGet, "/oss/"+ossID+"/versions", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteOssVersion_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)

	vid := uuid.NewString()
	query := regexp.QuoteMeta("DELETE FROM oss_versions WHERE id = ?")
	mock.ExpectExec(query).WithArgs(vid).WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodDelete, "/oss/"+uuid.NewString()+"/versions/"+vid, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOssVersion_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+uuid.NewString()+"/versions/"+uuid.NewString(), strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
