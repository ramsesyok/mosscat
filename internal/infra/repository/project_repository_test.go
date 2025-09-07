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

func TestProjectRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectRepository{DB: db}

	f := domrepo.ProjectFilter{Code: "PRJ", Page: 1, Size: 10}

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM projects WHERE project_code LIKE ?")
	mock.ExpectQuery(countQuery).WithArgs("%PRJ%").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT id, project_code, name, department, manager, delivery_date, description, created_at, updated_at, (SELECT COUNT(*) FROM project_usages u WHERE u.project_id = projects.id) FROM projects WHERE project_code LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?")
	now := dbtime.DBTime{Time: time.Now()}
	rows := sqlmock.NewRows([]string{"id", "project_code", "name", "department", "manager", "delivery_date", "description", "created_at", "updated_at", "count"}).
		AddRow(uuid.NewString(), "PRJ-1", "Proj", nil, nil, nil, nil, now, now, 0)
	mock.ExpectQuery(listQuery).WithArgs("%PRJ%", 10, 0).WillReturnRows(rows)

	res, total, err := repo.Search(context.Background(), f)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, res, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectRepository{DB: db}

	id := uuid.NewString()
	query := regexp.QuoteMeta("SELECT id, project_code, name, department, manager, delivery_date, description, created_at, updated_at, (SELECT COUNT(*) FROM project_usages u WHERE u.project_id = projects.id) FROM projects WHERE id = ?")
	now := dbtime.DBTime{Time: time.Now()}
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "project_code", "name", "department", "manager", "delivery_date", "description", "created_at", "updated_at", "count"}).
		AddRow(id, "PRJ-1", "Proj", nil, nil, nil, nil, now, now, 0))

	p, err := repo.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, p.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectRepository{DB: db}

	p := &model.Project{ID: uuid.NewString(), ProjectCode: "PRJ-1", Name: "Proj", CreatedAt: dbtime.DBTime{Time: time.Now()}, UpdatedAt: dbtime.DBTime{Time: time.Now()}}

	query := regexp.QuoteMeta("INSERT INTO projects (id, project_code, name, department, manager, delivery_date, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(query).WithArgs(p.ID, p.ProjectCode, p.Name, p.Department, p.Manager, p.DeliveryDate, p.Description, p.CreatedAt, p.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), p)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectRepository{DB: db}

	p := &model.Project{ID: uuid.NewString(), Name: "Proj", UpdatedAt: dbtime.DBTime{Time: time.Now()}}

	query := regexp.QuoteMeta("UPDATE projects SET name = ?, department = ?, manager = ?, delivery_date = ?, description = ?, updated_at = ? WHERE id = ?")
	mock.ExpectExec(query).WithArgs(p.Name, p.Department, p.Manager, p.DeliveryDate, p.Description, p.UpdatedAt, p.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), p)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ProjectRepository{DB: db}

	id := uuid.NewString()
	query := regexp.QuoteMeta("DELETE FROM projects WHERE id = ?")
	mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
