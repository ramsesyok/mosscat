package handler

import (
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	infrarepo "github.com/ramsesyok/mosscat/internal/infra/repository"
)

func TestSearchAuditLogs_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.AuditLogRepository{DB: db}
	h := &Handler{AuditRepo: repo}
	e := setupEcho(h)

	et := "PROJECT"
	eid := "p1"
	from := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	to := from.Add(time.Hour)

	query := regexp.QuoteMeta("SELECT id, entity_type, entity_id, action, user_name, summary, created_at FROM audit_logs WHERE entity_type = ? AND entity_id = ? AND created_at >= ? AND created_at <= ? ORDER BY created_at DESC")
	rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "user_name", "summary", "created_at"}).
		AddRow(uuid.NewString(), et, eid, "CREATE", "user1", nil, from).
		AddRow(uuid.NewString(), et, eid, "UPDATE", "user2", "note", to)
	mock.ExpectQuery(query).WithArgs(et, eid, from, to).WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/audit?entityType="+et+"&entityId="+eid+"&from="+from.Format(time.RFC3339)+"&to="+to.Format(time.RFC3339), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res struct {
		Items []map[string]any `json:"items"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Len(t, res.Items, 2)
	_, ok := res.Items[0]["summary"]
	require.False(t, ok)
	require.Equal(t, "note", res.Items[1]["summary"].(string))
}

func TestSearchAuditLogs_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.AuditLogRepository{DB: db}
	h := &Handler{AuditRepo: repo}
	e := setupEcho(h)

	et := "PROJECT"
	query := regexp.QuoteMeta("SELECT id, entity_type, entity_id, action, user_name, summary, created_at FROM audit_logs WHERE entity_type = ? ORDER BY created_at DESC")
	mock.ExpectQuery(query).WithArgs(et).WillReturnError(driver.ErrBadConn)

	req := httptest.NewRequest(http.MethodGet, "/audit?entityType="+et, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}
