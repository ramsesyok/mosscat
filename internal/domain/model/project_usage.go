package model

import "github.com/ramsesyok/mosscat/pkg/dbtime"

// ProjectUsage はプロジェクト内での OSS 利用状況を表す。
type ProjectUsage struct {
	ID                      string
	ProjectID               string
	OssID                   string
	OssVersionID            string
	UsageRole               string
	ScopeStatus             string
	InclusionNote           *string
	DirectDependency        bool
	AddedAt                 dbtime.DBTime
	EvaluatedAt             *dbtime.DBTime
	EvaluatedBy             *string
	Modified                bool
	ModificationDescription *string
	SupplierType            *string
	ForkOriginURL           *string
}

// ProjectUsageFilter は repository パッケージで定義される。
