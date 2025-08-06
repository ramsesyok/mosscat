package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDB_WithinTx_Commit(t *testing.T) {
	rawDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer rawDB.Close()

	d := &DB{DB: rawDB}
	mock.ExpectBegin()
	mock.ExpectCommit()

	err = d.WithinTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return nil
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_WithinTx_Rollback(t *testing.T) {
	rawDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer rawDB.Close()

	d := &DB{DB: rawDB}
	mock.ExpectBegin()
	mock.ExpectRollback()

	err = d.WithinTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return errors.New("fail")
	})
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_WithinTx_BeginError(t *testing.T) {
	rawDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer rawDB.Close()

	d := &DB{DB: rawDB}
	beginErr := errors.New("begin fail")
	mock.ExpectBegin().WillReturnError(beginErr)

	err = d.WithinTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return nil
	})
	require.ErrorIs(t, err, beginErr)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOpen(t *testing.T) {
	t.Run("SQLite", func(t *testing.T) {
		db, err := Open("file:test?mode=memory&cache=shared")
		require.NoError(t, err)
		require.NoError(t, db.Close())
	})

	t.Run("Postgres", func(t *testing.T) {
		db, err := Open("postgres://user:pass@localhost/dbname?sslmode=disable")
		require.NoError(t, err)
		require.NoError(t, db.Close())
	})
}
