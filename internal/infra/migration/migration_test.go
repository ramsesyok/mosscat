package migration

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestApply_SQLite(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test?mode=memory&cache=shared")
	require.NoError(t, err)
	defer db.Close()

	err = Apply(db, "file:test?mode=memory&cache=shared")
	require.NoError(t, err)

	var name string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='tags'`).Scan(&name)
	require.NoError(t, err)
	require.Equal(t, "tags", name)

	var cnt int
	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = 'admin'`).Scan(&cnt)
	require.NoError(t, err)
	require.Equal(t, 1, cnt)

	exe, err := os.Executable()
	require.NoError(t, err)
	p := filepath.Join(filepath.Dir(exe), "admin.initial.password")
	data, err := os.ReadFile(p)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	os.Remove(p)
}
