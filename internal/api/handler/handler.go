package handler

// handler.go - ハンドラ共通構造体
// (エンドポイント定義はなし)

import (
	domrepo "github.com/ramsesyok/mosscat/internal/domain/repository"
)

type Handler struct {
	AuditRepo             domrepo.AuditLogRepository
	ScopePolicyRepo       domrepo.ScopePolicyRepository
	OssComponentRepo      domrepo.OssComponentRepository
	OssComponentLayerRepo domrepo.OssComponentLayerRepository
	OssComponentTagRepo   domrepo.OssComponentTagRepository
	TagRepo               domrepo.TagRepository
	OssVersionRepo        domrepo.OssVersionRepository
	ProjectRepo           domrepo.ProjectRepository
	ProjectUsageRepo      domrepo.ProjectUsageRepository
	UserRepo              domrepo.UserRepository
}
