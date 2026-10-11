package middleware

import (
	"net/http"

	"familytree/internal/auth"
)

// SessionCookie is the name of the session cookie.
const SessionCookie = "session"

// RequireLogin redirects to the login page unless a valid session is present.
func RequireLogin(mgr *auth.Manager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookie)
			if err != nil {
				redirectToLogin(w, r)
				return
			}
			if _, ok := mgr.Verify(cookie.Value); !ok {
				// Clear the bad cookie so the browser stops sending it.
				http.SetCookie(w, &http.Cookie{Name: SessionCookie, Value: "", Path: "/", MaxAge: -1})
				redirectToLogin(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
