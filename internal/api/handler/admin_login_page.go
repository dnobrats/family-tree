package handler

import "net/http"

// AdminLoginPage hiển thị trang đăng nhập admin
func (h *Handler) AdminLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		csrfToken := csrfTokenFromRequest(r)
		w.Write([]byte(`
<!doctype html>
<html>
<head><title>Admin Login</title></head>
<body>
<h2>Đăng nhập quản trị</h2>
<form method="post" action="/admin/login">
  <input type="hidden" name="csrf_token" value="` + csrfToken + `"/>
  <label>User:</label><br/>
  <input name="username"/><br/><br/>
  <label>Password:</label><br/>
  <input type="password" name="password"/><br/><br/>
  <button>Login</button>
</form>
</body>
</html>
		`))
	}
}
