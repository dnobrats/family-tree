package handler

import (
	"net/http"
	"time"

	"genealogy-be/internal/auth"
)

// AdminLoginPost xử lý đăng nhập admin
func (h *Handler) AdminLoginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if err := h.service.Admin.Login(r.Context(), username, password); err != nil {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateSessionToken(username, 24*time.Hour)
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   86400,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}
