// Package httpx is the transport layer: it wires HTTP routes to the service and
// renders HTML pages / JSON. All routes live in Routes(); everything else is a
// handler or a small helper.
package httpx

import (
	"io/fs"
	"net/http"
	"time"

	"familytree/internal/auth"
	"familytree/internal/config"
	"familytree/internal/http/middleware"
	"familytree/internal/service"
	"familytree/web"
)

// Server bundles the dependencies shared by every handler.
type Server struct {
	svc    *service.Service
	auth   *auth.Manager
	cfg    *config.Config
	render *Renderer
}

// NewServer builds the HTTP server.
func NewServer(svc *service.Service, authMgr *auth.Manager, cfg *config.Config, render *Renderer) *Server {
	return &Server{svc: svc, auth: authMgr, cfg: cfg, render: render}
}

// Routes registers every route and returns the fully wrapped handler.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// ---- static assets (embedded) ----
	staticFS, _ := fs.Sub(web.FS, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	// ---- public pages ----
	mux.HandleFunc("GET /{$}", s.handleHome)
	mux.HandleFunc("GET /tree", s.handleTreePage)
	mux.HandleFunc("GET /persons/{id}", s.handlePersonPage)
	mux.HandleFunc("GET /healthz", s.handleHealth)

	// ---- public JSON API ----
	mux.HandleFunc("GET /api/tree", s.apiTree)
	mux.HandleFunc("GET /api/clans", s.apiClans)
	mux.HandleFunc("GET /api/clans/{id}/tree", s.apiClanTree)
	mux.HandleFunc("GET /api/persons", s.apiPersonsSearch)
	mux.HandleFunc("GET /api/persons/{id}", s.apiPerson)

	// ---- admin: authentication ----
	mux.HandleFunc("GET /admin/login", s.handleLoginPage)
	mux.Handle("POST /admin/login", middleware.RateLimit(0.5, 5)(
		middleware.RequireCSRF(http.HandlerFunc(s.handleLoginPost))))
	mux.Handle("POST /admin/logout", middleware.RequireCSRF(http.HandlerFunc(s.handleLogout)))
	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/", http.StatusSeeOther)
	})

	// ---- admin: protected pages ----
	protect := func(h http.HandlerFunc) http.Handler {
		return middleware.RequireLogin(s.auth)(middleware.RequireCSRF(h))
	}
	mux.Handle("GET /admin/{$}", protect(s.handleDashboard))

	mux.Handle("GET /admin/persons", protect(s.handlePersonList))
	mux.Handle("GET /admin/persons/new", protect(s.handlePersonNew))
	mux.Handle("POST /admin/persons/new", protect(s.handlePersonCreate))
	mux.Handle("GET /admin/persons/{id}/edit", protect(s.handlePersonEdit))
	mux.Handle("POST /admin/persons/{id}", protect(s.handlePersonUpdate))
	mux.Handle("POST /admin/persons/{id}/delete", protect(s.handlePersonDelete))
	mux.Handle("GET /admin/persons/{id}/spouses", protect(s.handleSpousePage))
	mux.Handle("POST /admin/persons/{id}/spouses", protect(s.handleSpouseAdd))
	mux.Handle("POST /admin/spouses/{id}/delete", protect(s.handleSpouseDelete))

	mux.Handle("GET /admin/clans", protect(s.handleClanList))
	mux.Handle("POST /admin/clans", protect(s.handleClanCreate))
	mux.Handle("GET /admin/clans/{id}/edit", protect(s.handleClanEdit))
	mux.Handle("POST /admin/clans/{id}", protect(s.handleClanUpdate))
	mux.Handle("POST /admin/clans/{id}/delete", protect(s.handleClanDelete))

	// Global middleware, applied outermost first.
	return middleware.Chain(mux,
		middleware.SecurityHeaders,
		middleware.RateLimit(20, 40),
		middleware.EnsureCSRFCookie,
	)
}

// base builds the shared page header values.
func (s *Server) base(r *http.Request, title string) pageBase {
	return pageBase{
		Title:     title,
		CSRFToken: csrfToken(r),
		Flash:     r.URL.Query().Get("msg"),
		Year:      time.Now().Year(),
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}
