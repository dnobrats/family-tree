package handler

import (
	"net/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func AdminLoginPost(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		user := r.Form.Get("username")
		pass := r.Form.Get("password")

		var hash string
		err := db.QueryRow(
			r.Context(),
			`SELECT password_hash FROM admin_user WHERE username=$1`,
			user,
		).Scan(&hash)
		if err != nil {
			http.Error(w, "invalid login", http.StatusUnauthorized)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) != nil {
			http.Error(w, "invalid login", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: "admin",
			Path:  "/",
		})

		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}
