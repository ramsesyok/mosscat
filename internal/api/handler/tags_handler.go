package handler

// tags_handler.go - /tags に関するハンドラ処理

import (
	"net/http"
	"time"

	"github.com/ramsesyok/mosscat/pkg/dbtime"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/domain/model"
)

func toTag(m model.Tag) gen.Tag {
	return gen.Tag{Id: uuid.MustParse(m.ID), Name: m.Name, CreatedAt: func() *time.Time {
		if m.CreatedAt == nil {
			return nil
		}
		t := m.CreatedAt.TimeValue()
		return &t
	}()}
}

// タグ一覧
// (GET /tags)
func (h *Handler) ListTags(ctx echo.Context) error {
	tags, err := h.TagRepo.List(ctx.Request().Context())
	if err != nil {
		return err
	}
	res := make([]gen.Tag, len(tags))
	for i, t := range tags {
		res[i] = toTag(t)
	}
	return ctx.JSON(http.StatusOK, res)
}

// タグ作成
// (POST /tags)
func (h *Handler) CreateTag(ctx echo.Context) error {
	var req gen.TagCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	now := dbtime.DBTime{Time: time.Now()}
	tag := &model.Tag{ID: uuid.NewString(), Name: req.Name, CreatedAt: &now}
	if err := h.TagRepo.Create(ctx.Request().Context(), tag); err != nil {
		return err
	}
	res := toTag(*tag)
	return ctx.JSON(http.StatusCreated, res)
}

// タグ削除
// (DELETE /tags/{tagId})
func (h *Handler) DeleteTag(ctx echo.Context, tagId openapi_types.UUID) error {
	if err := h.TagRepo.Delete(ctx.Request().Context(), tagId.String()); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
