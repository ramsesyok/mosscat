package api

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/api/handler"
	"github.com/ramsesyok/mosscat/pkg/auth"
	problem "github.com/ramsesyok/mosscat/pkg/response"
)

// RegisterRoutes registers handlers with JWT auth using generated RegisterHandlersWithBaseURL.
func RegisterRoutes(e *echo.Echo, h *handler.Handler) {
	// JWT認証設定
	cfg := echojwt.Config{
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
		ContextKey:    "authUser",
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return &auth.Claims{} },
		ErrorHandler: func(c echo.Context, err error) error {
			return problem.Unauthorized(c, "UNAUTHORIZED", "token invalid or expired")
		},
	}
	authRequired := echojwt.WithConfig(cfg)

	// /api プレフィックス付きでAPIハンドラーを登録
	gen.RegisterHandlersWithBaseURL(e, h, "/api")

	// 認証が必要なエンドポイントにJWTミドルウェアを適用
	wrapper := gen.ServerInterfaceWrapper{Handler: h}
	
	// ログインは認証不要
	// その他のエンドポイントに認証を適用
	apiGroup := e.Group("/api", authRequired)
	apiGroup.POST("/auth/logout", wrapper.Logout, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.GET("/audit", wrapper.SearchAuditLogs, auth.RolesRequired("ADMIN"))
	apiGroup.GET("/me", wrapper.GetCurrentUser, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.GET("/oss", wrapper.ListOssComponents, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.POST("/oss", wrapper.CreateOssComponent, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.DELETE("/oss/:ossId", wrapper.DeprecateOssComponent, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.GET("/oss/:ossId", wrapper.GetOssComponent, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.PATCH("/oss/:ossId", wrapper.UpdateOssComponent, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.GET("/oss/:ossId/versions", wrapper.ListOssVersions, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.POST("/oss/:ossId/versions", wrapper.CreateOssVersion, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.DELETE("/oss/:ossId/versions/:versionId", wrapper.DeleteOssVersion, auth.RolesRequired("ADMIN"))
	apiGroup.GET("/oss/:ossId/versions/:versionId", wrapper.GetOssVersion, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.PATCH("/oss/:ossId/versions/:versionId", wrapper.UpdateOssVersion, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.GET("/projects", wrapper.ListProjects, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.POST("/projects", wrapper.CreateProject, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.DELETE("/projects/:projectId", wrapper.DeleteProject, auth.RolesRequired("ADMIN"))
	apiGroup.GET("/projects/:projectId", wrapper.GetProject, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.PATCH("/projects/:projectId", wrapper.UpdateProject, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.GET("/projects/:projectId/export", wrapper.ExportProjectArtifacts, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.GET("/projects/:projectId/usages", wrapper.ListProjectUsages, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.POST("/projects/:projectId/usages", wrapper.CreateProjectUsage, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.DELETE("/projects/:projectId/usages/:usageId", wrapper.DeleteProjectUsage, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.PATCH("/projects/:projectId/usages/:usageId", wrapper.UpdateProjectUsage, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.PATCH("/projects/:projectId/usages/:usageId/scope", wrapper.UpdateProjectUsageScope, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.GET("/scope/policy", wrapper.GetScopePolicy, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.PATCH("/scope/policy", wrapper.UpdateScopePolicy, auth.RolesRequired("ADMIN"))
	apiGroup.GET("/tags", wrapper.ListTags, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	apiGroup.POST("/tags", wrapper.CreateTag, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.DELETE("/tags/:tagId", wrapper.DeleteTag, auth.RolesRequired("EDITOR", "ADMIN"))
	apiGroup.GET("/users", wrapper.ListUsers, auth.RolesRequired("ADMIN"))
	apiGroup.POST("/users", wrapper.CreateUser, auth.RolesRequired("ADMIN"))
	apiGroup.DELETE("/users/:userId", wrapper.DeleteUser, auth.RolesRequired("ADMIN"))
	apiGroup.GET("/users/:userId", wrapper.GetUser, auth.RolesRequired("ADMIN"))
	apiGroup.PATCH("/users/:userId", wrapper.UpdateUser, auth.RolesRequired("ADMIN"))
}
