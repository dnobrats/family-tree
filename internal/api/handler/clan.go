package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ClanTreeHandler trả về cây gia phả theo chi
func (h *Handler) ClanTreeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid clan ID")
			return
		}

		resp, err := h.service.Genealogy.GetClanTree(r.Context(), id)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, resp)
	}
}
