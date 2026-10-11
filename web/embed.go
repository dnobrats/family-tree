// Package web embeds the HTML templates and static assets (CSS/JS) into the
// binary, so the app ships as a single executable with no files to copy.
package web

import "embed"

// FS holds templates/ and static/ so both the server and the template renderer
// can read them from memory.
//
//go:embed templates static
var FS embed.FS
