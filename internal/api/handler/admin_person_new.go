package handler

import "net/http"

func AdminNewPerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
<!doctype html>
<html>
<head><title>Thêm người mới</title></head>
<body>
<h2>Thêm người mới</h2>

<form method="post" action="/admin/persons/new">
  <label>Họ tên</label><br/>
  <input name="full_name" required/><br/><br/>

  <label>Giới tính</label><br/>
  <select name="gender">
    <option value="1">Nam</option>
    <option value="2">Nữ</option>
  </select><br/><br/>

  <label>Năm sinh</label><br/>
  <input name="birth_year"/><br/><br/>

  <label>Cha (ID)</label><br/>
  <input name="father_id"/><br/><br/>

  <label>Mẹ (ID)</label><br/>
  <input name="mother_id"/><br/><br/>

  <label>Chi (ID)</label><br/>
  <input name="clan_id"/><br/><br/>

  <label>
    <input type="checkbox" name="is_alive" value="1" checked/>
    Còn sống
  </label><br/><br/>

  <button>Lưu</button>
</form>

</body>
</html>
		`))
	}
}
