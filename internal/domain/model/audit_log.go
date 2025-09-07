package model

import "github.com/ramsesyok/mosscat/pkg/dbtime"

// AuditLog は監査ログの 1 レコードを表す。
type AuditLog struct {
	ID         string
	EntityType string
	EntityID   string
	Action     string
	UserName   string
	Summary    *string
	CreatedAt  dbtime.DBTime
}
