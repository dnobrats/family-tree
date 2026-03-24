package handler

import (
	"net/http"
	"strconv"
	"time"

	"genealogy-be/internal/repository"

	"github.com/go-chi/chi/v5"
)

// AdminUpdatePerson xử lý cập nhật thông tin person
func (h *Handler) AdminUpdatePerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Lỗi đọc form", http.StatusBadRequest)
			return
		}

		// Get ID from URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid person ID", http.StatusBadRequest)
			return
		}

		// Parse fields
		fullName := r.Form.Get("full_name")
		if fullName == "" {
			http.Error(w, "Họ tên là bắt buộc", http.StatusBadRequest)
			return
		}

		gender, err := strconv.Atoi(r.Form.Get("gender"))
		if err != nil {
			http.Error(w, "Giới tính không hợp lệ", http.StatusBadRequest)
			return
		}

		// Resolve IDs - chỉ báo lỗi nếu nhập tên không tìm thấy
		fatherStr := r.Form.Get("father")
		if fatherStr == "" {
			http.Error(w, "Cha là bắt buộc", http.StatusBadRequest)
			return
		}
		fatherID, err := h.service.Person.ResolvePersonID(r.Context(), fatherStr)
		if err != nil && fatherStr != "" {
			http.Error(w, "Không tìm thấy cha với tên/ID: "+fatherStr, http.StatusBadRequest)
			return
		}

		motherStr := r.Form.Get("mother")
		if motherStr == "" {
			http.Error(w, "Mẹ là bắt buộc", http.StatusBadRequest)
			return
		}
		motherID, err := h.service.Person.ResolvePersonID(r.Context(), motherStr)
		if err != nil && motherStr != "" {
			http.Error(w, "Không tìm thấy mẹ với tên/ID: "+motherStr, http.StatusBadRequest)
			return
		}

		clanStr := r.Form.Get("clan")
		if clanStr == "" {
			http.Error(w, "Chi là bắt buộc", http.StatusBadRequest)
			return
		}
		clanID, err := h.service.Person.ResolveClanID(r.Context(), clanStr)
		if err != nil && clanStr != "" {
			http.Error(w, "Không tìm thấy chi với tên/ID: "+clanStr, http.StatusBadRequest)
			return
		}

		birthDateSolar, err := parseDatePtr(r.Form.Get("birth_date_solar"))
		if err != nil {
			http.Error(w, "Ngày sinh dương lịch không hợp lệ (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		birthDateLunar, err := parseDatePtr(r.Form.Get("birth_date_lunar"))
		if err != nil {
			http.Error(w, "Ngày sinh âm lịch không hợp lệ (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		birthYear, err := parseIntPtr(r.Form.Get("birth_year"))
		if err != nil {
			http.Error(w, "Năm sinh không hợp lệ", http.StatusBadRequest)
			return
		}
		if birthYear == nil && birthDateLunar != nil {
			if t, parseErr := time.Parse("2006-01-02", *birthDateLunar); parseErr == nil {
				y := t.Year()
				birthYear = &y
			}
		}

		isAlive := r.Form.Get("is_alive") == "1"
		deathDateSolar, err := parseDatePtr(r.Form.Get("death_date_solar"))
		if err != nil {
			http.Error(w, "Ngày mất dương lịch không hợp lệ (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		deathDateLunar, err := parseDatePtr(r.Form.Get("death_date_lunar"))
		if err != nil {
			http.Error(w, "Ngày mất âm lịch không hợp lệ (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		graveLocation := parseStringPtr(r.Form.Get("grave_location"))
		if isAlive {
			deathDateSolar = nil
			deathDateLunar = nil
			graveLocation = nil
		}

		in := repository.PersonInput{
			FullName:       fullName,
			Gender:         gender,
			BirthYear:      birthYear,
			BirthDateSolar: birthDateSolar,
			BirthDateLunar: birthDateLunar,
			FatherID:       fatherID,
			MotherID:       motherID,
			ClanID:         clanID,
			IsAlive:        isAlive,
			DeathDateSolar: deathDateSolar,
			DeathDateLunar: deathDateLunar,
			Address:        parseStringPtr(r.Form.Get("address")),
			Phone:          parseStringPtr(r.Form.Get("phone")),
			Occupation:     parseStringPtr(r.Form.Get("occupation")),
			AvatarURL:      parseStringPtr(r.Form.Get("avatar_url")),
			GraveLocation:  graveLocation,
			Note:           parseStringPtr(r.Form.Get("note")),
		}

		// TODO: Thêm validation logic vào service layer

		if err := h.service.Person.UpdatePerson(r.Context(), id, in); err != nil {
			http.Error(w, "Lỗi cập nhật", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}
