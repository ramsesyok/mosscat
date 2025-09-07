package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/domain/service"
	"github.com/ramsesyok/mosscat/pkg/auth"
	problem "github.com/ramsesyok/mosscat/pkg/response"
)

// Login issues JWT token.
func (h *Handler) Login(ctx echo.Context) error {
	var req gen.LoginJSONBody
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	svc := service.AuthService{UserRepo: h.UserRepo}
	u, err := svc.Authenticate(ctx.Request().Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserDisabled):
			return problem.Forbidden(ctx, "USER_DISABLED", "user disabled")
		case errors.Is(err, service.ErrInvalidCredential):
			return problem.Unauthorized(ctx, "INVALID_CREDENTIAL", "invalid username or password")
		default:
			return err
		}
	}

	token, exp, err := auth.GenerateToken(u)
	if err != nil {
		return err
	}
	res := gen.LoginResponse{AccessToken: token, ExpiresIn: exp}
	return ctx.JSON(http.StatusOK, res)
}

// Logout is a no-op placeholder.
func (h *Handler) Logout(ctx echo.Context) error {
	// TODO: write audit log of logout event
	return ctx.NoContent(http.StatusNoContent)
}

// GetCurrentUser returns current user info based on JWT claims.
func (h *Handler) GetCurrentUser(ctx echo.Context) error {
	claims := auth.GetClaims(ctx)
	if claims == nil {
		return problem.Unauthorized(ctx, "UNAUTHORIZED", "no claims")
	}
	u, err := h.UserRepo.Get(ctx.Request().Context(), claims.Sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return err
	}
	res := toUser(*u)
	return ctx.JSON(http.StatusOK, res)
}
