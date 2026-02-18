package handler

import (
	"net/http"
	"strconv"

	"genealogy-be/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminSavePerson(db *pgxpool.Pool, update bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()

		gender, _ := strconv.Atoi(r.Form.Get("gender"))

		in := service.PersonInput{
			FullName:  r.Form.Get("full_name"),
			Gender:    gender,
			BirthYear: parseIntPtr(r.Form.Get("birth_year")),
			FatherID:  parseInt64Ptr(r.Form.Get("father_id")),
			MotherID:  parseInt64Ptr(r.Form.Get("mother_id")),
			ClanID:    parseInt64Ptr(r.Form.Get("clan_id")),
			IsAlive:   r.Form.Get("is_alive") == "1",
		}

		if update {
			idStr := chi.URLParam(r, "id")
			id, _ := strconv.ParseInt(idStr, 10, 64)
			_ = service.UpdatePerson(r.Context(), db, id, in)
		} else {
			_ = service.CreatePerson(r.Context(), db, in)
		}

		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}
