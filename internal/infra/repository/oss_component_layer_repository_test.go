package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOssComponentLayerRepository_ListByOssID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentLayerRepository{DB: db}

	ossID := uuid.NewString()
	query := regexp.QuoteMeta(`SELECT layer FROM oss_component_layers WHERE oss_id = ? ORDER BY layer`)
	rows := sqlmock.NewRows([]string{"layer"}).AddRow("LIB").AddRow("DB")
	mock.ExpectQuery(query).WithArgs(ossID).WillReturnRows(rows)

	layers, err := repo.ListByOssID(context.Background(), ossID)
	require.NoError(t, err)
	require.Equal(t, []string{"LIB", "DB"}, layers)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssComponentLayerRepository_Replace(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentLayerRepository{DB: db}

	ossID := uuid.NewString()
	layers := []string{"LIB", "DB"}

	mock.ExpectBegin()
	delQuery := regexp.QuoteMeta(`DELETE FROM oss_component_layers WHERE oss_id = ?`)
	mock.ExpectExec(delQuery).WithArgs(ossID).WillReturnResult(sqlmock.NewResult(1, 1))
	insQuery := regexp.QuoteMeta(`INSERT INTO oss_component_layers (oss_id, layer) VALUES (?, ?)`)
	for _, l := range layers {
		mock.ExpectExec(insQuery).WithArgs(ossID, l).WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()

	err = repo.Replace(context.Background(), ossID, layers)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssComponentLayerRepository_Replace_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentLayerRepository{DB: db}
	ossID := uuid.NewString()

	mock.ExpectBegin()
	delQuery := regexp.QuoteMeta(`DELETE FROM oss_component_layers WHERE oss_id = ?`)
	mock.ExpectExec(delQuery).WithArgs(ossID).WillReturnError(errors.New("del"))
	mock.ExpectRollback()

	err = repo.Replace(context.Background(), ossID, []string{"LIB"})
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssComponentLayerRepository_Replace_InsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentLayerRepository{DB: db}
	ossID := uuid.NewString()

	mock.ExpectBegin()
	delQuery := regexp.QuoteMeta(`DELETE FROM oss_component_layers WHERE oss_id = ?`)
	mock.ExpectExec(delQuery).WithArgs(ossID).WillReturnResult(sqlmock.NewResult(1, 1))
	insQuery := regexp.QuoteMeta(`INSERT INTO oss_component_layers (oss_id, layer) VALUES (?, ?)`)
	mock.ExpectExec(insQuery).WithArgs(ossID, "LIB").WillReturnError(errors.New("ins"))
	mock.ExpectRollback()

	err = repo.Replace(context.Background(), ossID, []string{"LIB"})
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
