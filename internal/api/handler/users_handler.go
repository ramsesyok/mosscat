package handler

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
	"github.com/ramsesyok/mosscat/pkg/auth"
)

func toUser(m model.User) gen.User {
	uid := uuid.MustParse(m.ID)
	res := gen.User{
		Id:        uid,
		Username:  m.Username,
		Active:    m.Active,
		CreatedAt: m.CreatedAt.TimeValue(),
		UpdatedAt: m.UpdatedAt.TimeValue(),
	}
	if m.DisplayName != nil {
		res.DisplayName = m.DisplayName
	}
	if m.Email != nil {
		email := openapi_types.Email(*m.Email)
		res.Email = &email
	}
	if len(m.Roles) > 0 {
		roles := make([]gen.Role, len(m.Roles))
		for i, r := range m.Roles {
			roles[i] = gen.Role(r)
		}
		res.Roles = roles
	}
	return res
}

// ---- Users ----

// ListUsers ユーザー一覧 (GET /users)
func (h *Handler) ListUsers(ctx echo.Context, params gen.ListUsersParams) error {
	page := 1
	if params.Page != nil {
		page = int(*params.Page)
	}
	size := 50
	if params.Size != nil {
		size = int(*params.Size)
	}
	f := domrepo.UserFilter{Page: page, Size: size}
	if params.Username != nil {
		f.Username = *params.Username
	}
	if params.Role != nil {
		f.Role = string(*params.Role)
	}
	users, total, err := h.UserRepo.Search(ctx.Request().Context(), f)
	if err != nil {
		return err
	}
	items := make([]gen.User, len(users))
	for i, u := range users {
		items[i] = toUser(u)
	}
	res := gen.PagedResultUser{Items: &items, Page: &page, Size: &size, Total: &total}
	return ctx.JSON(http.StatusOK, res)
}

// CreateUser ユーザー作成 (POST /users)
func (h *Handler) CreateUser(ctx echo.Context) error {
	var req gen.UserCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	now := dbtime.DBTime{Time: time.Now()}
	id := uuid.NewString()
	active := true
	if req.Active != nil {
		active = *req.Active
	}
	var emailStr *string
	if req.Email != nil {
		v := string(*req.Email)
		emailStr = &v
	}
	roles := make([]string, len(req.Roles))
	for i, r := range req.Roles {
		roles[i] = string(r)
	}
	u := &model.User{
		ID:           id,
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		Email:        emailStr,
		PasswordHash: "",
		Roles:        roles,
		Active:       active,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if req.Password != nil {
		hash, err := auth.Hash(*req.Password)
		if err != nil {
			return err
		}
		u.PasswordHash = hash
	}
	if err := h.UserRepo.Create(ctx.Request().Context(), u); err != nil {
		return err
	}
	res := toUser(*u)
	return ctx.JSON(http.StatusCreated, res)
}

// GetUser ユーザー詳細 (GET /users/{userId})
func (h *Handler) GetUser(ctx echo.Context, userId openapi_types.UUID) error {
	u, err := h.UserRepo.Get(ctx.Request().Context(), userId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return err
	}
	res := toUser(*u)
	return ctx.JSON(http.StatusOK, res)
}

// UpdateUser ユーザー更新 (PATCH /users/{userId})
func (h *Handler) UpdateUser(ctx echo.Context, userId openapi_types.UUID) error {
	var req gen.UserUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	u, err := h.UserRepo.Get(ctx.Request().Context(), userId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return err
	}
	if req.DisplayName != nil {
		u.DisplayName = req.DisplayName
	}
	if req.Email != nil {
		v := string(*req.Email)
		u.Email = &v
	}
	if req.Password != nil {
		hash, err := auth.Hash(*req.Password)
		if err != nil {
			return err
		}
		u.PasswordHash = hash
	}
	if req.Roles != nil {
		roles := make([]string, len(*req.Roles))
		for i, r := range *req.Roles {
			roles[i] = string(r)
		}
		u.Roles = roles
	}
	if req.Active != nil {
		u.Active = *req.Active
	}
	u.UpdatedAt = dbtime.DBTime{Time: time.Now()}
	if err := h.UserRepo.Update(ctx.Request().Context(), u); err != nil {
		return err
	}
	res := toUser(*u)
	return ctx.JSON(http.StatusOK, res)
}

// DeleteUser ユーザー削除 (DELETE /users/{userId})
func (h *Handler) DeleteUser(ctx echo.Context, userId openapi_types.UUID) error {
	if err := h.UserRepo.Delete(ctx.Request().Context(), userId.String()); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
