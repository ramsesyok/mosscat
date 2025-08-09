package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

// --- stub implementations ---
type stubOssComponentRepo struct {
	createFn func(context.Context, *model.OssComponent) error
	getFn    func(context.Context, string) (*model.OssComponent, error)
	updateFn func(context.Context, *model.OssComponent) error
}

func (s *stubOssComponentRepo) Search(ctx context.Context, f domrepo.OssComponentFilter) ([]model.OssComponent, int, error) {
	return nil, 0, nil
}
func (s *stubOssComponentRepo) Create(ctx context.Context, c *model.OssComponent) error {
	if s.createFn != nil {
		return s.createFn(ctx, c)
	}
	return nil
}
func (s *stubOssComponentRepo) Get(ctx context.Context, id string) (*model.OssComponent, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, nil
}
func (s *stubOssComponentRepo) Update(ctx context.Context, c *model.OssComponent) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, c)
	}
	return nil
}

type stubOssComponentLayerRepo struct {
	replaceFn func(context.Context, string, []string) error
	listFn    func(context.Context, string) ([]string, error)
}

func (s *stubOssComponentLayerRepo) ListByOssID(ctx context.Context, id string) ([]string, error) {
	if s.listFn != nil {
		return s.listFn(ctx, id)
	}
	return nil, nil
}
func (s *stubOssComponentLayerRepo) Replace(ctx context.Context, id string, layers []string) error {
	if s.replaceFn != nil {
		return s.replaceFn(ctx, id, layers)
	}
	return nil
}

type stubOssComponentTagRepo struct {
	replaceFn func(context.Context, string, []string) error
	listFn    func(context.Context, string) ([]model.Tag, error)
}

func (s *stubOssComponentTagRepo) Replace(ctx context.Context, id string, tagIDs []string) error {
	if s.replaceFn != nil {
		return s.replaceFn(ctx, id, tagIDs)
	}
	return nil
}
func (s *stubOssComponentTagRepo) ListByOssID(ctx context.Context, id string) ([]model.Tag, error) {
	if s.listFn != nil {
		return s.listFn(ctx, id)
	}
	return nil, nil
}

type stubOssVersionRepo struct {
	searchFn func(context.Context, domrepo.OssVersionFilter) ([]model.OssVersion, int, error)
	createFn func(context.Context, *model.OssVersion) error
	deleteFn func(context.Context, string) error
	getFn    func(context.Context, string) (*model.OssVersion, error)
	updateFn func(context.Context, *model.OssVersion) error
}

func (s *stubOssVersionRepo) Search(ctx context.Context, f domrepo.OssVersionFilter) ([]model.OssVersion, int, error) {
	if s.searchFn != nil {
		return s.searchFn(ctx, f)
	}
	return nil, 0, nil
}
func (s *stubOssVersionRepo) Create(ctx context.Context, v *model.OssVersion) error {
	if s.createFn != nil {
		return s.createFn(ctx, v)
	}
	return nil
}
func (s *stubOssVersionRepo) Delete(ctx context.Context, id string) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}
func (s *stubOssVersionRepo) Get(ctx context.Context, id string) (*model.OssVersion, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, nil
}
func (s *stubOssVersionRepo) Update(ctx context.Context, v *model.OssVersion) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, v)
	}
	return nil
}

// --- tests ---
func TestToOssComponent_AllFields(t *testing.T) {
	now := dbtime.DBTime{Time: time.Now()}
	hp := "http://hp"
	repo := "http://repo"
	desc := "d"
	pl := "Go"
	role := "RUNTIME_REQUIRED"
	tagID := uuid.NewString()
	comp := model.OssComponent{ID: uuid.NewString(), Name: "Redis", NormalizedName: "redis", HomepageURL: &hp, RepositoryURL: &repo, Description: &desc, PrimaryLanguage: &pl, Layers: []string{"LIB"}, DefaultUsageRole: &role, Deprecated: false, CreatedAt: now, UpdatedAt: now, Tags: []model.Tag{{ID: tagID, Name: "db"}}}
	res := toOssComponent(comp)
	require.Equal(t, comp.ID, res.Id.String())
	require.NotNil(t, res.NormalizedName)
	require.Equal(t, "redis", *res.NormalizedName)
	require.NotNil(t, res.Layers)
	require.Equal(t, gen.Layer("LIB"), (*res.Layers)[0])
	require.NotNil(t, res.Tags)
	require.Equal(t, tagID, (*res.Tags)[0].Id.String())
	require.Equal(t, role, string(*res.DefaultUsageRole))
	require.Equal(t, &hp, res.HomepageUrl)
	require.Equal(t, &repo, res.RepositoryUrl)
	require.Equal(t, &desc, res.Description)
	require.Equal(t, &pl, res.PrimaryLanguage)
}

