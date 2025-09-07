package model

import "github.com/ramsesyok/mosscat/pkg/dbtime"

// OssComponent は OSS コンポーネントを表すドメインモデル。
type OssComponent struct {
	ID               string
	Name             string
	NormalizedName   string
	HomepageURL      *string
	RepositoryURL    *string
	Description      *string
	PrimaryLanguage  *string
	Layers           []string
	DefaultUsageRole *string
	Deprecated       bool
	CreatedAt        dbtime.DBTime
	UpdatedAt        dbtime.DBTime
	Tags             []Tag
}

// Tag はコンポーネントに付与される分類タグ。
type Tag struct {
	ID        string
	Name      string
	CreatedAt *dbtime.DBTime
}
