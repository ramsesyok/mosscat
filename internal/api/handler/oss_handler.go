package handler

// oss_handler.go - /oss に関するハンドラ処理

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

func toOssComponent(m model.OssComponent) gen.OssComponent {
	uid := uuid.MustParse(m.ID)
	res := gen.OssComponent{
		Id:         uid,
		Name:       m.Name,
		Deprecated: m.Deprecated,
		CreatedAt:  m.CreatedAt.TimeValue(),
		UpdatedAt:  m.UpdatedAt.TimeValue(),
	}
	if m.NormalizedName != "" {
		res.NormalizedName = &m.NormalizedName
	}
	if m.HomepageURL != nil {
		res.HomepageUrl = m.HomepageURL
	}
	if m.RepositoryURL != nil {
		res.RepositoryUrl = m.RepositoryURL
	}
	if m.Description != nil {
		res.Description = m.Description
	}
	if m.PrimaryLanguage != nil {
		res.PrimaryLanguage = m.PrimaryLanguage
	}
	if len(m.Layers) > 0 {
		layers := make([]gen.Layer, len(m.Layers))
		for i, l := range m.Layers {
			layers[i] = gen.Layer(l)
		}
		res.Layers = &layers
	}
	if m.DefaultUsageRole != nil {
		val := gen.UsageRole(*m.DefaultUsageRole)
		res.DefaultUsageRole = &val
	}
	if len(m.Tags) > 0 {
		tags := make([]gen.Tag, len(m.Tags))
		for i, t := range m.Tags {
			tags[i] = gen.Tag{Id: uuid.MustParse(t.ID), Name: t.Name}
		}
		res.Tags = &tags
	}
	return res
}

func toOssVersion(m model.OssVersion) gen.OssVersion {
	uid := uuid.MustParse(m.ID)
	ossUID := uuid.MustParse(m.OssID)
	res := gen.OssVersion{
		Id:           uid,
		OssId:        ossUID,
		Version:      m.Version,
		Modified:     m.Modified,
		ReviewStatus: gen.ReviewStatus(m.ReviewStatus),
		ScopeStatus:  gen.ScopeStatus(m.ScopeStatus),
		CreatedAt:    m.CreatedAt.TimeValue(),
		UpdatedAt:    m.UpdatedAt.TimeValue(),
	}
	if m.ReleaseDate != nil {
		res.ReleaseDate = &openapi_types.Date{Time: m.ReleaseDate.TimeValue()}
	}
	if m.LicenseExpressionRaw != nil {
		res.LicenseExpressionRaw = m.LicenseExpressionRaw
	}
	if m.LicenseConcluded != nil {
		res.LicenseConcluded = m.LicenseConcluded
	}
	if m.Purl != nil {
		res.Purl = m.Purl
	}
	if len(m.CpeList) > 0 {
		list := make([]string, len(m.CpeList))
		copy(list, m.CpeList)
		res.CpeList = &list
	}
	if m.HashSha256 != nil {
		res.HashSha256 = m.HashSha256
	}
	if m.ModificationDescription != nil {
		res.ModificationDescription = m.ModificationDescription
	}
	if m.LastReviewedAt != nil {
		t := m.LastReviewedAt.TimeValue()
		res.LastReviewedAt = &t
	}
	if m.SupplierType != nil {
		val := gen.SupplierType(*m.SupplierType)
		res.SupplierType = &val
	}
	if m.ForkOriginURL != nil {
		res.ForkOriginUrl = m.ForkOriginURL
	}
	return res
}

