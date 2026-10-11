// Package middleware contains the HTTP middlewares shared by all routes:
// security headers, per-IP rate limiting, session auth and CSRF protection.
package middleware

import "net/http"

// Middleware wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares so the first listed runs outermost.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
