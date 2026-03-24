package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// AdminDeletePerson xử lý xóa person
func (h *Handler) AdminDeletePerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid person ID", http.StatusBadRequest)
			return
		}

		if err := h.service.Person.DeletePerson(r.Context(), id); err != nil {
			http.Error(w, "Không thể xóa người này. Có thể họ có con cháu.", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