func TestToOssVersion_AllFields(t *testing.T) {
	now := dbtime.DBTime{Time: time.Now()}
	rel := now
	licRaw := "MIT"
	licConc := "MIT"
	purl := "pkg:"
	hash := "h"
	version := model.OssVersion{ID: uuid.NewString(), OssID: uuid.NewString(), Version: "1", ReleaseDate: &rel, LicenseExpressionRaw: &licRaw, LicenseConcluded: &licConc, Purl: &purl, CpeList: []string{"c"}, HashSha256: &hash, ReviewStatus: "draft", LastReviewedAt: &now, ScopeStatus: "IN_SCOPE", CreatedAt: now, UpdatedAt: now}
	res := toOssVersion(version)
	require.Equal(t, version.ID, res.Id.String())
	require.NotNil(t, res.ReleaseDate)
	require.NotNil(t, res.CpeList)
	require.Equal(t, "c", (*res.CpeList)[0])
	require.Equal(t, gen.ScopeStatus(version.ScopeStatus), res.ScopeStatus)
	require.Equal(t, gen.ReviewStatus(version.ReviewStatus), res.ReviewStatus)
	require.NotNil(t, res.HashSha256)
	require.Equal(t, hash, *res.HashSha256)
}

func TestCreateOssComponent(t *testing.T) {
	var created *model.OssComponent
	tagID := uuid.NewString()
	compRepo := &stubOssComponentRepo{createFn: func(ctx context.Context, c *model.OssComponent) error { created = c; return nil }}
	layerCalled := false
	layerRepo := &stubOssComponentLayerRepo{replaceFn: func(ctx context.Context, id string, layers []string) error {
		layerCalled = true
		require.Equal(t, created.ID, id)
		require.Equal(t, []string{"LIB"}, layers)
		return nil
	}}
	tagCalled := false
	tagRepo := &stubOssComponentTagRepo{replaceFn: func(ctx context.Context, id string, ids []string) error {
		tagCalled = true
		require.Equal(t, created.ID, id)
		require.Equal(t, []string{tagID}, ids)
		return nil
	}, listFn: func(ctx context.Context, id string) ([]model.Tag, error) {
		return []model.Tag{{ID: tagID, Name: "db"}}, nil
	}}
	h := &Handler{OssComponentRepo: compRepo, OssComponentLayerRepo: layerRepo, OssComponentTagRepo: tagRepo}
	e := setupEcho(h)
	body := `{"name":"Redis","layers":["LIB"],"tagIds":["` + tagID + `"]}`
	req := httptest.NewRequest(http.MethodPost, "/oss", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.True(t, layerCalled)
	require.True(t, tagCalled)
	require.NotNil(t, created)
	var res gen.OssComponent
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, created.ID, res.Id.String())
}

func TestListOssVersions(t *testing.T) {
	ossID := uuid.NewString()
	now := dbtime.DBTime{Time: time.Now()}
	repo := &stubOssVersionRepo{searchFn: func(ctx context.Context, f domrepo.OssVersionFilter) ([]model.OssVersion, int, error) {
		require.Equal(t, ossID, f.OssID)
		return []model.OssVersion{{ID: uuid.NewString(), OssID: ossID, Version: "1", ReviewStatus: "draft", ScopeStatus: "IN_SCOPE", CreatedAt: now, UpdatedAt: now}}, 1, nil
	}}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodGet, "/oss/"+ossID+"/versions", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var res gen.PagedResultOssVersion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotNil(t, res.Items)
	require.Len(t, *res.Items, 1)
}

