package api

import (
	"net/http"

	"genealogy-be/internal/api/handler"
	"genealogy-be/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	// ================= LOGGING =================
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

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

		// CREATE — phải đăng ký TRƯỚC route {id} để tránh "new" bị match vào {id}
		r.Get("/persons/new", handler.AdminNewPerson())
		r.Post("/persons/new", handler.AdminCreatePerson(db))

		// UPDATE
		r.Get("/persons/{id}", handler.AdminEditPerson(db))
		r.Post("/persons/{id}", handler.AdminUpdatePerson(db))
	})

	return r
}