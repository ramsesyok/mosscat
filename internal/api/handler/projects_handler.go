package handler

// projects_handler.go - /projects に関するハンドラ処理

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/domain/model"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

func toProject(m model.Project) gen.Project {
	uid := uuid.MustParse(m.ID)
	res := gen.Project{
		Id:          uid,
		ProjectCode: m.ProjectCode,
		Name:        m.Name,
		CreatedAt:   m.CreatedAt.TimeValue(),
		UpdatedAt:   m.UpdatedAt.TimeValue(),
	}
	if m.Department != nil {
		res.Department = m.Department
	}
	if m.Manager != nil {
		res.Manager = m.Manager
	}
	if m.DeliveryDate != nil {
		res.DeliveryDate = &openapi_types.Date{Time: m.DeliveryDate.TimeValue()}
	}
	if m.Description != nil {
		res.Description = m.Description
	}
	if m.OssUsageCount != 0 {
		res.OssUsageCount = &m.OssUsageCount
	}
	return res
}

func toProjectUsage(m model.ProjectUsage) gen.ProjectUsage {
	uid := uuid.MustParse(m.ID)
	pid := uuid.MustParse(m.ProjectID)
	oid := uuid.MustParse(m.OssID)
	vid := uuid.MustParse(m.OssVersionID)
	res := gen.ProjectUsage{
		Id:               uid,
		ProjectId:        pid,
		OssId:            oid,
		OssVersionId:     vid,
		UsageRole:        gen.UsageRole(m.UsageRole),
		ScopeStatus:      gen.ScopeStatus(m.ScopeStatus),
		DirectDependency: m.DirectDependency,
		AddedAt:          m.AddedAt.TimeValue(),
		Modified:         m.Modified,
	}
	if m.InclusionNote != nil {
		res.InclusionNote = m.InclusionNote
	}
	if m.EvaluatedAt != nil {
		t := m.EvaluatedAt.TimeValue()
		res.EvaluatedAt = &t
	}
	if m.EvaluatedBy != nil {
		res.EvaluatedBy = m.EvaluatedBy
	}
	if m.ModificationDescription != nil {
		res.ModificationDescription = m.ModificationDescription
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

func initialScopeStatus(policy *model.ScopePolicy, role string) string {
	switch role {
	case "BUILD_ONLY", "DEV_ONLY", "TEST_ONLY":
		return string(gen.OUTSCOPE)
	case "SERVER_ENV":
		if policy != nil && policy.ServerEnvIncluded {
			return string(gen.INSCOPE)
		}
		return string(gen.OUTSCOPE)
	case "RUNTIME_REQUIRED":
		if policy != nil && policy.RuntimeRequiredDefaultInScope {
			return string(gen.INSCOPE)
		}
		return string(gen.REVIEWNEEDED)
	default:
		return string(gen.INSCOPE)
	}
}

// プロジェクト一覧
// (GET /projects)
func (h *Handler) ListProjects(ctx echo.Context, params gen.ListProjectsParams) error {
	page := 1
	if params.Page != nil {
		page = int(*params.Page)
	}
	size := 50
	if params.Size != nil {
		size = int(*params.Size)
	}
	f := domrepo.ProjectFilter{Page: page, Size: size}
	if params.Code != nil {
		f.Code = *params.Code
	}
	if params.Name != nil {
		f.Name = *params.Name
	}

	projects, total, err := h.ProjectRepo.Search(ctx.Request().Context(), f)
	if err != nil {
		return err
	}

	items := make([]gen.Project, len(projects))
	for i, p := range projects {
		items[i] = toProject(p)
	}
	res := gen.PagedResultProject{
		Items: &items,
		Page:  &page,
		Size:  &size,
		Total: &total,
	}
	return ctx.JSON(http.StatusOK, res)
}

// プロジェクト作成
// (POST /projects)
func (h *Handler) CreateProject(ctx echo.Context) error {
	var req gen.ProjectCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	now := dbtime.DBTime{Time: time.Now()}
	p := &model.Project{
		ID:          uuid.NewString(),
		ProjectCode: req.ProjectCode,
		Name:        req.Name,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if req.Department != nil {
		p.Department = req.Department
	}
	if req.Manager != nil {
		p.Manager = req.Manager
	}
	if req.DeliveryDate != nil {
		t := dbtime.DBTime{Time: req.DeliveryDate.Time}
		p.DeliveryDate = &t
	}
	if req.Description != nil {
		p.Description = req.Description
	}

	if err := h.ProjectRepo.Create(ctx.Request().Context(), p); err != nil {
		return err
	}

	res := toProject(*p)
	return ctx.JSON(http.StatusCreated, res)
}

// プロジェクト削除 (論理予定)
// (DELETE /projects/{projectId})
func (h *Handler) DeleteProject(ctx echo.Context, projectId openapi_types.UUID) error {
	if err := h.ProjectRepo.Delete(ctx.Request().Context(), projectId.String()); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

// プロジェクト詳細
// (GET /projects/{projectId})
func (h *Handler) GetProject(ctx echo.Context, projectId openapi_types.UUID) error {
	p, err := h.ProjectRepo.Get(ctx.Request().Context(), projectId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "project not found")
		}
		return err
	}
	res := toProject(*p)
	return ctx.JSON(http.StatusOK, res)
}

// プロジェクト更新
// (PATCH /projects/{projectId})
func (h *Handler) UpdateProject(ctx echo.Context, projectId openapi_types.UUID) error {
	var req gen.ProjectUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	p, err := h.ProjectRepo.Get(ctx.Request().Context(), projectId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "project not found")
		}
		return err
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Department != nil {
		p.Department = req.Department
	}
	if req.Manager != nil {
		p.Manager = req.Manager
	}
	if req.DeliveryDate != nil {
		t := dbtime.DBTime{Time: req.DeliveryDate.Time}
		p.DeliveryDate = &t
	}
	if req.Description != nil {
		p.Description = req.Description
	}
	p.UpdatedAt = dbtime.DBTime{Time: time.Now()}

	if err := h.ProjectRepo.Update(ctx.Request().Context(), p); err != nil {
		return err
	}

	res := toProject(*p)
	return ctx.JSON(http.StatusOK, res)
}

// プロジェクト納品用エクスポート (プレーホルダ)
// (GET /projects/{projectId}/export)
func (*Handler) ExportProjectArtifacts(ctx echo.Context, projectId openapi_types.UUID, params gen.ExportProjectArtifactsParams) error {
	// 一時的な実装
	return ctx.JSON(http.StatusOK, map[string]any{"placeholder": "todo"})
}

// プロジェクト中利用 OSS 一覧
// (GET /projects/{projectId}/usages)
func (h *Handler) ListProjectUsages(ctx echo.Context, projectId openapi_types.UUID, params gen.ListProjectUsagesParams) error {
	page := 1
	if params.Page != nil {
		page = int(*params.Page)
	}
	size := 50
	if params.Size != nil {
		size = int(*params.Size)
	}
	f := domrepo.ProjectUsageFilter{ProjectID: projectId.String(), Page: page, Size: size}
	if params.ScopeStatus != nil {
		f.ScopeStatus = string(*params.ScopeStatus)
	}
	if params.UsageRole != nil {
		f.UsageRole = string(*params.UsageRole)
	}
	if params.Direct != nil {
		f.Direct = params.Direct
	}

	usages, total, err := h.ProjectUsageRepo.Search(ctx.Request().Context(), f)
	if err != nil {
		return err
	}
	items := make([]gen.ProjectUsage, len(usages))
	for i, u := range usages {
		items[i] = toProjectUsage(u)
	}
	res := gen.PagedResultProjectUsage{
		Items: &items,
		Page:  &page,
		Size:  &size,
		Total: &total,
	}
	return ctx.JSON(http.StatusOK, res)
}

// プロジェクト利用追加
// (POST /projects/{projectId}/usages)
func (h *Handler) CreateProjectUsage(ctx echo.Context, projectId openapi_types.UUID) error {
	var req gen.ProjectUsageCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	now := dbtime.DBTime{Time: time.Now()}
	policy, err := h.ScopePolicyRepo.Get(ctx.Request().Context())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	scope := initialScopeStatus(policy, string(req.UsageRole))
	direct := true
	if req.DirectDependency != nil {
		direct = *req.DirectDependency
	}
	modified := false
	if req.Modified != nil {
		modified = *req.Modified
	}
	u := &model.ProjectUsage{
		ID:                      uuid.NewString(),
		ProjectID:               projectId.String(),
		OssID:                   req.OssId.String(),
		OssVersionID:            req.OssVersionId.String(),
		UsageRole:               string(req.UsageRole),
		DirectDependency:        direct,
		InclusionNote:           req.InclusionNote,
		ScopeStatus:             scope,
		AddedAt:                 now,
		Modified:                modified,
		ModificationDescription: req.ModificationDescription,
	}
	if req.SupplierType != nil {
		val := string(*req.SupplierType)
		u.SupplierType = &val
	}
	if req.ForkOriginUrl != nil {
		u.ForkOriginURL = req.ForkOriginUrl
	}
	if err := h.ProjectUsageRepo.Create(ctx.Request().Context(), u); err != nil {
		return err
	}
	res := toProjectUsage(*u)
	return ctx.JSON(http.StatusCreated, res)
}

// 利用削除
// (DELETE /projects/{projectId}/usages/{usageId})
func (h *Handler) DeleteProjectUsage(ctx echo.Context, projectId openapi_types.UUID, usageId openapi_types.UUID) error {
	if err := h.ProjectUsageRepo.Delete(ctx.Request().Context(), usageId.String()); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

// 利用情報更新
// (PATCH /projects/{projectId}/usages/{usageId})
func (h *Handler) UpdateProjectUsage(ctx echo.Context, projectId openapi_types.UUID, usageId openapi_types.UUID) error {
	var req gen.ProjectUsageUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	u := &model.ProjectUsage{
		ID:        usageId.String(),
		ProjectID: projectId.String(),
	}
	if req.OssVersionId != nil {
		u.OssVersionID = req.OssVersionId.String()
	}
	if req.UsageRole != nil {
		u.UsageRole = string(*req.UsageRole)
	}
	if req.DirectDependency != nil {
		u.DirectDependency = *req.DirectDependency
	}
	if req.InclusionNote != nil {
		u.InclusionNote = req.InclusionNote
	}
	if req.ScopeStatus != nil {
		u.ScopeStatus = string(*req.ScopeStatus)
	}
	if req.Modified != nil {
		u.Modified = *req.Modified
	}
	if req.ModificationDescription != nil {
		u.ModificationDescription = req.ModificationDescription
	}
	if req.SupplierType != nil {
		val := string(*req.SupplierType)
		u.SupplierType = &val
	}
	if req.ForkOriginUrl != nil {
		u.ForkOriginURL = req.ForkOriginUrl
	}

	if err := h.ProjectUsageRepo.Update(ctx.Request().Context(), u); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]any{})
}

// スコープ判定更新
// (PATCH /projects/{projectId}/usages/{usageId}/scope)
func (h *Handler) UpdateProjectUsageScope(ctx echo.Context, projectId openapi_types.UUID, usageId openapi_types.UUID) error {
	var req gen.ScopeStatusUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	now := dbtime.DBTime{Time: time.Now()}
	evaluatedBy := "api-user"
	if err := h.ProjectUsageRepo.UpdateScope(ctx.Request().Context(), usageId.String(), string(req.ScopeStatus), req.ReasonNote, now, &evaluatedBy); err != nil {
		return err
	}
	u := model.ProjectUsage{
		ID:            usageId.String(),
		ProjectID:     projectId.String(),
		ScopeStatus:   string(req.ScopeStatus),
		InclusionNote: req.ReasonNote,
		EvaluatedAt:   &now,
		EvaluatedBy:   &evaluatedBy,
	}
	_ = u // 実装完了までのプレースホルダ
	return ctx.JSON(http.StatusOK, map[string]any{})
}