func TestCreateOssVersion(t *testing.T) {
	ossID := uuid.NewString()
	var created *model.OssVersion
	repo := &stubOssVersionRepo{createFn: func(ctx context.Context, v *model.OssVersion) error { created = v; return nil }}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	body := `{"version":"1.0.0"}`
	req := httptest.NewRequest(http.MethodPost, "/oss/"+ossID+"/versions", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, created)
	require.Equal(t, ossID, created.OssID)
	var res gen.OssVersion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, created.ID, res.Id.String())
}

func TestCreateOssVersion_WithOptions(t *testing.T) {
	ossID := uuid.NewString()
	var created *model.OssVersion
	repo := &stubOssVersionRepo{createFn: func(ctx context.Context, v *model.OssVersion) error {
		created = v
		return nil
	}}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	now := dbtime.DBTime{Time: time.Now().UTC().Truncate(time.Second)}
	body := `{"version":"1.0.0","releaseDate":"` + now.Format("2006-01-02") + `","purl":"pkg:","cpeList":["c"],"modified":true,"supplierType":"ORIGIN"}`
	req := httptest.NewRequest(http.MethodPost, "/oss/"+ossID+"/versions", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, created)
	require.Equal(t, ossID, created.OssID)
	require.Equal(t, "pkg:", *created.Purl)
	require.Equal(t, []string{"c"}, created.CpeList)
	require.NotNil(t, created.ReleaseDate)
}

func TestCreateOssVersion_InvalidBody(t *testing.T) {
	ossID := uuid.NewString()
	repo := &stubOssVersionRepo{}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodPost, "/oss/"+ossID+"/versions", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteOssVersion(t *testing.T) {
	ossID := uuid.NewString()
	vid := uuid.NewString()
	called := false
	repo := &stubOssVersionRepo{deleteFn: func(ctx context.Context, id string) error { called = true; require.Equal(t, vid, id); return nil }}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodDelete, "/oss/"+ossID+"/versions/"+vid, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.True(t, called)
}

func TestUpdateOssVersion_NotFound(t *testing.T) {
	ossID := uuid.NewString()
	vid := uuid.NewString()
	repo := &stubOssVersionRepo{getFn: func(ctx context.Context, id string) (*model.OssVersion, error) {
		return nil, sql.ErrNoRows
	}}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	body := `{"reviewStatus":"approved"}`
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+ossID+"/versions/"+vid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestOssComponentStubMethods(t *testing.T) {
	h := &Handler{}
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())
	require.NoError(t, h.DeprecateOssComponent(ctx, uuid.New()))
}

func TestUpdateOssComponent_NotFound(t *testing.T) {
	oid := uuid.NewString()
	repo := &stubOssComponentRepo{getFn: func(ctx context.Context, id string) (*model.OssComponent, error) {
		return nil, sql.ErrNoRows
	}}
	h := &Handler{OssComponentRepo: repo}
	e := setupEcho(h)
	body := `{"description":"x"}`
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+oid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateOssComponent(t *testing.T) {
	oid := uuid.NewString()
	now := dbtime.DBTime{Time: time.Now()}
	existing := model.OssComponent{ID: oid, Name: "Redis", NormalizedName: "redis", Deprecated: false, CreatedAt: now, UpdatedAt: now}
	var updated *model.OssComponent
	repo := &stubOssComponentRepo{getFn: func(ctx context.Context, id string) (*model.OssComponent, error) {
		require.Equal(t, oid, id)
		c := existing
		return &c, nil
	}, updateFn: func(ctx context.Context, c *model.OssComponent) error {
		updated = c
		return nil
	}}
	h := &Handler{OssComponentRepo: repo, OssComponentLayerRepo: &stubOssComponentLayerRepo{listFn: func(ctx context.Context, id string) ([]string, error) {
		return nil, nil
	}}, OssComponentTagRepo: &stubOssComponentTagRepo{listFn: func(ctx context.Context, id string) ([]model.Tag, error) { return nil, nil }}}
	e := setupEcho(h)
	body := `{"description":"patched"}`
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+oid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, updated)
	require.Equal(t, "patched", *updated.Description)
}

func TestUpdateOssComponent_WithFields(t *testing.T) {
	oid := uuid.NewString()
	now := dbtime.DBTime{Time: time.Now()}
	existing := model.OssComponent{ID: oid, Name: "Old", NormalizedName: "old", Deprecated: false, CreatedAt: now, UpdatedAt: now}
	tagID := uuid.NewString()
	var updated *model.OssComponent
	var replacedLayers []string
	var replacedTags []string
	repo := &stubOssComponentRepo{getFn: func(ctx context.Context, id string) (*model.OssComponent, error) {
		c := existing
		return &c, nil
	}, updateFn: func(ctx context.Context, c *model.OssComponent) error {
		updated = c
		return nil
	}}
	layerRepo := &stubOssComponentLayerRepo{replaceFn: func(ctx context.Context, id string, layers []string) error {
		replacedLayers = layers
		return nil
	}, listFn: func(ctx context.Context, id string) ([]string, error) { return replacedLayers, nil }}
	tagRepo := &stubOssComponentTagRepo{replaceFn: func(ctx context.Context, id string, ids []string) error {
		replacedTags = ids
		return nil
	}, listFn: func(ctx context.Context, id string) ([]model.Tag, error) {
		return []model.Tag{{ID: tagID, Name: "db"}}, nil
	}}
	h := &Handler{OssComponentRepo: repo, OssComponentLayerRepo: layerRepo, OssComponentTagRepo: tagRepo}
	e := setupEcho(h)
	body := `{"name":"New","homepageUrl":"http://hp","repositoryUrl":"http://repo","description":"d","primaryLanguage":"Go","layers":["LIB"],"defaultUsageRole":"RUNTIME_REQUIRED","deprecated":true,"tagIds":["` + tagID + `"]}`
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+oid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, updated)
	require.Equal(t, "New", updated.Name)
	require.Equal(t, "new", updated.NormalizedName)
	require.Equal(t, "http://hp", *updated.HomepageURL)
	require.Equal(t, "http://repo", *updated.RepositoryURL)
	require.Equal(t, []string{"LIB"}, replacedLayers)
	require.Equal(t, []string{tagID}, replacedTags)
}

