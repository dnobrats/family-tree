package handler

import "strconv"

func parseInt64Ptr(v string) *int64 {
	if v == "" {
		return nil
	}
	x, _ := strconv.ParseInt(v, 10, 64)
	return &x
}

func parseIntPtr(v string) *int {
	if v == "" {
		return nil
	}
	x, _ := strconv.Atoi(v)
	return &x
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
