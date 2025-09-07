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

func TestAuditLogRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &AuditLogRepository{DB: db}

	f := domrepo.AuditLogFilter{EntityType: func() *string { s := "PROJECT"; return &s }()}

	query := regexp.QuoteMeta("SELECT id, entity_type, entity_id, action, user_name, summary, created_at FROM audit_logs WHERE entity_type = ? ORDER BY created_at DESC")
	now := dbtime.DBTime{Time: time.Now()}
	rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "user_name", "summary", "created_at"}).
		AddRow(uuid.NewString(), "PROJECT", "1", "CREATE", "user", "created", now)
	mock.ExpectQuery(query).WithArgs("PROJECT").WillReturnRows(rows)

	logs, err := repo.Search(context.Background(), f)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditLogRepository_Search_AllFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &AuditLogRepository{DB: db}

	et := "PROJECT"
	eid := "1"
	from := dbtime.DBTime{Time: time.Now().Add(-time.Hour)}
	to := dbtime.DBTime{Time: time.Now()}
	f := domrepo.AuditLogFilter{EntityType: &et, EntityID: &eid, From: &from, To: &to}

	query := regexp.QuoteMeta("SELECT id, entity_type, entity_id, action, user_name, summary, created_at FROM audit_logs WHERE entity_type = ? AND entity_id = ? AND created_at >= ? AND created_at <= ? ORDER BY created_at DESC")
	rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "user_name", "summary", "created_at"}).AddRow(uuid.NewString(), et, eid, "CREATE", "user", nil, to)
	mock.ExpectQuery(query).WithArgs(et, eid, from, to).WillReturnRows(rows)

	logs, err := repo.Search(context.Background(), f)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditLogRepository_Search_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &AuditLogRepository{DB: db}
	et := "PROJECT"
	f := domrepo.AuditLogFilter{EntityType: &et}

	query := regexp.QuoteMeta("SELECT id, entity_type, entity_id, action, user_name, summary, created_at FROM audit_logs WHERE entity_type = ? ORDER BY created_at DESC")
	mock.ExpectQuery(query).WithArgs(et).WillReturnError(errors.New("fail"))

	_, err = repo.Search(context.Background(), f)
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditLogRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &AuditLogRepository{DB: db}

	l := &model.AuditLog{
		ID:         uuid.NewString(),
		EntityType: "PROJECT",
		EntityID:   "1",
		Action:     "CREATE",
		UserName:   "user",
		CreatedAt:  dbtime.DBTime{Time: time.Now()},
	}

	query := regexp.QuoteMeta("INSERT INTO audit_logs (id, entity_type, entity_id, action, user_name, summary, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(query).
		WithArgs(l.ID, l.EntityType, l.EntityID, l.Action, l.UserName, l.Summary, l.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), l)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
