package repository

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

func setupSQLiteDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", "file:test?mode=memory&cache=shared&_loc=auto")
	require.NoError(t, err)
	_, file, _, _ := runtime.Caller(0)
	// Apply migration 0001
	path1 := filepath.Join(filepath.Dir(file), "..", "..", "..", "migrations", "0001_create_tables.up.sql")
	sqlBytes1, err := os.ReadFile(path1)
	require.NoError(t, err)
	sqlStr1 := strings.ReplaceAll(string(sqlBytes1), "TIMESTAMPTZ", "TIMESTAMP")
	_, err = db.Exec(sqlStr1)
	require.NoError(t, err)

	// Apply migration 0003 (move modification fields)
	path3 := filepath.Join(filepath.Dir(file), "..", "..", "..", "migrations", "0003_move_modification_fields.up.sql")
	sqlBytes3, err := os.ReadFile(path3)
	require.NoError(t, err)
	_, err = db.Exec(string(sqlBytes3))
	require.NoError(t, err)
	return db
}

func TestRepositories_SQLite(t *testing.T) {
	ctx := context.Background()

	t.Run("TagRepository", func(t *testing.T) {
		db := setupSQLiteDB(t)
		defer db.Close()

		repo := &TagRepository{DB: db}
		now := dbtime.DBTime{Time: time.Now()}
		tag := &model.Tag{ID: uuid.NewString(), Name: "db", CreatedAt: &now}
		require.NoError(t, repo.Create(ctx, tag))
		tags, err := repo.List(ctx)
		require.NoError(t, err)
		require.Len(t, tags, 1)
		require.Equal(t, tag.ID, tags[0].ID)
		require.NoError(t, repo.Delete(ctx, tag.ID))
		tags, err = repo.List(ctx)
		require.NoError(t, err)
		require.Len(t, tags, 0)
	})

	t.Run("OssComponentRepositories", func(t *testing.T) {
		db := setupSQLiteDB(t)
		defer db.Close()
		compRepo := &OssComponentRepository{DB: db}
		layerRepo := &OssComponentLayerRepository{DB: db}
		tagRepo := &TagRepository{DB: db}
		compTagRepo := &OssComponentTagRepository{DB: db}

		now := dbtime.DBTime{Time: time.Now()}
		tag := &model.Tag{ID: uuid.NewString(), Name: "db", CreatedAt: &now}
		require.NoError(t, tagRepo.Create(ctx, tag))

		comp := &model.OssComponent{
			ID:             uuid.NewString(),
			Name:           "Redis",
			NormalizedName: "redis",
			Deprecated:     false,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		require.NoError(t, compRepo.Create(ctx, comp))
		require.NoError(t, layerRepo.Replace(ctx, comp.ID, []string{"LIB"}))
		require.NoError(t, compTagRepo.Replace(ctx, comp.ID, []string{tag.ID}))

		res, total, err := compRepo.Search(ctx, domrepo.OssComponentFilter{Name: "red", Page: 1, Size: 10})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, comp.ID, res[0].ID)

		res, total, err = compRepo.Search(ctx, domrepo.OssComponentFilter{Layers: []string{"LIB"}, Page: 1, Size: 10})
		require.NoError(t, err)
		require.Equal(t, 1, total)

		res, total, err = compRepo.Search(ctx, domrepo.OssComponentFilter{Tag: "db", Page: 1, Size: 10})
		require.NoError(t, err)
		require.Equal(t, 1, total)

		layers, err := layerRepo.ListByOssID(ctx, comp.ID)
		require.NoError(t, err)
		require.Equal(t, []string{"LIB"}, layers)

		tags, err := compTagRepo.ListByOssID(ctx, comp.ID)
		require.NoError(t, err)
		require.Len(t, tags, 1)
		require.Equal(t, tag.ID, tags[0].ID)
	})

	t.Run("OssVersionRepository", func(t *testing.T) {
		db := setupSQLiteDB(t)
		defer db.Close()
		compRepo := &OssComponentRepository{DB: db}
		verRepo := &OssVersionRepository{DB: db}

		now := dbtime.DBTime{Time: time.Now()}
		comp := &model.OssComponent{ID: uuid.NewString(), Name: "Redis", NormalizedName: "redis", CreatedAt: now, UpdatedAt: now}
		require.NoError(t, compRepo.Create(ctx, comp))

		ver := &model.OssVersion{
			ID:           uuid.NewString(),
			OssID:        comp.ID,
			Version:      "1.0.0",
			ReviewStatus: "draft",
			ScopeStatus:  "IN_SCOPE",
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		require.NoError(t, verRepo.Create(ctx, ver))

		got, err := verRepo.Get(ctx, ver.ID)
		require.NoError(t, err)
		require.Equal(t, ver.ID, got.ID)

		ver.ReviewStatus = "verified"
		ver.UpdatedAt = dbtime.DBTime{Time: time.Now()}
		require.NoError(t, verRepo.Update(ctx, ver))

		res, total, err := verRepo.Search(ctx, domrepo.OssVersionFilter{OssID: comp.ID, ReviewStatus: "verified", Page: 1, Size: 10})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, ver.ID, res[0].ID)

		require.NoError(t, verRepo.Delete(ctx, ver.ID))
		res, total, err = verRepo.Search(ctx, domrepo.OssVersionFilter{OssID: comp.ID, Page: 1, Size: 10})
		require.NoError(t, err)
		require.Equal(t, 0, total)
	})

	t.Run("ProjectRepository", func(t *testing.T) {
		db := setupSQLiteDB(t)
		defer db.Close()
		repo := &ProjectRepository{DB: db}

		now := dbtime.DBTime{Time: time.Now()}
		proj := &model.Project{ID: uuid.NewString(), ProjectCode: "P1", Name: "Proj", CreatedAt: now, UpdatedAt: now}
		require.NoError(t, repo.Create(ctx, proj))

		p, err := repo.Get(ctx, proj.ID)
		require.NoError(t, err)
		require.Equal(t, proj.ID, p.ID)

		proj.Name = "Updated"
		proj.UpdatedAt = dbtime.DBTime{Time: time.Now()}
		require.NoError(t, repo.Update(ctx, proj))

		res, total, err := repo.Search(ctx, domrepo.ProjectFilter{Name: "Upd", Page: 1, Size: 10})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Len(t, res, 1)

		require.NoError(t, repo.Delete(ctx, proj.ID))
	})

	t.Run("ProjectUsageRepository", func(t *testing.T) {
		db := setupSQLiteDB(t)
		defer db.Close()
		compRepo := &OssComponentRepository{DB: db}
		verRepo := &OssVersionRepository{DB: db}
		projRepo := &ProjectRepository{DB: db}
		usageRepo := &ProjectUsageRepository{DB: db}

		now := dbtime.DBTime{Time: time.Now()}
		comp := &model.OssComponent{ID: uuid.NewString(), Name: "Redis", NormalizedName: "redis", CreatedAt: now, UpdatedAt: now}
		require.NoError(t, compRepo.Create(ctx, comp))
		ver := &model.OssVersion{ID: uuid.NewString(), OssID: comp.ID, Version: "1.0.0", ReviewStatus: "draft", ScopeStatus: "IN_SCOPE", CreatedAt: now, UpdatedAt: now}
		require.NoError(t, verRepo.Create(ctx, ver))
		proj := &model.Project{ID: uuid.NewString(), ProjectCode: "P1", Name: "Proj", CreatedAt: now, UpdatedAt: now}
		require.NoError(t, projRepo.Create(ctx, proj))

		usage := &model.ProjectUsage{
			ID:               uuid.NewString(),
			ProjectID:        proj.ID,
			OssID:            comp.ID,
			OssVersionID:     ver.ID,
			UsageRole:        "RUNTIME_REQUIRED",
			ScopeStatus:      "IN_SCOPE",
			DirectDependency: true,
			AddedAt:          now,
			Modified:         false,
		}
		require.NoError(t, usageRepo.Create(ctx, usage))

		res, total, err := usageRepo.Search(ctx, domrepo.ProjectUsageFilter{ProjectID: proj.ID, Page: 1, Size: 10})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Equal(t, usage.ID, res[0].ID)

		usage.UsageRole = "DEV_TOOL"
		require.NoError(t, usageRepo.Update(ctx, usage))

		note := "out"
		now2 := dbtime.DBTime{Time: time.Now()}
		user := "tester"
		require.NoError(t, usageRepo.UpdateScope(ctx, usage.ID, "OUT_SCOPE", &note, now2, &user))

		require.NoError(t, usageRepo.Delete(ctx, usage.ID))
	})

	t.Run("ScopePolicyRepository", func(t *testing.T) {
		db := setupSQLiteDB(t)
		defer db.Close()
		repo := &ScopePolicyRepository{DB: db}
		now := dbtime.DBTime{Time: time.Now()}
		policy := &model.ScopePolicy{ID: uuid.NewString(), RuntimeRequiredDefaultInScope: true, ServerEnvIncluded: false, AutoMarkForksInScope: true, UpdatedAt: now, UpdatedBy: "user"}
		require.NoError(t, repo.Update(ctx, policy))
		p, err := repo.Get(ctx)
		require.NoError(t, err)
		require.Equal(t, policy.ID, p.ID)
	})

	t.Run("AuditLogRepository", func(t *testing.T) {
		db := setupSQLiteDB(t)
		defer db.Close()
		repo := &AuditLogRepository{DB: db}
		now := dbtime.DBTime{Time: time.Now()}
		l := &model.AuditLog{ID: uuid.NewString(), EntityType: "PROJECT", EntityID: "1", Action: "CREATE", UserName: "user", CreatedAt: now}
		require.NoError(t, repo.Create(ctx, l))
		et := "PROJECT"
		logs, err := repo.Search(ctx, domrepo.AuditLogFilter{EntityType: &et})
		require.NoError(t, err)
		require.Len(t, logs, 1)
	})
}
