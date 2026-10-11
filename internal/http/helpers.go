package httpx

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"familytree/internal/domain"
	"familytree/internal/http/middleware"
)

// csrfToken reads the CSRF cookie so forms can embed the token.
func csrfToken(r *http.Request) string {
	if c, err := r.Cookie(middleware.CSRFCookie); err == nil {
		return c.Value
	}
	return ""
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func queryInt64(r *http.Request, key string) int64 {
	v, _ := strconv.ParseInt(r.URL.Query().Get(key), 10, 64)
	return v
}

func queryInt(r *http.Request, key string) int {
	v, _ := strconv.Atoi(r.URL.Query().Get(key))
	return v
}

// ---- form field helpers -------------------------------------------------

func formString(r *http.Request, key string) string {
	return strings.TrimSpace(r.FormValue(key))
}

func formOptional(r *http.Request, key string) *string {
	v := strings.TrimSpace(r.FormValue(key))
	if v == "" {
		return nil
	}
	return &v
}

func formIntPtr(r *http.Request, key string) (*int, error) {
	v := strings.TrimSpace(r.FormValue(key))
	if v == "" {
		return nil, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return nil, domain.Invalid("Giá trị số không hợp lệ")
	}
	return &n, nil
}

func formIDPtr(r *http.Request, key string) (*int64, error) {
	v := strings.TrimSpace(r.FormValue(key))
	if v == "" || v == "0" {
		return nil, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, domain.Invalid("Lựa chọn không hợp lệ")
	}
	return &n, nil
}

func personInputFromForm(r *http.Request) (domain.PersonInput, error) {
	gender, err := strconv.Atoi(r.FormValue("gender"))
	if err != nil {
		return domain.PersonInput{}, domain.Invalid("Giới tính không hợp lệ")
	}
	birthYear, err := formIntPtr(r, "birth_year")
	if err != nil {
		return domain.PersonInput{}, err
	}
	deathYear, err := formIntPtr(r, "death_year")
	if err != nil {
		return domain.PersonInput{}, err
	}
	fatherID, err := formIDPtr(r, "father_id")
	if err != nil {
		return domain.PersonInput{}, err
	}
	motherID, err := formIDPtr(r, "mother_id")
	if err != nil {
		return domain.PersonInput{}, err
	}
	clanID, err := formIDPtr(r, "clan_id")
	if err != nil {
		return domain.PersonInput{}, err
	}

	return domain.PersonInput{
		FullName:       formString(r, "full_name"),
		Gender:         domain.Gender(gender),
		BirthYear:      birthYear,
		BirthDateSolar: formOptional(r, "birth_date_solar"),
		BirthDateLunar: formOptional(r, "birth_date_lunar"),
		IsAlive:        r.FormValue("is_alive") != "",
		DeathYear:      deathYear,
		DeathDateSolar: formOptional(r, "death_date_solar"),
		DeathDateLunar: formOptional(r, "death_date_lunar"),
		FatherID:       fatherID,
		MotherID:       motherID,
		ClanID:         clanID,
		Address:        formOptional(r, "address"),
		Phone:          formOptional(r, "phone"),
		Occupation:     formOptional(r, "occupation"),
		AvatarURL:      formOptional(r, "avatar_url"),
		GraveLocation:  formOptional(r, "grave_location"),
		Note:           formOptional(r, "note"),
	}, nil
}

func clanInputFromForm(r *http.Request) (domain.ClanInput, error) {
	parentID, err := formIDPtr(r, "parent_clan_id")
	if err != nil {
		return domain.ClanInput{}, err
	}
	rootID, err := formIDPtr(r, "root_person_id")
	if err != nil {
		return domain.ClanInput{}, err
	}
	return domain.ClanInput{
		Name:         formString(r, "name"),
		Description:  formOptional(r, "description"),
		ParentClanID: parentID,
		RootPersonID: rootID,
	}, nil
}

// formValuesFromPerson pre-fills the form for editing.
func formValuesFromPerson(p domain.Person) personFormValues {
	return personFormValues{
		FullName:       p.FullName,
		Gender:         strconv.Itoa(int(p.Gender)),
		BirthYear:      deref(p.BirthYear),
		BirthDateSolar: deref(p.BirthDateSolar),
		BirthDateLunar: deref(p.BirthDateLunar),
		IsAlive:        p.IsAlive,
		DeathYear:      deref(p.DeathYear),
		DeathDateSolar: deref(p.DeathDateSolar),
		DeathDateLunar: deref(p.DeathDateLunar),
		FatherID:       deref(p.FatherID),
		MotherID:       deref(p.MotherID),
		ClanID:         deref(p.ClanID),
		Address:        deref(p.Address),
		Phone:          deref(p.Phone),
		Occupation:     deref(p.Occupation),
		AvatarURL:      deref(p.AvatarURL),
		GraveLocation:  deref(p.GraveLocation),
		Note:           deref(p.Note),
	}
}

// formValuesFromRequest keeps the user's input after a validation error.
func formValuesFromRequest(r *http.Request) personFormValues {
	return personFormValues{
		FullName:       r.FormValue("full_name"),
		Gender:         r.FormValue("gender"),
		BirthYear:      r.FormValue("birth_year"),
		BirthDateSolar: r.FormValue("birth_date_solar"),
		BirthDateLunar: r.FormValue("birth_date_lunar"),
		IsAlive:        r.FormValue("is_alive") != "",
		DeathYear:      r.FormValue("death_year"),
		DeathDateSolar: r.FormValue("death_date_solar"),
		DeathDateLunar: r.FormValue("death_date_lunar"),
		FatherID:       r.FormValue("father_id"),
		MotherID:       r.FormValue("mother_id"),
		ClanID:         r.FormValue("clan_id"),
		Address:        r.FormValue("address"),
		Phone:          r.FormValue("phone"),
		Occupation:     r.FormValue("occupation"),
		AvatarURL:      r.FormValue("avatar_url"),
		GraveLocation:  r.FormValue("grave_location"),
		Note:           r.FormValue("note"),
	}
}

// ---- responses ----------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeAPIError maps a domain error to a JSON response without leaking internals.
func writeAPIError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	case errors.Is(err, domain.ErrValidation):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrConflict):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "already exists"})
	default:
		log.Printf("internal error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}

// fail renders a friendly HTML error page. Validation/not-found/conflict
// messages are shown; anything else is logged and hidden.
func (s *Server) fail(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	message := "Đã xảy ra lỗi, vui lòng thử lại"
	switch {
	case errors.Is(err, domain.ErrValidation):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, domain.ErrNotFound):
		status, message = http.StatusNotFound, "Không tìm thấy dữ liệu"
	case errors.Is(err, domain.ErrConflict):
		status, message = http.StatusConflict, "Dữ liệu đã tồn tại"
	default:
		log.Printf("internal error: %v", err)
	}
	data := errorPageData{pageBase: s.base(r, "Có lỗi"), Message: message}
	s.render.Render(w, status, "error", data)
}

// safeMessage returns a message that is safe to put in a redirect.
func safeMessage(err error) string {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return err.Error()
	case errors.Is(err, domain.ErrConflict):
		return "Dữ liệu đã tồn tại"
	case errors.Is(err, domain.ErrNotFound):
		return "Không tìm thấy dữ liệu"
	default:
		log.Printf("internal error: %v", err)
		return "Đã xảy ra lỗi, vui lòng thử lại"
	}
}

func redirectMsg(w http.ResponseWriter, r *http.Request, path, msg string) {
	if msg != "" {
		sep := "?"
		if strings.Contains(path, "?") {
			sep = "&"
		}
		path += sep + "msg=" + url.QueryEscape(msg)
	}
	http.Redirect(w, r, path, http.StatusSeeOther)
}
