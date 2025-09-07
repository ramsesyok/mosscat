package handler

// scope_policy_handler.go - /scope に関するハンドラ処理

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/domain/model"
)

func toScopePolicy(m model.ScopePolicy) gen.ScopePolicy {
	uid := uuid.MustParse(m.ID)
	return gen.ScopePolicy{
		Id:                            &uid,
		RuntimeRequiredDefaultInScope: &m.RuntimeRequiredDefaultInScope,
		ServerEnvIncluded:             &m.ServerEnvIncluded,
		AutoMarkForksInScope:          &m.AutoMarkForksInScope,
		UpdatedAt:                     func() *time.Time { t := m.UpdatedAt.TimeValue(); return &t }(),
		UpdatedBy:                     &m.UpdatedBy,
	}
}

// 現行スコープポリシー取得
// (GET /scope/policy)
func (h *Handler) GetScopePolicy(ctx echo.Context) error {
	p, err := h.ScopePolicyRepo.Get(ctx.Request().Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "policy not found")
		}
		return err
	}
	res := toScopePolicy(*p)
	return ctx.JSON(http.StatusOK, res)
}

// スコープポリシー更新 (管理者)
// (PATCH /scope/policy)
func (h *Handler) UpdateScopePolicy(ctx echo.Context) error {
	var req gen.ScopePolicyUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	ctxx := ctx.Request().Context()
	existing, err := h.ScopePolicyRepo.Get(ctxx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	now := dbtime.DBTime{Time: time.Now()}
	var p *model.ScopePolicy
	if existing != nil {
		p = existing
	} else {
		p = &model.ScopePolicy{ID: uuid.NewString()}
	}

	if req.RuntimeRequiredDefaultInScope != nil {
		p.RuntimeRequiredDefaultInScope = *req.RuntimeRequiredDefaultInScope
	}
	if req.ServerEnvIncluded != nil {
		p.ServerEnvIncluded = *req.ServerEnvIncluded
	}
	if req.AutoMarkForksInScope != nil {
		p.AutoMarkForksInScope = *req.AutoMarkForksInScope
	}
	p.UpdatedAt = now
	p.UpdatedBy = "api-user"

	if err := h.ScopePolicyRepo.Update(ctxx, p); err != nil {
		return err
	}
	res := toScopePolicy(*p)
	return ctx.JSON(http.StatusOK, res)
}
