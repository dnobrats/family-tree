package api

import (
	"net/http"

	"genealogy-be/internal/api/handler"
	"genealogy-be/internal/middleware"
	"genealogy-be/internal/repository"
	"genealogy-be/internal/service"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool) http.Handler {
	// Khởi tạo layers: Repository -> Service -> Handler
	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)

	r := chi.NewRouter()

	// ================= MIDDLEWARE =================
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// ================= PUBLIC PAGES =================
	r.Get("/docs", h.DocsHandler())
	r.Get("/tree", h.TreePageHandler())

	// ================= API (READ ONLY) =================
	r.Route("/api", func(r chi.Router) {
		r.Get("/tree", h.TreeHandler())
		r.Get("/clans/{id}/tree", h.ClanTreeHandler())
		r.Get("/persons/{id}", h.PersonHandler())
	})

	// ================= ADMIN AUTH =================
	r.With(middleware.EnsureCSRFCookie).Get("/admin/login", h.AdminLoginPage())
	r.With(middleware.RequireCSRF).Post("/admin/login", h.AdminLoginPost())

	// ================= ADMIN (PROTECTED) =================
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.RequireLogin)
		r.Use(middleware.EnsureCSRFCookie)
		r.Use(middleware.RequireCSRF)

		r.Get("/", h.AdminHome())

		// LIST - Danh sách thành viên
		r.Get("/persons", h.AdminPersonList())

		// CREATE — phải đăng ký TRƯỚC route {id} để tránh "new" bị match vào {id}
		r.Get("/persons/new", h.AdminNewPerson())
		r.Post("/persons/new", h.AdminCreatePerson())

		// SPOUSE MANAGEMENT
		r.Get("/persons/{id}/spouses", h.AdminSpouseManage())
		r.Post("/persons/{id}/spouses", h.AdminSpouseAdd())

		// UPDATE
		r.Get("/persons/{id}", h.AdminEditPerson())
		r.Post("/persons/{id}", h.AdminUpdatePerson())

		// DELETE
		r.Post("/persons/{id}/delete", h.AdminDeletePerson())
		r.Post("/spouses/{id}/delete", h.AdminSpouseDelete())
	})

	return r
}