func TestUpdateOssVersion(t *testing.T) {
	ossID := uuid.NewString()
	vid := uuid.NewString()
	now := dbtime.DBTime{Time: time.Now()}
	existing := model.OssVersion{ID: vid, OssID: ossID, Version: "1", ReviewStatus: "draft", ScopeStatus: "IN_SCOPE", CreatedAt: now, UpdatedAt: now}
	var updated *model.OssVersion
	repo := &stubOssVersionRepo{getFn: func(ctx context.Context, id string) (*model.OssVersion, error) {
		require.Equal(t, vid, id)
		v := existing
		return &v, nil
	}, updateFn: func(ctx context.Context, v *model.OssVersion) error { updated = v; return nil }}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	body := `{"reviewStatus":"approved"}`
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+ossID+"/versions/"+vid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, updated)
	require.Equal(t, "approved", updated.ReviewStatus)
}

func TestUpdateOssVersion_WithFields(t *testing.T) {
	ossID := uuid.NewString()
	vid := uuid.NewString()
	now := dbtime.DBTime{Time: time.Now()}
	existing := model.OssVersion{ID: vid, OssID: ossID, Version: "1", ReviewStatus: "draft", ScopeStatus: "IN_SCOPE", CreatedAt: now, UpdatedAt: now}
	var updated *model.OssVersion
	repo := &stubOssVersionRepo{getFn: func(ctx context.Context, id string) (*model.OssVersion, error) {
		v := existing
		return &v, nil
	}, updateFn: func(ctx context.Context, v *model.OssVersion) error {
		updated = v
		return nil
	}}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	body := `{"licenseExpressionRaw":"MIT","scopeStatus":"OUT_SCOPE"}`
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+ossID+"/versions/"+vid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, updated)
	require.Equal(t, "MIT", *updated.LicenseExpressionRaw)
	require.Equal(t, "OUT_SCOPE", updated.ScopeStatus)
}
