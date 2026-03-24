package handler

import (
	"net/http"
	"strconv"
)

// TreeHandler trả về cây gia phả từ root person
func (h *Handler) TreeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		root := r.URL.Query().Get("root")
		if root == "" {
			respondError(w, http.StatusBadRequest, "missing root parameter")
			return
		}

		rootID, err := strconv.ParseInt(root, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid root ID")
			return
		}

		resp, err := h.service.Genealogy.GetTree(r.Context(), rootID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, resp)
	}
}