// OSSコンポーネント一覧取得
// (GET /oss)
func (h *Handler) ListOssComponents(ctx echo.Context, params gen.ListOssComponentsParams) error {
	page := 1
	if params.Page != nil {
		page = int(*params.Page)
	}
	size := 50
	if params.Size != nil {
		size = int(*params.Size)
	}
	f := domrepo.OssComponentFilter{Page: page, Size: size}
	if params.Name != nil {
		f.Name = *params.Name
	}
	if params.Layers != nil && *params.Layers != "" {
		f.Layers = strings.Split(*params.Layers, ",")
	}
	if params.Tag != nil {
		f.Tag = *params.Tag
	}
	if params.InScopeOnly != nil {
		f.InScopeOnly = *params.InScopeOnly
	}

	comps, total, err := h.OssComponentRepo.Search(ctx.Request().Context(), f)
	if err != nil {
		return err
	}

	for i := range comps {
		layers, err := h.OssComponentLayerRepo.ListByOssID(ctx.Request().Context(), comps[i].ID)
		if err != nil {
			return err
		}
		comps[i].Layers = layers
		tags, err := h.OssComponentTagRepo.ListByOssID(ctx.Request().Context(), comps[i].ID)
		if err != nil {
			return err
		}
		comps[i].Tags = tags
	}

	items := make([]gen.OssComponent, len(comps))
	for i, c := range comps {
		items[i] = toOssComponent(c)
	}
	res := gen.PagedResultOssComponent{
		Items: &items,
		Page:  &page,
		Size:  &size,
		Total: &total,
	}
	return ctx.JSON(http.StatusOK, res)
}

// OSSコンポーネント作成
// (POST /oss)
func (h *Handler) CreateOssComponent(ctx echo.Context) error {
	var req gen.OssComponentCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	now := dbtime.DBTime{Time: time.Now()}
	id := uuid.NewString()
	comp := &model.OssComponent{
		ID:              id,
		Name:            req.Name,
		NormalizedName:  strings.ToLower(req.Name),
		HomepageURL:     req.HomepageUrl,
		RepositoryURL:   req.RepositoryUrl,
		Description:     req.Description,
		PrimaryLanguage: req.PrimaryLanguage,
		Deprecated:      false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if req.DefaultUsageRole != nil {
		val := string(*req.DefaultUsageRole)
		comp.DefaultUsageRole = &val
	}
	if req.Layers != nil {
		ls := make([]string, len(*req.Layers))
		for i, l := range *req.Layers {
			ls[i] = string(l)
		}
		comp.Layers = ls
	}

	if err := h.OssComponentRepo.Create(ctx.Request().Context(), comp); err != nil {
		return err
	}
	if len(comp.Layers) > 0 {
		if err := h.OssComponentLayerRepo.Replace(ctx.Request().Context(), comp.ID, comp.Layers); err != nil {
			return err
		}
	}
	if req.TagIds != nil {
		ids := make([]string, len(*req.TagIds))
		for i, tid := range *req.TagIds {
			ids[i] = tid.String()
		}
		if err := h.OssComponentTagRepo.Replace(ctx.Request().Context(), comp.ID, ids); err != nil {
			return err
		}
	}
	tags, err := h.OssComponentTagRepo.ListByOssID(ctx.Request().Context(), comp.ID)
	if err != nil {
		return err
	}
	comp.Tags = tags
	res := toOssComponent(*comp)
	return ctx.JSON(http.StatusCreated, res)
}

// OSSコンポーネントを非推奨 (deprecated=true) に設定
// (DELETE /oss/{ossId})
func (*Handler) DeprecateOssComponent(ctx echo.Context, ossId openapi_types.UUID) error {
	return nil
}

// OSSコンポーネント詳細
// (GET /oss/{ossId})
func (h *Handler) GetOssComponent(ctx echo.Context, ossId openapi_types.UUID) error {
	comp, err := h.OssComponentRepo.Get(ctx.Request().Context(), ossId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "component not found")
		}
		return err
	}
	layers, err := h.OssComponentLayerRepo.ListByOssID(ctx.Request().Context(), comp.ID)
	if err != nil {
		return err
	}
	comp.Layers = layers
	tags, err := h.OssComponentTagRepo.ListByOssID(ctx.Request().Context(), comp.ID)
	if err != nil {
		return err
	}
	comp.Tags = tags
	res := toOssComponent(*comp)
	return ctx.JSON(http.StatusOK, res)
}

