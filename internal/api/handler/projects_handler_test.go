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

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/domain/model"
	infrarepo "github.com/ramsesyok/mosscat/internal/infra/repository"
)

func TestListProjects(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	projRepo := &infrarepo.ProjectRepository{DB: db}
	h := &Handler{ProjectRepo: projRepo}
	e := setupEcho(h)

	id := uuid.NewString()
	now := dbtime.DBTime{Time: time.Now()}
	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM projects WHERE project_code LIKE ?")
	mock.ExpectQuery(countQuery).WithArgs("%PRJ%").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT id, project_code, name, department, manager, delivery_date, description, created_at, updated_at, (SELECT COUNT(*) FROM project_usages u WHERE u.project_id = projects.id) FROM projects WHERE project_code LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?")
	mock.ExpectQuery(listQuery).WithArgs("%PRJ%", 50, 0).WillReturnRows(sqlmock.NewRows([]string{"id", "project_code", "name", "department", "manager", "delivery_date", "description", "created_at", "updated_at", "count"}).AddRow(id, "PRJ-1", "Proj", nil, nil, nil, nil, now, now, 0))

	req := httptest.NewRequest(http.MethodGet, "/projects?code=PRJ", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.PagedResultProject
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotNil(t, res.Items)
	require.Len(t, *res.Items, 1)
}

func TestGetProject_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	projRepo := &infrarepo.ProjectRepository{DB: db}
	h := &Handler{ProjectRepo: projRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	query := regexp.QuoteMeta("SELECT id, project_code, name, department, manager, delivery_date, description, created_at, updated_at, (SELECT COUNT(*) FROM project_usages u WHERE u.project_id = projects.id) FROM projects WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(pid).WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/projects/"+pid, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProject_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	projRepo := &infrarepo.ProjectRepository{DB: db}
	h := &Handler{ProjectRepo: projRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	now := dbtime.DBTime{Time: time.Now()}
	query := regexp.QuoteMeta("SELECT id, project_code, name, department, manager, delivery_date, description, created_at, updated_at, (SELECT COUNT(*) FROM project_usages u WHERE u.project_id = projects.id) FROM projects WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(pid).WillReturnRows(sqlmock.NewRows([]string{"id", "project_code", "name", "department", "manager", "delivery_date", "description", "created_at", "updated_at", "count"}).AddRow(pid, "P1", "Proj", nil, nil, nil, nil, now, now, 0))

	req := httptest.NewRequest(http.MethodGet, "/projects/"+pid, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.Project
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, pid, res.Id.String())
}

func TestCreateProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	projRepo := &infrarepo.ProjectRepository{DB: db}
	h := &Handler{ProjectRepo: projRepo}
	e := setupEcho(h)

	reqBody := `{"projectCode":"P1","name":"Proj"}`
	// We don't check ID/time exactly; just expect exec with any args
	query := regexp.QuoteMeta("INSERT INTO projects (id, project_code, name, department, manager, delivery_date, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	// Since uuid/time generated inside handler, use sqlmock.AnyArg to match in expectation; we used above
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	projRepo := &infrarepo.ProjectRepository{DB: db}
	h := &Handler{ProjectRepo: projRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	query := regexp.QuoteMeta("DELETE FROM projects WHERE id = ?")
	mock.ExpectExec(query).WithArgs(pid).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodDelete, "/projects/"+pid, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListProjects_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.ProjectRepository{DB: db}
	h := &Handler{ProjectRepo: repo}
	e := setupEcho(h)

	query := regexp.QuoteMeta("SELECT COUNT(*) FROM projects WHERE project_code LIKE ?")
	mock.ExpectQuery(query).WithArgs("%P%").WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodGet, "/projects?code=P", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProject_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateProject_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	pid := uuid.NewString()
	req := httptest.NewRequest(http.MethodPatch, "/projects/"+pid, strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	projRepo := &infrarepo.ProjectRepository{DB: db}
	h := &Handler{ProjectRepo: projRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	now := dbtime.DBTime{Time: time.Now()}
	getQuery := regexp.QuoteMeta("SELECT id, project_code, name, department, manager, delivery_date, description, created_at, updated_at, (SELECT COUNT(*) FROM project_usages u WHERE u.project_id = projects.id) FROM projects WHERE id = ?")
	mock.ExpectQuery(getQuery).WithArgs(pid).WillReturnRows(sqlmock.NewRows([]string{"id", "project_code", "name", "department", "manager", "delivery_date", "description", "created_at", "updated_at", "count"}).AddRow(pid, "P1", "Proj", nil, nil, nil, nil, now, now, 0))
	updateQuery := regexp.QuoteMeta("UPDATE projects SET name = ?, department = ?, manager = ?, delivery_date = ?, description = ?, updated_at = ? WHERE id = ?")
	mock.ExpectExec(updateQuery).WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"name":"Upd"}`
	req := httptest.NewRequest(http.MethodPatch, "/projects/"+pid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListProjectUsages(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	usageRepo := &infrarepo.ProjectUsageRepository{DB: db}
	h := &Handler{ProjectUsageRepo: usageRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM project_usages WHERE project_id = ?")
	mock.ExpectQuery(countQuery).WithArgs(pid).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	listQuery := regexp.QuoteMeta("SELECT id, project_id, oss_id, oss_version_id, usage_role, scope_status, inclusion_note, direct_dependency, added_at, evaluated_at, evaluated_by, modified, modification_description, supplier_type, fork_origin_url FROM project_usages WHERE project_id = ? ORDER BY added_at DESC LIMIT ? OFFSET ?")
	now := dbtime.DBTime{Time: time.Now()}
	mock.ExpectQuery(listQuery).WithArgs(pid, 50, 0).WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "oss_id", "oss_version_id", "usage_role", "scope_status", "inclusion_note", "direct_dependency", "added_at", "evaluated_at", "evaluated_by", "modified", "modification_description", "supplier_type", "fork_origin_url"}).AddRow(uuid.NewString(), pid, uuid.NewString(), uuid.NewString(), "RUNTIME_REQUIRED", "IN_SCOPE", nil, true, now, nil, nil, false, nil, nil, nil))

	req := httptest.NewRequest(http.MethodGet, "/projects/"+pid+"/usages", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.PagedResultProjectUsage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotNil(t, res.Items)
	require.Len(t, *res.Items, 1)
}

func TestCreateProjectUsage(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	usageRepo := &infrarepo.ProjectUsageRepository{DB: db}
	policyRepo := &infrarepo.ScopePolicyRepository{DB: db}
	h := &Handler{ProjectUsageRepo: usageRepo, ScopePolicyRepo: policyRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	policyQuery := regexp.QuoteMeta("SELECT id, runtime_required_default_in_scope, server_env_included, auto_mark_forks_in_scope, updated_at, updated_by FROM scope_policies LIMIT 1")
	now := dbtime.DBTime{Time: time.Now()}
	mock.ExpectQuery(policyQuery).WillReturnRows(sqlmock.NewRows([]string{"id", "runtime_required_default_in_scope", "server_env_included", "auto_mark_forks_in_scope", "updated_at", "updated_by"}).AddRow(uuid.NewString(), true, false, false, now, "user"))
	createQuery := regexp.QuoteMeta("INSERT INTO project_usages (id, project_id, oss_id, oss_version_id, usage_role, scope_status, inclusion_note, direct_dependency, added_at, evaluated_at, evaluated_by, modified, modification_description, supplier_type, fork_origin_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(createQuery).WillReturnResult(sqlmock.NewResult(1, 1))

	reqBody := `{"ossId":"` + uuid.NewString() + `","ossVersionId":"` + uuid.NewString() + `","usageRole":"RUNTIME_REQUIRED"}`
	req := httptest.NewRequest(http.MethodPost, "/projects/"+pid+"/usages", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProjectUsage_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	pid := uuid.NewString()
	req := httptest.NewRequest(http.MethodPost, "/projects/"+pid+"/usages", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateProjectUsageScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	usageRepo := &infrarepo.ProjectUsageRepository{DB: db}
	h := &Handler{ProjectUsageRepo: usageRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	uid := uuid.NewString()
	updateQuery := regexp.QuoteMeta("UPDATE project_usages SET scope_status = ?, inclusion_note = ?, evaluated_at = ?, evaluated_by = ? WHERE id = ?")
	mock.ExpectExec(updateQuery).WithArgs("OUT_SCOPE", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), uid).WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"scopeStatus":"OUT_SCOPE","reasonNote":"bad"}`
	req := httptest.NewRequest(http.MethodPatch, "/projects/"+pid+"/usages/"+uid+"/scope", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProjectUsageScope_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	pid := uuid.NewString()
	uid := uuid.NewString()
	req := httptest.NewRequest(http.MethodPatch, "/projects/"+pid+"/usages/"+uid+"/scope", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateProjectUsage(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	usageRepo := &infrarepo.ProjectUsageRepository{DB: db}
	h := &Handler{ProjectUsageRepo: usageRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	uid := uuid.NewString()
	updateQuery := regexp.QuoteMeta("UPDATE project_usages SET oss_version_id = ?, usage_role = ?, direct_dependency = ?, inclusion_note = ?, scope_status = ?, evaluated_at = ?, evaluated_by = ?, modified = ?, modification_description = ?, supplier_type = ?, fork_origin_url = ? WHERE id = ?")
	mock.ExpectExec(updateQuery).WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"usageRole":"DEV_ONLY"}`
	req := httptest.NewRequest(http.MethodPatch, "/projects/"+pid+"/usages/"+uid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProjectUsage_InvalidBody(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	pid := uuid.NewString()
	uid := uuid.NewString()
	req := httptest.NewRequest(http.MethodPatch, "/projects/"+pid+"/usages/"+uid, strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteProjectUsage(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	usageRepo := &infrarepo.ProjectUsageRepository{DB: db}
	h := &Handler{ProjectUsageRepo: usageRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	uid := uuid.NewString()
	deleteQuery := regexp.QuoteMeta("DELETE FROM project_usages WHERE id = ?")
	mock.ExpectExec(deleteQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodDelete, "/projects/"+pid+"/usages/"+uid, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProjectUsage_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	usageRepo := &infrarepo.ProjectUsageRepository{DB: db}
	h := &Handler{ProjectUsageRepo: usageRepo}
	e := setupEcho(h)

	pid := uuid.NewString()
	uid := uuid.NewString()
	query := regexp.QuoteMeta("DELETE FROM project_usages WHERE id = ?")
	mock.ExpectExec(query).WithArgs(uid).WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodDelete, "/projects/"+pid+"/usages/"+uid, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExportProjectArtifacts(t *testing.T) {
	h := &Handler{}
	e := setupEcho(h)
	pid := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/projects/"+pid+"/export?format=csv", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestInitialScopeStatus(t *testing.T) {
	policy := &model.ScopePolicy{RuntimeRequiredDefaultInScope: true, ServerEnvIncluded: false}
	require.Equal(t, string(gen.INSCOPE), initialScopeStatus(policy, "RUNTIME_REQUIRED"))
	require.Equal(t, string(gen.OUTSCOPE), initialScopeStatus(policy, "BUILD_ONLY"))
	require.Equal(t, string(gen.OUTSCOPE), initialScopeStatus(policy, "SERVER_ENV"))

	// nil policy should result in REVIEW_NEEDED for runtime role
	require.Equal(t, string(gen.REVIEWNEEDED), initialScopeStatus(nil, "RUNTIME_REQUIRED"))
}

func TestToProjectUsage(t *testing.T) {
	now := dbtime.DBTime{Time: time.Now()}
	u := model.ProjectUsage{ID: uuid.NewString(), ProjectID: uuid.NewString(), OssID: uuid.NewString(), OssVersionID: uuid.NewString(), UsageRole: "RUNTIME_REQUIRED", ScopeStatus: "IN_SCOPE", DirectDependency: true, AddedAt: now}
	res := toProjectUsage(u)
	require.Equal(t, u.ID, res.Id.String())
	require.Equal(t, u.ProjectID, res.ProjectId.String())
}

func TestToProject(t *testing.T) {
	now := dbtime.DBTime{Time: time.Now()}
	desc := "desc"
	p := model.Project{ID: uuid.NewString(), ProjectCode: "P1", Name: "Proj", Description: &desc, CreatedAt: now, UpdatedAt: now}
	res := toProject(p)
	require.Equal(t, p.ID, res.Id.String())
	require.Equal(t, p.ProjectCode, res.ProjectCode)
	require.NotNil(t, res.Description)
}
