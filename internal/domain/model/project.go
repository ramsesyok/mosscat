package model

import "github.com/ramsesyok/mosscat/pkg/dbtime"

// Project はデリバリーユニットのプロジェクトを表すモデル。
type Project struct {
	ID            string
	ProjectCode   string
	Name          string
	Department    *string
	Manager       *string
	DeliveryDate  *dbtime.DBTime
	Description   *string
	OssUsageCount int
	CreatedAt     dbtime.DBTime
	UpdatedAt     dbtime.DBTime
}
