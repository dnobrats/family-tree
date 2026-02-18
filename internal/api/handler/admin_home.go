package handler

import "net/http"

func AdminHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
<!doctype html>
<html>
<head><title>Admin Dashboard</title></head>
<body>
<h2>Admin Dashboard</h2>
<ul>
  <li><a href="/admin/persons/new">➕ Thêm người</a></li>
  <li><a href="/tree">🌳 Xem cây gia phả</a></li>
</ul>
</body>
</html>
		`))
	}
}
