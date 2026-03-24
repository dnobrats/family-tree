package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func parseInt64Ptr(v string) *int64 {
	if v == "" {
		return nil
	}
	x, _ := strconv.ParseInt(v, 10, 64)
	return &x
}

func parseIntPtr(v string) (*int, error) {
	if v == "" {
		return nil, nil
	}
	x, err := strconv.Atoi(v)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func selected(ok bool) string {
	if ok {
		return "selected"
	}
	return ""
}

func checked(ok bool) string {
	if ok {
		return "checked"
	}
	return ""
}

func intPtrToStr(v *int) string {
	if v == nil {
		return ""
	}
	return strconv.Itoa(*v)
}

func int64PtrToStr(v *int64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatInt(*v, 10)
}

func parseStringPtr(v string) *string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	trimmed := strings.TrimSpace(v)
	return &trimmed
}

func strPtrToStr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func parseDatePtr(v string) (*string, error) {
	if strings.TrimSpace(v) == "" {
		return nil, nil
	}
	trimmed := strings.TrimSpace(v)
	if _, err := time.Parse("2006-01-02", trimmed); err != nil {
		return nil, err
	}
	return &trimmed, nil
}

func csrfTokenFromRequest(r *http.Request) string {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}
