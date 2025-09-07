package model

import "github.com/ramsesyok/mosscat/pkg/dbtime"

// OssVersion は OSS コンポーネントのバージョン情報を表す。
type OssVersion struct {
	ID                      string
	OssID                   string
	Version                 string
	ReleaseDate             *dbtime.DBTime
	LicenseExpressionRaw    *string
	LicenseConcluded        *string
	Purl                    *string
	CpeList                 []string
	HashSha256     *string
	ReviewStatus   string
	LastReviewedAt *dbtime.DBTime
	ScopeStatus    string
	CreatedAt      dbtime.DBTime
	UpdatedAt               dbtime.DBTime
}
