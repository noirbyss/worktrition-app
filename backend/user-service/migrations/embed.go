package migrations

import "embed"

// Files contains embedded SQL migrations used during service startup.
//
//go:embed *.up.sql
var Files embed.FS
