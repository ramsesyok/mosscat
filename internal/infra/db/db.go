package db

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// DB は sql.DB をラップしトランザクション管理を提供する。
type DB struct {
	*sql.DB
}

// Open は DSN に基づき DB 接続を確立する。
// DSN が "postgres://" または "postgresql://" で始まる場合は postgres ドライバを使用する。
// それ以外は SQLite とみなす。
func Open(dsn string) (*DB, error) {
	driver := "sqlite3"
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		driver = "postgres"
	} else if !strings.Contains(dsn, "_loc=") {
		if strings.Contains(dsn, "?") {
			dsn += "&_loc=auto"
		} else {
			dsn += "?_loc=auto"
		}
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

// WithinTx はトランザクション内で fn を実行する。
func (d *DB) WithinTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}
