package api

import (
	"net/http"

	"genealogy-be/internal/api/handler"
	"genealogy-be/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	// ================= PUBLIC =================
	r.Get("/docs", handler.DocsHandler())
	r.Get("/tree", handler.TreePageHandler())

	// ================= API (READ ONLY) =================
	r.Route("/api", func(r chi.Router) {
		r.Get("/tree", handler.TreeHandler(db))
		r.Get("/clans/{id}/tree", handler.ClanTreeHandler(db))
		r.Get("/persons/{id}", handler.PersonHandler(db))
	})

	// ================= ADMIN AUTH =================
	r.Get("/admin/login", handler.AdminLoginPage())
	r.Post("/admin/login", handler.AdminLoginPost(db))

	// ================= ADMIN (PROTECTED) =================
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.RequireLogin)

		r.Get("/", handler.AdminHome())

		// Manual input (CMS)
		r.Get("/persons/new", handler.AdminNewPerson())
		r.Post("/persons/new", handler.AdminSavePerson(db, false))
		r.Post("/persons/{id}", handler.AdminSavePerson(db, true))

		// (sau này)
		// r.Get("/persons", handler.AdminListPersons())
	})

	return r
}
