package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
)

const csrfCookieName = "csrf_token"

func EnsureCSRFCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie(csrfCookieName); err == nil {
			next.ServeHTTP(w, r)
			return
		}

		token, err := newCSRFToken()
		if err != nil {
			http.Error(w, "failed to initialize csrf", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     csrfCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: false,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   86400,
		})

		next.ServeHTTP(w, r)
	})
}

func RequireCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(csrfCookieName)
		if err != nil || cookie.Value == "" {
			http.Error(w, "csrf token missing", http.StatusForbidden)
			return
		}

		_ = r.ParseForm()
		token := r.Form.Get("csrf_token")
		if token == "" {
			token = r.Header.Get("X-CSRF-Token")
		}
		if token == "" {
			http.Error(w, "csrf token missing", http.StatusForbidden)
			return
		}

		if subtle.ConstantTimeCompare([]byte(token), []byte(cookie.Value)) != 1 {
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
