package httpx

import (
	"net/http"
	"strconv"

	"familytree/internal/domain"
	"familytree/internal/http/middleware"
)

// ---- authentication -----------------------------------------------------

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	data := loginPageData{
		pageBase: s.base(r, "Đăng nhập"),
		Error:    r.URL.Query().Get("error"),
	}
	s.render.Render(w, http.StatusOK, "login", data)
}

func (s *Server) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	username := r.FormValue("username")
	if err := s.svc.Authenticate(r.Context(), username, r.FormValue("password")); err != nil {
		data := loginPageData{pageBase: s.base(r, "Đăng nhập"), Error: "Tên đăng nhập hoặc mật khẩu không đúng"}
		s.render.Render(w, http.StatusUnauthorized, "login", data)
		return
	}

	token, err := s.auth.Issue(username)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	http.SetCookie(w, s.sessionCookie(token))
	http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: middleware.SessionCookie, Value: "", Path: "/", MaxAge: -1})
	redirectMsg(w, r, "/admin/login", "Đã đăng xuất")
}

func (s *Server) sessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     middleware.SessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   s.cfg.CookieSecure,
		MaxAge:   int(s.auth.TTL().Seconds()),
	}
}

// ---- dashboard ----------------------------------------------------------

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	total := 0
	if page, err := s.svc.ListPersons(r.Context(), 1, 1); err == nil {
		total = page.Total
	}
	clans, _ := s.svc.ListClans(r.Context())

	data := dashboardData{
		pageBase:    s.base(r, "Quản trị"),
		PersonCount: total,
		ClanCount:   len(clans),
	}
	data.Admin = true
	s.render.Render(w, http.StatusOK, "dashboard", data)
}

// ---- people -------------------------------------------------------------

func (s *Server) handlePersonList(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.ListPersons(r.Context(), queryInt(r, "page"), 20)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	page := &personPage{
		Items:      result.Items,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
	data := personsData{pageBase: s.base(r, "Danh sách thành viên"), Page: page}
	data.Admin = true
	s.render.Render(w, http.StatusOK, "persons", data)
}

func (s *Server) handlePersonNew(w http.ResponseWriter, r *http.Request) {
	values := personFormValues{IsAlive: true, Gender: "1"}
	s.renderPersonForm(w, r, 0, values, false, "")
}

func (s *Server) handlePersonEdit(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		s.fail(w, r, domain.Invalid("ID không hợp lệ"))
		return
	}
	detail, err := s.svc.GetPersonDetail(r.Context(), id)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	s.renderPersonForm(w, r, id, formValuesFromPerson(detail.Person), true, "")
}

func (s *Server) handlePersonCreate(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	in, err := personInputFromForm(r)
	if err == nil {
		_, err = s.svc.CreatePerson(r.Context(), in)
	}
	if err != nil {
		s.renderPersonForm(w, r, 0, formValuesFromRequest(r), false, safeMessage(err))
		return
	}
	redirectMsg(w, r, "/admin/persons", "Đã thêm thành viên mới")
}

func (s *Server) handlePersonUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		s.fail(w, r, domain.Invalid("ID không hợp lệ"))
		return
	}
	_ = r.ParseForm()

	in, err := personInputFromForm(r)
	if err == nil {
		err = s.svc.UpdatePerson(r.Context(), id, in)
	}
	if err != nil {
		s.renderPersonForm(w, r, id, formValuesFromRequest(r), true, safeMessage(err))
		return
	}
	redirectMsg(w, r, "/admin/persons", "Đã cập nhật thông tin")
}

