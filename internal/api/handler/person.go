package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// PersonHandler trả về chi tiết thông tin 1 người
func (h *Handler) PersonHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid person ID")
			return
		}

		resp, err := h.service.Genealogy.GetPersonDetail(r.Context(), id)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, resp)
	}
}
