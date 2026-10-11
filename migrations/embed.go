// Package migrations embeds the numbered SQL schema files so the server can
// apply them automatically at startup (no external migration tool needed).
//
// To change the schema: add a new file named e.g. `0002_add_photo.sql`.
// Never edit an already-applied file — the runner records versions and skips
// files it has already executed.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
