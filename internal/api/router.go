package api

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/api/handler"
	"github.com/ramsesyok/mosscat/pkg/auth"
	problem "github.com/ramsesyok/mosscat/pkg/response"
	"os"
)

// RegisterRoutes registers handlers with JWT auth.
func RegisterRoutes(e *echo.Echo, h *handler.Handler) {
	cfg := echojwt.Config{
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
		ContextKey:    "authUser",
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return &auth.Claims{} },
		ErrorHandler: func(c echo.Context, err error) error {
			return problem.Unauthorized(c, "UNAUTHORIZED", "token invalid or expired")
		},
	}
	authRequired := echojwt.WithConfig(cfg)

	wrapper := gen.ServerInterfaceWrapper{Handler: h}

	e.POST("/auth/login", wrapper.Login)
	e.POST("/auth/logout", wrapper.Logout, authRequired, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))

	g := e.Group("", authRequired)
	g.GET("/audit", wrapper.SearchAuditLogs, auth.RolesRequired("ADMIN"))
	g.GET("/me", wrapper.GetCurrentUser, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.GET("/oss", wrapper.ListOssComponents, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.POST("/oss", wrapper.CreateOssComponent, auth.RolesRequired("EDITOR", "ADMIN"))
	g.DELETE("/oss/:ossId", wrapper.DeprecateOssComponent, auth.RolesRequired("EDITOR", "ADMIN"))
	g.GET("/oss/:ossId", wrapper.GetOssComponent, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.PATCH("/oss/:ossId", wrapper.UpdateOssComponent, auth.RolesRequired("EDITOR", "ADMIN"))
	g.GET("/oss/:ossId/versions", wrapper.ListOssVersions, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.POST("/oss/:ossId/versions", wrapper.CreateOssVersion, auth.RolesRequired("EDITOR", "ADMIN"))
	g.DELETE("/oss/:ossId/versions/:versionId", wrapper.DeleteOssVersion, auth.RolesRequired("ADMIN"))
	g.GET("/oss/:ossId/versions/:versionId", wrapper.GetOssVersion, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.PATCH("/oss/:ossId/versions/:versionId", wrapper.UpdateOssVersion, auth.RolesRequired("EDITOR", "ADMIN"))
	g.GET("/projects", wrapper.ListProjects, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.POST("/projects", wrapper.CreateProject, auth.RolesRequired("EDITOR", "ADMIN"))
	g.DELETE("/projects/:projectId", wrapper.DeleteProject, auth.RolesRequired("ADMIN"))
	g.GET("/projects/:projectId", wrapper.GetProject, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.PATCH("/projects/:projectId", wrapper.UpdateProject, auth.RolesRequired("EDITOR", "ADMIN"))
	g.GET("/projects/:projectId/export", wrapper.ExportProjectArtifacts, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.GET("/projects/:projectId/usages", wrapper.ListProjectUsages, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.POST("/projects/:projectId/usages", wrapper.CreateProjectUsage, auth.RolesRequired("EDITOR", "ADMIN"))
	g.DELETE("/projects/:projectId/usages/:usageId", wrapper.DeleteProjectUsage, auth.RolesRequired("EDITOR", "ADMIN"))
	g.PATCH("/projects/:projectId/usages/:usageId", wrapper.UpdateProjectUsage, auth.RolesRequired("EDITOR", "ADMIN"))
	g.PATCH("/projects/:projectId/usages/:usageId/scope", wrapper.UpdateProjectUsageScope, auth.RolesRequired("EDITOR", "ADMIN"))
	g.GET("/scope/policy", wrapper.GetScopePolicy, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.PATCH("/scope/policy", wrapper.UpdateScopePolicy, auth.RolesRequired("ADMIN"))
	g.GET("/tags", wrapper.ListTags, auth.RolesRequired("VIEWER", "EDITOR", "ADMIN"))
	g.POST("/tags", wrapper.CreateTag, auth.RolesRequired("EDITOR", "ADMIN"))
	g.DELETE("/tags/:tagId", wrapper.DeleteTag, auth.RolesRequired("EDITOR", "ADMIN"))
	g.GET("/users", wrapper.ListUsers, auth.RolesRequired("ADMIN"))
	g.POST("/users", wrapper.CreateUser, auth.RolesRequired("ADMIN"))
	g.DELETE("/users/:userId", wrapper.DeleteUser, auth.RolesRequired("ADMIN"))
	g.GET("/users/:userId", wrapper.GetUser, auth.RolesRequired("ADMIN"))
	g.PATCH("/users/:userId", wrapper.UpdateUser, auth.RolesRequired("ADMIN"))
}
