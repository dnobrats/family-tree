package handler

func AdminSavePerson(db *pgxpool.Pool, update bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		in := service.PersonInput{
			FullName: r.Form.Get("full_name"),
			Gender:   atoi(r.Form.Get("gender")),
			BirthYear: parseIntPtr(r.Form.Get("birth_year")),
			FatherID: parseInt64Ptr(r.Form.Get("father_id")),
			MotherID: parseInt64Ptr(r.Form.Get("mother_id")),
			ClanID:   parseInt64Ptr(r.Form.Get("clan_id")),
			IsAlive:  r.Form.Get("is_alive") == "1",
		}

		if update {
			id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			_ = service.UpdatePerson(r.Context(), db, id, in)
		} else {
			_ = service.CreatePerson(r.Context(), db, in)
		}

		http.Redirect(w, r, "/admin/persons", http.StatusFound)
	}
}
