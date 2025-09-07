//go:generate go tool oapi-codegen -config cfg.yaml internal/api/openapi.yaml
package main

import (
	"context"
	"flag"
	"log"
	"net"
	"runtime"

	"github.com/getkin/kin-openapi/openapi3filter"
	apirouter "github.com/ramsesyok/mosscat/internal/api"
	gen "github.com/ramsesyok/mosscat/internal/api/gen"
	"github.com/ramsesyok/mosscat/internal/api/handler"
	"github.com/ramsesyok/mosscat/internal/config"
	infradb "github.com/ramsesyok/mosscat/internal/infra/db"
	"github.com/ramsesyok/mosscat/internal/infra/migration"
	infrarepo "github.com/ramsesyok/mosscat/internal/infra/repository"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
)

const serviceName = "mosscat"

func runServer(host, port, dsn string, origins []string) error {
	// OASテンプレートの読み込み
	swagger, err := gen.GetSwagger()
	if err != nil {
		return err
	}
	// ホスト名での検証は行わない
	swagger.Servers = nil

	dbConn, err := infradb.Open(dsn)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	if err := migration.Apply(dbConn.DB, dsn); err != nil {
		return err
	}

	h := handler.Handler{
		AuditRepo:             &infrarepo.AuditLogRepository{DB: dbConn.DB},
		ScopePolicyRepo:       &infrarepo.ScopePolicyRepository{DB: dbConn.DB},
		OssComponentRepo:      &infrarepo.OssComponentRepository{DB: dbConn.DB},
		OssComponentLayerRepo: &infrarepo.OssComponentLayerRepository{DB: dbConn.DB},
		OssComponentTagRepo:   &infrarepo.OssComponentTagRepository{DB: dbConn.DB},
		TagRepo:               &infrarepo.TagRepository{DB: dbConn.DB},
		OssVersionRepo:        &infrarepo.OssVersionRepository{DB: dbConn.DB},
		ProjectRepo:           &infrarepo.ProjectRepository{DB: dbConn.DB},
		ProjectUsageRepo:      &infrarepo.ProjectUsageRepository{DB: dbConn.DB},
		UserRepo:              &infrarepo.UserRepository{DB: dbConn.DB},
	}

	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: origins,
	}))
	
	// 静的ファイルの配信設定
	e.Static("/", "dist")

	// APIグループを作成してOpenAPI検証を適用
	apiGroup := e.Group("/api")
	// OASテンプレートで指定したスキーマによる検証を行う
	// 認証は別ミドルウェアで行うため、バリデータ側ではセキュリティチェックをスキップする
	apiGroup.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, _ *openapi3filter.AuthenticationInput) error {
				return nil
			},
		},
	}))

	apirouter.RegisterRoutes(e, &h)
	return e.Start(net.JoinHostPort(host, port))
}

func main() {
	cfgPath := flag.String("config", "", "config file path")
	svcFlag := flag.String("service", "", "windows service control (install|uninstall)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if runtime.GOOS == "windows" {
		switch *svcFlag {
		case "install":
			if err := installService(serviceName, "OSS Catalog service"); err != nil {
				log.Fatalf("install failed: %v", err)
			}
			return
		case "uninstall":
			if err := removeService(serviceName); err != nil {
				log.Fatalf("uninstall failed: %v", err)
			}
			return
		}

		isSvc, err := isWindowsService()
		if err == nil && isSvc {
			if err := runService(serviceName, cfg.Server.Host, cfg.Server.Port, cfg.DB.DSN, cfg.Server.AllowedOrigins); err != nil {
				log.Fatalf("service run failed: %v", err)
			}
			return
		}
	}

	if err := runServer(cfg.Server.Host, cfg.Server.Port, cfg.DB.DSN, cfg.Server.AllowedOrigins); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