// OSSコンポーネント更新 (部分)
// (PATCH /oss/{ossId})
func (h *Handler) UpdateOssComponent(ctx echo.Context, ossId openapi_types.UUID) error {
	var req gen.OssComponentUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	comp, err := h.OssComponentRepo.Get(ctx.Request().Context(), ossId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "component not found")
		}
		return err
	}
	if req.Name != nil {
		comp.Name = *req.Name
		comp.NormalizedName = strings.ToLower(*req.Name)
	}
	if req.HomepageUrl != nil {
		comp.HomepageURL = req.HomepageUrl
	}
	if req.RepositoryUrl != nil {
		comp.RepositoryURL = req.RepositoryUrl
	}
	if req.Description != nil {
		comp.Description = req.Description
	}
	if req.PrimaryLanguage != nil {
		comp.PrimaryLanguage = req.PrimaryLanguage
	}
	if req.DefaultUsageRole != nil {
		val := string(*req.DefaultUsageRole)
		comp.DefaultUsageRole = &val
	}
	if req.Deprecated != nil {
		comp.Deprecated = *req.Deprecated
	}
	comp.UpdatedAt = dbtime.DBTime{Time: time.Now()}
	if err := h.OssComponentRepo.Update(ctx.Request().Context(), comp); err != nil {
		return err
	}
	if req.Layers != nil {
		ls := make([]string, len(*req.Layers))
		for i, l := range *req.Layers {
			ls[i] = string(l)
		}
		if err := h.OssComponentLayerRepo.Replace(ctx.Request().Context(), comp.ID, ls); err != nil {
			return err
		}
	}
	layers, err := h.OssComponentLayerRepo.ListByOssID(ctx.Request().Context(), comp.ID)
	if err != nil {
		return err
	}
	comp.Layers = layers
	if req.TagIds != nil {
		ids := make([]string, len(*req.TagIds))
		for i, id := range *req.TagIds {
			ids[i] = id.String()
		}
		if err := h.OssComponentTagRepo.Replace(ctx.Request().Context(), comp.ID, ids); err != nil {
			return err
		}
	}
	tags, err := h.OssComponentTagRepo.ListByOssID(ctx.Request().Context(), comp.ID)
	if err != nil {
		return err
	}
	comp.Tags = tags
	res := toOssComponent(*comp)
	return ctx.JSON(http.StatusOK, res)
}

// 指定 OSS のバージョン一覧
// (GET /oss/{ossId}/versions)
func (h *Handler) ListOssVersions(ctx echo.Context, ossId openapi_types.UUID, params gen.ListOssVersionsParams) error {
	page := 1
	if params.Page != nil {
		page = int(*params.Page)
	}
	size := 50
	if params.Size != nil {
		size = int(*params.Size)
	}
	f := domrepo.OssVersionFilter{
		OssID: ossId.String(),
		Page:  page,
		Size:  size,
	}
	if params.ReviewStatus != nil {
		f.ReviewStatus = string(*params.ReviewStatus)
	}
	if params.ScopeStatus != nil {
		f.ScopeStatus = string(*params.ScopeStatus)
	}
	vers, total, err := h.OssVersionRepo.Search(ctx.Request().Context(), f)
	if err != nil {
		return err
	}
	items := make([]gen.OssVersion, len(vers))
	for i, v := range vers {
		items[i] = toOssVersion(v)
	}
	res := gen.PagedResultOssVersion{
		Items: &items,
		Page:  &page,
		Size:  &size,
		Total: &total,
	}
	return ctx.JSON(http.StatusOK, res)
}

