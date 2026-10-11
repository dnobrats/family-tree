package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
)

// CSRF uses the double-submit cookie pattern: a random token is stored in a
// readable cookie and must be echoed back in a form field or header.
const (
	CSRFCookie = "csrf_token"
	CSRFField  = "csrf_token"
	CSRFHeader = "X-CSRF-Token"
)

// EnsureCSRFCookie sets a CSRF cookie when the request has none. The new cookie
// is also added to the request so templates rendered in the SAME request can
// read the token (otherwise the first form load would embed an empty token).
func EnsureCSRFCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie(CSRFCookie); err != nil {
			if token, err := newCSRFToken(); err == nil {
				cookie := &http.Cookie{
					Name:     CSRFCookie,
					Value:    token,
					Path:     "/",
					HttpOnly: false, // must be readable by the double-submit JS
					SameSite: http.SameSiteLaxMode,
					MaxAge:   12 * 3600,
				}
				http.SetCookie(w, cookie)
				r.AddCookie(cookie)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RequireCSRF validates the token on state-changing methods.
func RequireCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(CSRFCookie)
		if err != nil || cookie.Value == "" {
			http.Error(w, "csrf token missing", http.StatusForbidden)
			return
		}

		_ = r.ParseForm()
		token := r.Form.Get(CSRFField)
		if token == "" {
			token = r.Header.Get(CSRFHeader)
		}
		if token == "" || subtle.ConstantTimeCompare([]byte(token), []byte(cookie.Value)) != 1 {
			http.Error(w, "csrf token invalid", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func newCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