func (s *Server) handlePersonDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}
	if err := s.svc.DeletePerson(r.Context(), id); err != nil {
		http.Error(w, safeMessage(err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) renderPersonForm(w http.ResponseWriter, r *http.Request, id int64, values personFormValues, isEdit bool, errMsg string) {
	people, _ := s.svc.ListAllPersons(r.Context())
	clans, _ := s.svc.ListClans(r.Context())

	title := "Thêm người mới"
	if isEdit {
		title = "Sửa thông tin"
	}
	status := http.StatusOK
	if errMsg != "" {
		status = http.StatusBadRequest
	}

	data := personFormData{
		pageBase: s.base(r, title),
		PersonID: id,
		Values:   values,
		People:   people,
		Clans:    clans,
		IsEdit:   isEdit,
		Error:    errMsg,
	}
	data.Admin = true
	s.render.Render(w, status, "person_form", data)
}

// ---- spouses ------------------------------------------------------------

func (s *Server) handleSpousePage(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		s.fail(w, r, domain.Invalid("ID không hợp lệ"))
		return
	}
	detail, err := s.svc.GetPersonDetail(r.Context(), id)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	people, _ := s.svc.ListAllPersons(r.Context())

	data := spousesData{
		pageBase: s.base(r, "Vợ/Chồng của "+detail.Person.FullName),
		Person:   detail.Person,
		Spouses:  detail.Spouses,
		People:   people,
	}
	data.Admin = true
	s.render.Render(w, http.StatusOK, "spouses", data)
}

func (s *Server) handleSpouseAdd(w http.ResponseWriter, r *http.Request) {
	personID, err := parseID(r)
	if err != nil {
		s.fail(w, r, domain.Invalid("ID không hợp lệ"))
		return
	}
	_ = r.ParseForm()

	partnerID, err := formIDPtr(r, "partner_id")
	if err != nil || partnerID == nil {
		s.fail(w, r, domain.Invalid("Vui lòng chọn vợ/chồng"))
		return
	}
	marriageYear, err := formIntPtr(r, "marriage_year")
	if err != nil {
		s.fail(w, r, err)
		return
	}

	in := domain.SpouseInput{
		PersonID:     personID,
		PartnerID:    *partnerID,
		MarriageYear: marriageYear,
		Note:         formOptional(r, "note"),
	}
	if err := s.svc.AddSpouse(r.Context(), in); err != nil {
		s.fail(w, r, err)
		return
	}
	redirectMsg(w, r, "/admin/persons/"+strconv.FormatInt(personID, 10)+"/spouses", "Đã thêm vợ/chồng")
}

func (s *Server) handleSpouseDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}
	if err := s.svc.DeleteSpouse(r.Context(), id); err != nil {
		http.Error(w, safeMessage(err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ---- clans --------------------------------------------------------------

func (s *Server) handleClanList(w http.ResponseWriter, r *http.Request) {
	clans, err := s.svc.ListClans(r.Context())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	people, _ := s.svc.ListAllPersons(r.Context())

	var edit *domain.Clan
	if idStr := r.URL.Query().Get("edit"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			if c, err := s.svc.GetClan(r.Context(), id); err == nil {
				edit = c
			}
		}
	}

	data := clansData{
		pageBase: s.base(r, "Quản lý chi"),
		Clans:    clans,
		People:   people,
		Edit:     edit,
		Error:    r.URL.Query().Get("error"),
	}
	data.Admin = true
	if edit != nil {
		data.EditID = edit.ID
		if edit.ParentClanID != nil {
			data.EditParentID = *edit.ParentClanID
		}
		if edit.RootPersonID != nil {
			data.EditRootID = *edit.RootPersonID
		}
	}
	s.render.Render(w, http.StatusOK, "clans", data)
}

func (s *Server) handleClanCreate(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	in, err := clanInputFromForm(r)
	if err == nil {
		_, err = s.svc.CreateClan(r.Context(), in)
	}
	if err != nil {
		redirectMsg(w, r, "/admin/clans", safeMessage(err))
		return
	}
	redirectMsg(w, r, "/admin/clans", "Đã tạo chi mới")
}

func (s *Server) handleClanEdit(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		redirectMsg(w, r, "/admin/clans", "ID không hợp lệ")
		return
	}
	http.Redirect(w, r, "/admin/clans?edit="+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

func (s *Server) handleClanUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		redirectMsg(w, r, "/admin/clans", "ID không hợp lệ")
		return
	}
	_ = r.ParseForm()

	in, err := clanInputFromForm(r)
	if err == nil {
		err = s.svc.UpdateClan(r.Context(), id, in)
	}
	if err != nil {
		redirectMsg(w, r, "/admin/clans?edit="+strconv.FormatInt(id, 10), safeMessage(err))
		return
	}
	redirectMsg(w, r, "/admin/clans", "Đã cập nhật chi")
}

func (s *Server) handleClanDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		redirectMsg(w, r, "/admin/clans", "ID không hợp lệ")
		return
	}
	if err := s.svc.DeleteClan(r.Context(), id); err != nil {
		redirectMsg(w, r, "/admin/clans", safeMessage(err))
		return
	}
	redirectMsg(w, r, "/admin/clans", "Đã xóa chi")
}
