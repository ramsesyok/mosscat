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

func TestOssVersionRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssVersionRepository{DB: db}

	f := domrepo.OssVersionFilter{OssID: uuid.NewString(), Page: 1, Size: 10}

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM oss_versions WHERE oss_id = ?")
	mock.ExpectQuery(countQuery).WithArgs(f.OssID).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, review_status, last_reviewed_at, scope_status, created_at, updated_at FROM oss_versions WHERE oss_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?")
	now := dbtime.DBTime{Time: time.Now()}
	rows := sqlmock.NewRows([]string{"id", "oss_id", "version", "release_date", "license_expression_raw", "license_concluded", "purl", "cpe_list", "hash_sha256", "review_status", "last_reviewed_at", "scope_status", "created_at", "updated_at"}).
		AddRow(uuid.NewString(), f.OssID, "1.0.0", now, nil, nil, nil, pq.StringArray{"cpe:/a"}, nil, "draft", nil, "IN_SCOPE", now, now)
	mock.ExpectQuery(listQuery).WithArgs(f.OssID, 10, 0).WillReturnRows(rows)

	res, total, err := repo.Search(context.Background(), f)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, res, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssVersionRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssVersionRepository{DB: db}

	v := &model.OssVersion{
		ID:           uuid.NewString(),
		OssID:        uuid.NewString(),
		Version:      "1.0.0",
		ReviewStatus: "draft",
		ScopeStatus:  "IN_SCOPE",
		CreatedAt:    dbtime.DBTime{Time: time.Now()},
		UpdatedAt:    dbtime.DBTime{Time: time.Now()},
	}

	query := regexp.QuoteMeta("INSERT INTO oss_versions (id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, review_status, last_reviewed_at, scope_status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(query).
		WithArgs(v.ID, v.OssID, v.Version, v.ReleaseDate, v.LicenseExpressionRaw, v.LicenseConcluded, v.Purl, sqlmock.AnyArg(), v.HashSha256, v.ReviewStatus, v.LastReviewedAt, v.ScopeStatus, v.CreatedAt, v.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), v)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssVersionRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssVersionRepository{DB: db}

	v := &model.OssVersion{
		ID:           uuid.NewString(),
		ReleaseDate:  func() *dbtime.DBTime { t := dbtime.DBTime{Time: time.Now()}; return &t }(),
		ReviewStatus: "verified",
		ScopeStatus:  "IN_SCOPE",
		UpdatedAt:    dbtime.DBTime{Time: time.Now()},
	}

	query := regexp.QuoteMeta("UPDATE oss_versions SET release_date = ?, license_expression_raw = ?, license_concluded = ?, purl = ?, cpe_list = ?, hash_sha256 = ?, review_status = ?, last_reviewed_at = ?, scope_status = ?, updated_at = ? WHERE id = ?")
	mock.ExpectExec(query).
		WithArgs(v.ReleaseDate, v.LicenseExpressionRaw, v.LicenseConcluded, v.Purl, sqlmock.AnyArg(), v.HashSha256, v.ReviewStatus, v.LastReviewedAt, v.ScopeStatus, v.UpdatedAt, v.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), v)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssVersionRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssVersionRepository{DB: db}

	id := uuid.NewString()

	query := regexp.QuoteMeta("DELETE FROM oss_versions WHERE id = ?")
	mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
