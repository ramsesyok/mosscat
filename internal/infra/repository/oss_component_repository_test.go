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
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

func TestOssComponentRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentRepository{DB: db}

	f := domrepo.OssComponentFilter{Name: "redis", Page: 1, Size: 10}

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM oss_components oc WHERE normalized_name LIKE ?")
	mock.ExpectQuery(countQuery).WithArgs("%redis%").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT oc.id, oc.name, oc.normalized_name, oc.homepage_url, oc.repository_url, oc.description, oc.primary_language, oc.default_usage_role, oc.deprecated, oc.created_at, oc.updated_at FROM oss_components oc WHERE normalized_name LIKE ? ORDER BY oc.created_at DESC LIMIT ? OFFSET ?")
	now := dbtime.DBTime{Time: time.Now()}
	mock.ExpectQuery(listQuery).WithArgs("%redis%", 10, 0).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "normalized_name", "homepage_url", "repository_url", "description", "primary_language", "default_usage_role", "deprecated", "created_at", "updated_at"}).AddRow(uuid.NewString(), "Redis", "redis", nil, nil, nil, nil, nil, false, now, now))

	res, total, err := repo.Search(context.Background(), f)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, res, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssComponentRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentRepository{DB: db}

	c := &model.OssComponent{
		ID:             uuid.NewString(),
		Name:           "Redis",
		NormalizedName: "redis",
		Deprecated:     false,
		CreatedAt:      dbtime.DBTime{Time: time.Now()},
		UpdatedAt:      dbtime.DBTime{Time: time.Now()},
	}

	query := regexp.QuoteMeta("INSERT INTO oss_components (id, name, normalized_name, homepage_url, repository_url, description, primary_language, default_usage_role, deprecated, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(query).WithArgs(c.ID, c.Name, c.NormalizedName, c.HomepageURL, c.RepositoryURL, c.Description, c.PrimaryLanguage, c.DefaultUsageRole, c.Deprecated, c.CreatedAt, c.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), c)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssComponentRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentRepository{DB: db}

	id := uuid.NewString()
	now := dbtime.DBTime{Time: time.Now()}
	query := regexp.QuoteMeta("SELECT id, name, normalized_name, homepage_url, repository_url, description, primary_language, default_usage_role, deprecated, created_at, updated_at FROM oss_components WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "normalized_name", "homepage_url", "repository_url", "description", "primary_language", "default_usage_role", "deprecated", "created_at", "updated_at"}).AddRow(id, "Redis", "redis", nil, nil, nil, nil, nil, false, now, now))

	res, err := repo.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, res.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssComponentRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentRepository{DB: db}
	c := &model.OssComponent{
		ID:               uuid.NewString(),
		Name:             "Redis",
		NormalizedName:   "redis",
		HomepageURL:      ptrString("http://hp"),
		RepositoryURL:    ptrString("http://repo"),
		Description:      ptrString("desc"),
		PrimaryLanguage:  ptrString("Go"),
		DefaultUsageRole: ptrString("RUNTIME_REQUIRED"),
		Deprecated:       true,
		UpdatedAt:        dbtime.DBTime{Time: time.Now()},
	}

	query := regexp.QuoteMeta("UPDATE oss_components SET name = ?, normalized_name = ?, homepage_url = ?, repository_url = ?, description = ?, primary_language = ?, default_usage_role = ?, deprecated = ?, updated_at = ? WHERE id = ?")
	mock.ExpectExec(query).WithArgs(c.Name, c.NormalizedName, c.HomepageURL, c.RepositoryURL, c.Description, c.PrimaryLanguage, c.DefaultUsageRole, c.Deprecated, c.UpdatedAt, c.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), c)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func ptrString(s string) *string { return &s }
