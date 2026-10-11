package httpx

import (
	"net/http"

	"familytree/internal/domain"
)

// apiTree returns the nested family tree.
func (s *Server) apiTree(w http.ResponseWriter, r *http.Request) {
	node, err := s.svc.GetTree(r.Context(), queryInt64(r, "root"), queryInt(r, "maxGen"))
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, domain.TreeResponse{Root: node})
}

// apiClanTree returns the tree rooted at a clan's head person.
func (s *Server) apiClanTree(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid clan id"})
		return
	}
	clan, node, err := s.svc.GetClanTree(r.Context(), id, queryInt(r, "maxGen"))
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, domain.ClanTreeResponse{Clan: clan, Root: node})
}

// apiPerson returns a person's full profile.
func (s *Server) apiPerson(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid person id"})
		return
	}
	detail, err := s.svc.GetPersonDetail(r.Context(), id)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

// apiPersonsSearch powers the admin autocomplete: GET /api/persons?q=...
func (s *Server) apiPersonsSearch(w http.ResponseWriter, r *http.Request) {
	people, err := s.svc.SearchPersons(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		writeAPIError(w, err)
		return
	}
	if people == nil {
		people = []domain.PersonBrief{}
	}
	writeJSON(w, http.StatusOK, people)
}

// apiClans lists all clans.
func (s *Server) apiClans(w http.ResponseWriter, r *http.Request) {
	clans, err := s.svc.ListClans(r.Context())
	if err != nil {
		writeAPIError(w, err)
		return
	}
	if clans == nil {
		clans = []domain.Clan{}
	}
	writeJSON(w, http.StatusOK, clans)
}
