package middleware

import "net/http"

// SecurityHeaders adds conservative security headers to every response.
//
// The CSP forbids inline scripts, so all JS must live in /static files (see
// web/static/js). Inline <style> and style="" are allowed for simplicity.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "same-origin")
		h.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"img-src 'self' data: https:; "+
				"style-src 'self' 'unsafe-inline'; "+
				"script-src 'self'; "+
				"font-src 'self' data:; "+
				"base-uri 'self'; form-action 'self'")
		next.ServeHTTP(w, r)
	})
}
