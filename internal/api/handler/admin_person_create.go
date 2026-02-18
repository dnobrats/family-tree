package handler

import (
	"net/http"
	"strconv"

	"genealogy-be/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminCreatePerson(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse form
		_ = r.ParseForm()

		// parse fields
		gender, _ := strconv.Atoi(r.Form.Get("gender"))

		// resolve IDs
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

		// execute insert
		if err := service.CreatePerson(r.Context(), db, in); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}
