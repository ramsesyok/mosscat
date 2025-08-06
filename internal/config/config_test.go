package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoad_Default(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)
	require.Equal(t, "0.0.0.0", cfg.Server.Host)
	require.Equal(t, "8080", cfg.Server.Port)
	require.Equal(t, []string{"*"}, cfg.Server.AllowedOrigins)
	require.Equal(t, "file:mosscat.db?mode=memory&cache=shared", cfg.DB.DSN)
}

func TestLoad_FromFile(t *testing.T) {
	data := []byte("server:\n  host: 127.0.0.1\n  port: '9090'\n  allowed_origins:\n    - http://example.com\n    - http://foo.com\ndb:\n  dsn: postgres://user:pass@localhost/db")
	f, err := os.CreateTemp(t.TempDir(), "cfg*.yaml")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	_, err = f.Write(data)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	cfg, err := Load(f.Name())
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", cfg.Server.Host)
	require.Equal(t, "9090", cfg.Server.Port)
	require.Equal(t, []string{"http://example.com", "http://foo.com"}, cfg.Server.AllowedOrigins)
	require.Equal(t, "postgres://user:pass@localhost/db", cfg.DB.DSN)
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/no/such/file.yaml")
	require.Error(t, err)
}
