package handler

import (
	"net/http"
	"strconv"

	"genealogy-be/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminUpdatePerson(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse form
		_ = r.ParseForm()

		// id from URL
		idStr := chi.URLParam(r, "id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		// parse fields
		gender, _ := strconv.Atoi(r.Form.Get("gender"))
		fatherID, err := service.ResolvePersonID(r.Context(), db, r.Form.Get("father"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		motherID, err := service.ResolvePersonID(r.Context(), db, r.Form.Get("mother"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		clanID, err := service.ResolveClanID(r.Context(), db, r.Form.Get("clan"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		in := service.PersonInput{
			FullName:  r.Form.Get("full_name"),
			Gender:    gender,
			BirthYear: parseIntPtr(r.Form.Get("birth_year")),
			FatherID:  fatherID,
			MotherID:  motherID,
			ClanID:    clanID,
			IsAlive:   r.Form.Get("is_alive") == "1",
		}

		// validate ages
		if err := service.ValidateParentAge(r.Context(), db, in.BirthYear, in.FatherID, "Cha"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := service.ValidateParentAge(r.Context(), db, in.BirthYear, in.MotherID, "Mẹ"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// execute update
		if err := service.UpdatePerson(r.Context(), db, id, in); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}
