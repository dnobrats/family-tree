package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"genealogy-be/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PersonHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid person id", http.StatusBadRequest)
			return
		}

		resp, err := service.GetPersonDetail(r.Context(), db, id)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

