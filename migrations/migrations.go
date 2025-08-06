package migrations

import "embed"

// FS contains SQL migration files.
//
//go:embed *.sql
var FS embed.FS
