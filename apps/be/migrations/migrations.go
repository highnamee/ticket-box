package migrations

import "embed"

// EmbedFS contains all SQL migration files embedded into the Go binary
//
//go:embed *.sql
var EmbedFS embed.FS
