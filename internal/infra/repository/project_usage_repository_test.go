package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

func TestProjectUsageRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectUsageRepository{DB: db}

	pid := uuid.NewString()
	f := domrepo.ProjectUsageFilter{ProjectID: pid, Page: 1, Size: 10}

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM project_usages WHERE project_id = ?")
	mock.ExpectQuery(countQuery).WithArgs(pid).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT id, project_id, oss_id, oss_version_id, usage_role, scope_status, inclusion_note, direct_dependency, added_at, evaluated_at, evaluated_by, modified, modification_description, supplier_type, fork_origin_url FROM project_usages WHERE project_id = ? ORDER BY added_at DESC LIMIT ? OFFSET ?")
	now := dbtime.DBTime{Time: time.Now()}
	rows := sqlmock.NewRows([]string{"id", "project_id", "oss_id", "oss_version_id", "usage_role", "scope_status", "inclusion_note", "direct_dependency", "added_at", "evaluated_at", "evaluated_by", "modified", "modification_description", "supplier_type", "fork_origin_url"}).
		AddRow(uuid.NewString(), pid, uuid.NewString(), uuid.NewString(), "RUNTIME_REQUIRED", "IN_SCOPE", nil, true, now, nil, nil, false, nil, nil, nil)
	mock.ExpectQuery(listQuery).WithArgs(pid, 10, 0).WillReturnRows(rows)

	res, total, err := repo.Search(context.Background(), f)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, res, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectUsageRepository_Search_WithFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectUsageRepository{DB: db}
	pid := uuid.NewString()
	direct := true
	f := domrepo.ProjectUsageFilter{ProjectID: pid, UsageRole: "RUNTIME_REQUIRED", Direct: &direct, Page: 1, Size: 10}

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM project_usages WHERE project_id = ? AND usage_role = ? AND direct_dependency = ?")
	mock.ExpectQuery(countQuery).WithArgs(pid, "RUNTIME_REQUIRED", direct).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	listQuery := regexp.QuoteMeta("SELECT id, project_id, oss_id, oss_version_id, usage_role, scope_status, inclusion_note, direct_dependency, added_at, evaluated_at, evaluated_by, modified, modification_description, supplier_type, fork_origin_url FROM project_usages WHERE project_id = ? AND usage_role = ? AND direct_dependency = ? ORDER BY added_at DESC LIMIT ? OFFSET ?")
	mock.ExpectQuery(listQuery).WithArgs(pid, "RUNTIME_REQUIRED", direct, 10, 0).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, _, err = repo.Search(context.Background(), f)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectUsageRepository_Search_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectUsageRepository{DB: db}
	pid := uuid.NewString()
	f := domrepo.ProjectUsageFilter{ProjectID: pid, Page: 1, Size: 10}

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM project_usages WHERE project_id = ?")
	mock.ExpectQuery(countQuery).WithArgs(pid).WillReturnError(errors.New("fail"))

	_, _, err = repo.Search(context.Background(), f)
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectUsageRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectUsageRepository{DB: db}

	u := &model.ProjectUsage{
		ID:               uuid.NewString(),
		ProjectID:        uuid.NewString(),
		OssID:            uuid.NewString(),
		OssVersionID:     uuid.NewString(),
		UsageRole:        "RUNTIME_REQUIRED",
		ScopeStatus:      "IN_SCOPE",
		DirectDependency: true,
		AddedAt:          dbtime.DBTime{Time: time.Now()},
		Modified:         false,
	}

	query := regexp.QuoteMeta("INSERT INTO project_usages (id, project_id, oss_id, oss_version_id, usage_role, scope_status, inclusion_note, direct_dependency, added_at, evaluated_at, evaluated_by, modified, modification_description, supplier_type, fork_origin_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(query).
		WithArgs(u.ID, u.ProjectID, u.OssID, u.OssVersionID, u.UsageRole, u.ScopeStatus, u.InclusionNote, u.DirectDependency, u.AddedAt, u.EvaluatedAt, u.EvaluatedBy, u.Modified, u.ModificationDescription, u.SupplierType, u.ForkOriginURL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), u)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectUsageRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectUsageRepository{DB: db}

	u := &model.ProjectUsage{
		ID:               uuid.NewString(),
		OssVersionID:     uuid.NewString(),
		UsageRole:        "BUNDLED_SOURCE",
		DirectDependency: false,
		InclusionNote:    func() *string { s := "note"; return &s }(),
		ScopeStatus:      "IN_SCOPE",
		EvaluatedAt:      func() *dbtime.DBTime { t := dbtime.DBTime{Time: time.Now()}; return &t }(),
		EvaluatedBy:      func() *string { s := "user"; return &s }(),
		Modified:         true,
	}

	query := regexp.QuoteMeta("UPDATE project_usages SET oss_version_id = ?, usage_role = ?, direct_dependency = ?, inclusion_note = ?, scope_status = ?, evaluated_at = ?, evaluated_by = ?, modified = ?, modification_description = ?, supplier_type = ?, fork_origin_url = ? WHERE id = ?")
	mock.ExpectExec(query).
		WithArgs(u.OssVersionID, u.UsageRole, u.DirectDependency, u.InclusionNote, u.ScopeStatus, u.EvaluatedAt, u.EvaluatedBy, u.Modified, u.ModificationDescription, u.SupplierType, u.ForkOriginURL, u.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), u)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectUsageRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectUsageRepository{DB: db}

	id := uuid.NewString()
	query := regexp.QuoteMeta("DELETE FROM project_usages WHERE id = ?")
	mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectUsageRepository_UpdateScope(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectUsageRepository{DB: db}

	id := uuid.NewString()
	inclusion := "reason"
	evalBy := "user"
	now := dbtime.DBTime{Time: time.Now()}

	query := regexp.QuoteMeta("UPDATE project_usages SET scope_status = ?, inclusion_note = ?, evaluated_at = ?, evaluated_by = ? WHERE id = ?")
	mock.ExpectExec(query).WithArgs("OUT_SCOPE", &inclusion, now, &evalBy, id).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateScope(context.Background(), id, "OUT_SCOPE", &inclusion, now, &evalBy)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
