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

		fatherID, err := service.ResolvePersonID(r.Context(), db, r.Form.Get("father"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		motherID, err := service.ResolvePersonID(r.Context(), db, r.Form.Get("mother"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		clanID, err := service.ResolveClanID(r.Context(), db, r.Form.Get("clan"))
		if err != nil {
			http.Error(w, err.Error(), 400)
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

		// Validate tuổi cha
		if err := service.ValidateParentAge(
			r.Context(),
			db,
			in.BirthYear,
			in.FatherID,
			"Cha",
		); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// Validate tuổi mẹ
		if err := service.ValidateParentAge(
			r.Context(),
			db,
			in.BirthYear,
			in.MotherID,
			"Mẹ",
		); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		if update {
			idStr := chi.URLParam(r, "id")
			id, _ := strconv.ParseInt(idStr, 10, 64)

			if err := service.UpdatePerson(r.Context(), db, id, in); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			if err := service.CreatePerson(r.Context(), db, in); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		http.Redirect(w, r, "/admin", http.StatusFound)

	}
}
