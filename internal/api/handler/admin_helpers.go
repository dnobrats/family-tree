package handler

import "strconv"

func parseInt64Ptr(v string) *int64 {
	if v == "" {
		return nil
	}
	x, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil
	}
	return &x
}

func parseIntPtr(v string) *int {
	if v == "" {
		return nil
	}
	x, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	return &x
}
