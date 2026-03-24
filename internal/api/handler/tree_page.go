package handler

import (
	"embed"
	"net/http"
)

//go:embed static/tree.html
var treeFS embed.FS

// TreePageHandler hiển thị trang tree visualization
func (h *Handler) TreePageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := treeFS.ReadFile("static/tree.html")
		if err != nil {
			http.Error(w, "tree page not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	}
}
