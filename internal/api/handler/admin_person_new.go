package handler

import "net/http"

// AdminNewPerson hiển thị form tạo person mới
func (h *Handler) AdminNewPerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		csrfToken := csrfTokenFromRequest(r)

		success := r.URL.Query().Get("success") == "1"

		w.Write([]byte(`<!doctype html>
<html>
<head>
  <title>Thêm người mới</title>
  <style>
    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      background: #f8fafc;
      margin: 0;
      padding: 20px;
    }
    form {
      max-width: 600px;
      margin: 0 auto;
      background: white;
      padding: 30px;
      border-radius: 12px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }
    h2 {
      color: #1e293b;
      margin-top: 0;
      border-bottom: 3px solid #3b82f6;
      padding-bottom: 15px;
    }
    label {
      display: block;
      margin-bottom: 5px;
      font-weight: 600;
      color: #475569;
    }
    input, select {
      width: 100%;
      padding: 10px;
      border: 1px solid #e2e8f0;
      border-radius: 6px;
      margin-bottom: 15px;
      box-sizing: border-box;
    }
    button, .btn {
      padding: 10px 20px;
      border: none;
      border-radius: 6px;
      cursor: pointer;
      font-weight: 600;
      transition: all 0.2s;
    }
    button[type="submit"] {
      background: #3b82f6;
      color: white;
    }
    button[type="submit"]:hover {
      background: #2563eb;
    }
  </style>
</head>
<body>

`))

		if success {
			w.Write([]byte(`
<div id="success-popup" style="
  position: fixed;
  top: 20px;
  right: 20px;
  background: #4CAF50;
  color: white;
  padding: 14px 20px;
  border-radius: 6px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.2);
  z-index: 9999;
">
  ✅ Thêm thành viên thành công
</div>

<script>
setTimeout(function () {
  const el = document.getElementById("success-popup");
  if (el) el.remove();
}, 3000);

// xoá query param để F5 không hiện lại
const url = new URL(window.location);
url.searchParams.delete("success");
window.history.replaceState({}, document.title, url.pathname);
</script>
`))
		}

		w.Write([]byte(`
<h2>Thêm người mới</h2>

<form method="post" action="/admin/persons/new">
  <input type="hidden" name="csrf_token" value="` + csrfToken + `"/>
  <label>Họ tên</label><br/>
  <input name="full_name" required/><br/><br/>

  <label>Giới tính</label><br/>
  <select name="gender">
    <option value="1">Nam</option>
    <option value="2">Nữ</option>
  </select><br/><br/>

  <label>Ngày sinh dương lịch (birth_date_solar)</label><br/>
  <input type="date" name="birth_date_solar"/><br/><br/>

  <label>Ngày sinh âm lịch (birth_date_lunar)</label><br/>
  <input type="date" name="birth_date_lunar"/><br/><br/>

  <label>Năm sinh (fallback)</label><br/>
  <input name="birth_year" type="number"/><br/><br/>

  <label>Cha (ID hoặc tên)</label><br/>
  <input name="father" required/><br/><br/>

  <label>Mẹ (ID hoặc tên)</label><br/>
  <input name="mother" required/><br/><br/>

  <label>Chi (ID hoặc tên trưởng chi)</label><br/>
  <input name="clan" required/><br/><br/>

  <label>
    <input id="is_alive_checkbox" type="checkbox" name="is_alive" value="1" checked/>
    Còn sống
  </label><br/><br/>

  <div id="death_fields" style="display:none;">
    <label>Ngày mất dương lịch (death_date_solar)</label><br/>
    <input type="date" name="death_date_solar"/><br/><br/>

    <label>Ngày mất âm lịch (death_date_lunar)</label><br/>
    <input type="date" name="death_date_lunar"/><br/><br/>

    <label>Nơi an táng (grave_location)</label><br/>
    <input name="grave_location"/><br/><br/>
  </div>

  <label>Số điện thoại</label><br/>
  <input name="phone"/><br/><br/>

  <label>Nghề nghiệp</label><br/>
  <input name="occupation"/><br/><br/>

  <label>Ảnh đại diện (URL)</label><br/>
  <input name="avatar_url"/><br/><br/>

  <label>Địa chỉ</label><br/>
  <input name="address"/><br/><br/>

  <label>Ghi chú</label><br/>
  <input name="note"/><br/><br/>

  <button type="submit" class="btn btn-primary">Lưu</button>
  <a href="/admin" style="
    display: inline-block;
    padding: 10px 20px;
    background: #64748b;
    color: white;
    text-decoration: none;
    border-radius: 6px;
    margin-left: 10px;
    font-weight: 600;
  ">← Quay về Admin</a>
</form>

<script>
  (function() {
    const alive = document.getElementById('is_alive_checkbox');
    const deathFields = document.getElementById('death_fields');
    function sync() {
      deathFields.style.display = alive.checked ? 'none' : 'block';
    }
    alive.addEventListener('change', sync);
    sync();
  })();
</script>

</body>
</html>
`))
	}
}
