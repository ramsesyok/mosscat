package migration

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"

	"github.com/ramsesyok/mosscat/internal/domain/model"
	"github.com/ramsesyok/mosscat/internal/infra/repository"
	"github.com/ramsesyok/mosscat/migrations"
	"github.com/ramsesyok/mosscat/pkg/auth"
	"github.com/ramsesyok/mosscat/pkg/dbtime"
)

// Apply runs all up migrations using golang-migrate. If the schema is already up to date,
// it does nothing.
func Apply(db *sql.DB, dsn string) error {
	var (
		driver database.Driver
		err    error
	)

	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	} else {
		driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
	}
	if err != nil {
		return err
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	pass, err := ensureAdminUser(db)
	if err != nil {
		return err
	}
	if pass != "" {
		if err := savePasswordFile(pass); err != nil {
			return err
		}
	}
	return nil
}

func ensureAdminUser(db *sql.DB) (string, error) {
	repo := &repository.UserRepository{DB: db}
	ctx := context.Background()
	if _, err := repo.FindByUsername(ctx, "admin"); err == nil {
		return "", nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	pass := randomString(16)
	hash, err := auth.Hash(pass)
	if err != nil {
		return "", err
	}
	now := dbtime.DBTime{Time: time.Now()}
	u := &model.User{
		ID:           uuid.NewString(),
		Username:     "admin",
		PasswordHash: hash,
		Roles:        []string{"ADMIN"},
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := repo.Create(ctx, u); err != nil {
		return "", err
	}
	return pass, nil
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		for i := range b {
			b[i] = letters[int(time.Now().UnixNano())%len(letters)]
		}
		return string(b)
	}
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}

func savePasswordFile(pass string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Dir(exe)
	path := filepath.Join(dir, "admin.initial.password")
	return os.WriteFile(path, []byte(pass), 0600)
}
