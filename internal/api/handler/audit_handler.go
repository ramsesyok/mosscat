package handler

// audit_handler.go - /audit に関するハンドラ処理

import (
	"net/http"

	"github.com/labstack/echo/v4"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
	"github.com/ramsesyok/mosscat/pkg/dbtime"
)

// 監査ログ簡易検索 (Phase1簡易)
// (GET /audit)
func (h *Handler) SearchAuditLogs(ctx echo.Context, params gen.SearchAuditLogsParams) error {
	var from, to *dbtime.DBTime
	if params.From != nil {
		v := dbtime.DBTime{Time: params.From.UTC()}
		from = &v
	}
	if params.To != nil {
		v := dbtime.DBTime{Time: params.To.UTC()}
		to = &v
	}
	filter := domrepo.AuditLogFilter{
		EntityType: params.EntityType,
		EntityID:   params.EntityId,
		From:       from,
		To:         to,
	}

	logs, err := h.AuditRepo.Search(ctx.Request().Context(), filter)
	if err != nil {
		return err
	}

	items := make([]map[string]any, len(logs))
	for i, l := range logs {
		item := map[string]any{
			"id":         l.ID,
			"entityType": l.EntityType,
			"entityId":   l.EntityID,
			"action":     l.Action,
			"at":         l.CreatedAt.TimeValue(),
			"user":       l.UserName,
		}
		if l.Summary != nil {
			item["summary"] = *l.Summary
		}
		items[i] = item
	}

	return ctx.JSON(http.StatusOK, map[string]any{"items": items})
}
