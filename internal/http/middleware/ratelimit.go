package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimit applies a per-client-IP token bucket so a single client cannot
// flood the server (and so login cannot be brute-forced too quickly).
func RateLimit(rps float64, burst int) Middleware {
	limiter := newIPLimiter(rate.Limit(rps), burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.allow(clientIP(r)) {
				http.Error(w, "too many requests, please slow down", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type ipLimiter struct {
	mu      sync.Mutex
	limit   rate.Limit
	burst   int
	clients map[string]*clientBucket
}

type clientBucket struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func newIPLimiter(limit rate.Limit, burst int) *ipLimiter {
	l := &ipLimiter{limit: limit, burst: burst, clients: map[string]*clientBucket{}}
	go l.cleanup()
	return l
}

func (l *ipLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.clients[ip]
	if !ok {
		b = &clientBucket{limiter: rate.NewLimiter(l.limit, l.burst)}
		l.clients[ip] = b
	}
	b.lastSeen = time.Now()
	return b.limiter.Allow()
}

func (l *ipLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		l.mu.Lock()
		for ip, b := range l.clients {
			if time.Since(b.lastSeen) > 15*time.Minute {
				delete(l.clients, ip)
			}
		}
		l.mu.Unlock()
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
