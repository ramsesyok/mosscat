package repository

import (
	"database/sql"
	"strings"

	"github.com/ramsesyok/mosscat/pkg/dbtime"
)

// whereClause は条件句の配列から WHERE 句文字列を生成する。
func whereClause(wheres []string) string {
	if len(wheres) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(wheres, " AND ")
}

// strPtr は NullString から *string を生成する。
func strPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// timePtr は NullTime から *time.Time を生成する。
func timePtr(nt sql.NullTime) *dbtime.DBTime {
	if nt.Valid {
		t := dbtime.DBTime{Time: nt.Time}
		return &t
	}
	return nil
}
