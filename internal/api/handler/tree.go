package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"genealogy-be/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TreeHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		root := r.URL.Query().Get("root")
		if root == "" {
			http.Error(w, "missing root", http.StatusBadRequest)
			return
		}

		rootID, _ := strconv.ParseInt(root, 10, 64)
		resp, err := service.GetTree(r.Context(), db, rootID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
