package httpx

import (
	"net/http"

	"familytree/internal/domain"
)

// handleHome sends visitors to the tree.
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tree", http.StatusSeeOther)
}

// handleTreePage renders the interactive family tree page.
func (s *Server) handleTreePage(w http.ResponseWriter, r *http.Request) {
	clans, err := s.svc.ListClans(r.Context())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	data := treePageData{
		pageBase: s.base(r, "Cây gia phả"),
		Clans:    clans,
		RootID:   queryInt64(r, "root"),
		MaxGen:   queryInt(r, "maxGen"),
	}
	s.render.Render(w, http.StatusOK, "tree", data)
}

// handlePersonPage renders a single person's profile.
func (s *Server) handlePersonPage(w http.ResponseWriter, r *http.Request) {
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
	data := personPageData{
		pageBase: s.base(r, detail.Person.FullName),
		Detail:   detail,
	}
	s.render.Render(w, http.StatusOK, "person", data)
}
