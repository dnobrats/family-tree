package handler

import "net/http"

func AdminLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
<!doctype html>
<html>
<head><title>Admin Login</title></head>
<body>
<h2>Đăng nhập quản trị</h2>
<form method="post" action="/admin/login">
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