// バージョン追加
// (POST /oss/{ossId}/versions)
func (h *Handler) CreateOssVersion(ctx echo.Context, ossId openapi_types.UUID) error {
	var req gen.OssVersionCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	now := dbtime.DBTime{Time: time.Now()}
	v := &model.OssVersion{
		ID:           uuid.NewString(),
		OssID:        ossId.String(),
		Version:      req.Version,
		ReviewStatus: "draft",
		ScopeStatus:  string(gen.INSCOPE),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if req.ReleaseDate != nil {
		t := dbtime.DBTime{Time: req.ReleaseDate.Time}
		v.ReleaseDate = &t
	}
	if req.LicenseExpressionRaw != nil {
		v.LicenseExpressionRaw = req.LicenseExpressionRaw
	}
	if req.Purl != nil {
		v.Purl = req.Purl
	}
	if req.CpeList != nil {
		v.CpeList = *req.CpeList
	}
	if req.HashSha256 != nil {
		v.HashSha256 = req.HashSha256
	}
	if req.Modified != nil {
		v.Modified = *req.Modified
	}
	if req.ModificationDescription != nil {
		v.ModificationDescription = req.ModificationDescription
	}
	if req.SupplierType != nil {
		val := string(*req.SupplierType)
		v.SupplierType = &val
	}
	if req.ForkOriginUrl != nil {
		v.ForkOriginURL = req.ForkOriginUrl
	}
	if err := h.OssVersionRepo.Create(ctx.Request().Context(), v); err != nil {
		return err
	}
	res := toOssVersion(*v)
	return ctx.JSON(http.StatusCreated, res)
}

// バージョン削除 (論理/物理は実装方針による)
// (DELETE /oss/{ossId}/versions/{versionId})
func (h *Handler) DeleteOssVersion(ctx echo.Context, ossId openapi_types.UUID, versionId openapi_types.UUID) error {
	if err := h.OssVersionRepo.Delete(ctx.Request().Context(), versionId.String()); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

// バージョン詳細
// (GET /oss/{ossId}/versions/{versionId})
func (h *Handler) GetOssVersion(ctx echo.Context, ossId openapi_types.UUID, versionId openapi_types.UUID) error {
	v, err := h.OssVersionRepo.Get(ctx.Request().Context(), versionId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "version not found")
		}
		return err
	}
	res := toOssVersion(*v)
	return ctx.JSON(http.StatusOK, res)
}

// バージョン更新
// (PATCH /oss/{ossId}/versions/{versionId})
func (h *Handler) UpdateOssVersion(ctx echo.Context, ossId openapi_types.UUID, versionId openapi_types.UUID) error {
	var req gen.OssVersionUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	v, err := h.OssVersionRepo.Get(ctx.Request().Context(), versionId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "version not found")
		}
		return err
	}
	if req.ReleaseDate != nil {
		t := dbtime.DBTime{Time: req.ReleaseDate.Time}
		v.ReleaseDate = &t
	}
	if req.LicenseExpressionRaw != nil {
		v.LicenseExpressionRaw = req.LicenseExpressionRaw
	}
	if req.LicenseConcluded != nil {
		v.LicenseConcluded = req.LicenseConcluded
	}
	if req.Purl != nil {
		v.Purl = req.Purl
	}
	if req.CpeList != nil {
		v.CpeList = *req.CpeList
	}
	if req.HashSha256 != nil {
		v.HashSha256 = req.HashSha256
	}
	if req.Modified != nil {
		v.Modified = *req.Modified
	}
	if req.ModificationDescription != nil {
		v.ModificationDescription = req.ModificationDescription
	}
	if req.ReviewStatus != nil {
		v.ReviewStatus = string(*req.ReviewStatus)
	}
	if req.ScopeStatus != nil {
		v.ScopeStatus = string(*req.ScopeStatus)
	}
	if req.SupplierType != nil {
		val := string(*req.SupplierType)
		v.SupplierType = &val
	}
	if req.ForkOriginUrl != nil {
		v.ForkOriginURL = req.ForkOriginUrl
	}
	v.UpdatedAt = dbtime.DBTime{Time: time.Now()}
	if err := h.OssVersionRepo.Update(ctx.Request().Context(), v); err != nil {
		return err
	}
	res := toOssVersion(*v)
	return ctx.JSON(http.StatusOK, res)
}
