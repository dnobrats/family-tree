package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireCSRF(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := RequireCSRF(ok)

	t.Run("rejects missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/login", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("want 403, got %d", rec.Code)
		}
	})

	t.Run("accepts matching header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/login", nil)
		req.AddCookie(&http.Cookie{Name: CSRFCookie, Value: "abc123"})
		req.Header.Set(CSRFHeader, "abc123")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("want 200, got %d", rec.Code)
		}
	})

	t.Run("rejects mismatched token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/login", nil)
		req.AddCookie(&http.Cookie{Name: CSRFCookie, Value: "abc123"})
		req.Header.Set(CSRFHeader, "different")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("want 403, got %d", rec.Code)
		}
	})

	t.Run("ignores safe methods", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tree", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("want 200, got %d", rec.Code)
		}
	})
}

func TestRateLimit_PerIP(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := RateLimit(0.0001, 2)(ok)

	call := func(ip string) int {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = ip
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}

	// Burst of 2 for the first IP.
	if call("10.0.0.1:1111") != http.StatusOK || call("10.0.0.1:2222") != http.StatusOK {
		t.Fatal("first two requests from one IP should pass")
	}
	if call("10.0.0.1:3333") != http.StatusTooManyRequests {
		t.Fatal("third request from same IP should be limited")
	}
	// A different IP has its own bucket.
	if call("10.0.0.2:1111") != http.StatusOK {
		t.Fatal("a different IP must not be affected")
	}
}
