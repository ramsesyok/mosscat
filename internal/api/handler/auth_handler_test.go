package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	apirouter "github.com/ramsesyok/mosscat/internal/api"
	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	handler "github.com/ramsesyok/mosscat/internal/api/handler"
	"github.com/ramsesyok/mosscat/internal/domain/model"
	infrarepo "github.com/ramsesyok/mosscat/internal/infra/repository"
	"github.com/ramsesyok/mosscat/pkg/auth"
)

func setupAuthEcho(h *handler.Handler) (*echo.Echo, string) {
	os.Setenv("JWT_SECRET", "testsecret")
	os.Setenv("JWT_EXPIRES_MIN", "1")
	e := echo.New()
	apirouter.RegisterRoutes(e, h)
	hash, _ := auth.Hash("pass")
	u := &model.User{ID: uuid.NewString(), Username: "admin", PasswordHash: hash, Roles: []string{"ADMIN"}, Active: true}
	token, _, _ := auth.GenerateToken(u)
	return e, token
}

func TestLoginAndAccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &handler.Handler{UserRepo: repo}
	e, _ := setupAuthEcho(h)

	now := time.Now()
	uid := uuid.NewString()
	query := regexp.QuoteMeta("SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE username = ?")
	hash, _ := auth.Hash("pass")
	mock.ExpectQuery(query).WithArgs("admin").WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "display_name", "email", "password_hash", "roles", "active", "created_at", "updated_at"}).
			AddRow(uid, "admin", nil, nil, hash, pq.StringArray{"ADMIN"}, true, now, now),
	)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"admin","password":"pass"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var res gen.LoginResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotEmpty(t, res.AccessToken)
	t.Logf("token=%s", res.AccessToken)

	// expectation for GetCurrentUser
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE id = ?")).
		WithArgs(uid).WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "display_name", "email", "password_hash", "roles", "active", "created_at", "updated_at"}).
			AddRow(uid, "admin", nil, nil, hash, pq.StringArray{"ADMIN"}, true, now, now),
	)
	req2 := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req2.Header.Set("Authorization", "Bearer "+res.AccessToken)
	rec2 := httptest.NewRecorder()
	e.ServeHTTP(rec2, req2)
	t.Logf("me body=%s", rec2.Body.String())
	require.Equal(t, http.StatusOK, rec2.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_InvalidPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &handler.Handler{UserRepo: repo}
	e, _ := setupAuthEcho(h)

	now := time.Now()
	uid := uuid.NewString()
	query := regexp.QuoteMeta("SELECT id, username, display_name, email, password_hash, roles, active, created_at, updated_at FROM users WHERE username = ?")
	hash, _ := auth.Hash("pass")
	mock.ExpectQuery(query).WithArgs("admin").WillReturnRows(
		sqlmock.NewRows([]string{"id", "username", "display_name", "email", "password_hash", "roles", "active", "created_at", "updated_at"}).
			AddRow(uid, "admin", nil, nil, hash, pq.StringArray{"ADMIN"}, true, now, now),
	)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"admin","password":"wrong"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRolesRequired_Forbidden(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &handler.Handler{UserRepo: repo}
	e, token := setupAuthEcho(h)

	e.GET("/api/admin", func(c echo.Context) error { return c.String(200, "ok") }, auth.RolesRequired("ADMIN"))

	req := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
	// token for admin already generated but we need viewer role token to test deny
	viewer := &model.User{ID: uuid.NewString(), Username: "view", PasswordHash: "pass", Roles: []string{"VIEWER"}, Active: true}
	vtoken, _, _ := auth.GenerateToken(viewer)
	req.Header.Set("Authorization", "Bearer "+vtoken)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	_ = token // silence unused
}

func TestExpiredToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := &infrarepo.UserRepository{DB: db}
	h := &handler.Handler{UserRepo: repo}
	os.Setenv("JWT_SECRET", "testsecret")
	os.Setenv("JWT_EXPIRES_MIN", "-1")
	e := echo.New()
	apirouter.RegisterRoutes(e, h)

	expiredUser := &model.User{ID: uuid.NewString(), Username: "a", PasswordHash: "p", Roles: []string{"ADMIN"}, Active: true}
	token, _, _ := auth.GenerateToken(expiredUser)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	t.Logf("expired body=%s status=%d", rec.Body.String(), rec.Code)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRouterRolesForbidden(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &infrarepo.UserRepository{DB: db}
	h := &handler.Handler{UserRepo: repo}
	e, _ := setupAuthEcho(h)

	viewer := &model.User{ID: uuid.NewString(), Username: "view", PasswordHash: "pass", Roles: []string{"VIEWER"}, Active: true}
	token, _, _ := auth.GenerateToken(viewer)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}
