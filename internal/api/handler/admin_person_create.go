package handler

import (
	"log"
	"net/http"
	"strconv"

	"genealogy-be/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminCreatePerson(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[AdminCreatePerson] Nhận request POST /admin/persons/new")

		// parse form
		if err := r.ParseForm(); err != nil {
			log.Printf("[AdminCreatePerson] Lỗi parse form: %v", err)
			http.Error(w, "Lỗi đọc form", http.StatusBadRequest)
			return
		}

		log.Printf("[AdminCreatePerson] Form data: %v", r.Form)

		// parse fields
		gender, _ := strconv.Atoi(r.Form.Get("gender"))

		// resolve IDs
		fatherRaw := r.Form.Get("father")
		motherRaw := r.Form.Get("mother")
		clanRaw := r.Form.Get("clan")
		log.Printf("[AdminCreatePerson] father=%q mother=%q clan=%q", fatherRaw, motherRaw, clanRaw)

		fatherID, err := service.ResolvePersonID(r.Context(), db, fatherRaw)
		if err != nil {
			log.Printf("[AdminCreatePerson] Lỗi resolve father: %v", err)
			http.Error(w, "Lỗi tìm cha: "+err.Error(), http.StatusBadRequest)
			return
		}
		motherID, err := service.ResolvePersonID(r.Context(), db, motherRaw)
		if err != nil {
			log.Printf("[AdminCreatePerson] Lỗi resolve mother: %v", err)
			http.Error(w, "Lỗi tìm mẹ: "+err.Error(), http.StatusBadRequest)
			return
		}
		clanID, err := service.ResolveClanID(r.Context(), db, clanRaw)
		if err != nil {
			log.Printf("[AdminCreatePerson] Lỗi resolve clan: %v", err)
			http.Error(w, "Lỗi tìm chi: "+err.Error(), http.StatusBadRequest)
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
		log.Printf("[AdminCreatePerson] PersonInput: %+v", in)

		// validate ages
		if err := service.ValidateParentAge(r.Context(), db, in.BirthYear, in.FatherID, "Cha"); err != nil {
			log.Printf("[AdminCreatePerson] Lỗi validate tuổi cha: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := service.ValidateParentAge(r.Context(), db, in.BirthYear, in.MotherID, "Mẹ"); err != nil {
			log.Printf("[AdminCreatePerson] Lỗi validate tuổi mẹ: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// execute insert
		if err := service.CreatePerson(r.Context(), db, in); err != nil {
			log.Printf("[AdminCreatePerson] Lỗi INSERT vào DB: %v", err)
			http.Error(w, "Lỗi lưu vào cơ sở dữ liệu: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[AdminCreatePerson] Tạo thành công người: %q", in.FullName)
		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}