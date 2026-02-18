package handler

import (
	"embed"
	"net/http"
)

//go:embed static/docs.html
var docsFS embed.FS

func DocsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := docsFS.ReadFile("static/docs.html")
		if err != nil {
			http.Error(w, "docs not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	}
}
