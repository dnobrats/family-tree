package handler

import (
	"html"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// AdminEditPerson hiển thị form edit person
func (h *Handler) AdminEditPerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
			return
		}

		person, err := h.service.Person.GetPersonForEdit(r.Context(), id)
		if err != nil {
			http.Error(w, "Không tìm thấy người", http.StatusNotFound)
			return
		}

		// Render form edit (đơn giản, không template engine)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		csrfToken := csrfTokenFromRequest(r)
		fullName := html.EscapeString(person.FullName)
		birthDateSolar := html.EscapeString(strPtrToStr(person.BirthDateSolar))
		birthDateLunar := html.EscapeString(strPtrToStr(person.BirthDateLunar))
		birthYear := html.EscapeString(intPtrToStr(person.BirthYear))
		father := html.EscapeString(int64PtrToStr(person.FatherID))
		mother := html.EscapeString(int64PtrToStr(person.MotherID))
		clan := html.EscapeString(int64PtrToStr(person.ClanID))
		deathDateSolar := html.EscapeString(strPtrToStr(person.DeathDateSolar))
		deathDateLunar := html.EscapeString(strPtrToStr(person.DeathDateLunar))
		graveLocation := html.EscapeString(strPtrToStr(person.GraveLocation))
		phone := html.EscapeString(strPtrToStr(person.Phone))
		occupation := html.EscapeString(strPtrToStr(person.Occupation))
		avatarURL := html.EscapeString(strPtrToStr(person.AvatarURL))
		address := html.EscapeString(strPtrToStr(person.Address))
		note := html.EscapeString(strPtrToStr(person.Note))

		w.Write([]byte(`
<!doctype html>
<html>
<head><title>Sửa người</title></head>
<body>
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
<h2>Sửa thông tin</h2>

<form method="post" action="/admin/persons/` + idStr + `">
  <input type="hidden" name="csrf_token" value="` + csrfToken + `" />
  <label>Họ tên</label><br/>
  <input name="full_name" value="` + fullName + `" required /><br/><br/>

  <label>Giới tính</label><br/>
  <select name="gender">
    <option value="1" ` + selected(person.Gender == 1) + `>Nam</option>
    <option value="2" ` + selected(person.Gender == 2) + `>Nữ</option>
  </select><br/><br/>

  <label>Ngày sinh dương lịch (birth_date_solar)</label><br/>
  <input type="date" name="birth_date_solar" value="` + birthDateSolar + `" /><br/><br/>

  <label>Ngày sinh âm lịch (birth_date_lunar)</label><br/>
  <input type="date" name="birth_date_lunar" value="` + birthDateLunar + `" /><br/><br/>

  <label>Năm sinh</label><br/>
  <input name="birth_year" value="` + birthYear + `" /><br/><br/>

  <label>Cha (ID hoặc tên)</label><br/>
  <input name="father" value="` + father + `" required /><br/><br/>

  <label>Mẹ (ID hoặc tên)</label><br/>
  <input name="mother" value="` + mother + `" required /><br/><br/>

  <label>Chi (ID hoặc trưởng chi)</label><br/>
  <input name="clan" value="` + clan + `" required /><br/><br/>

  <label>
    <input id="is_alive_checkbox" type="checkbox" name="is_alive" value="1" ` + checked(person.IsAlive) + ` />
    Còn sống
  </label><br/><br/>

  <div id="death_fields">
    <label>Ngày mất dương lịch (death_date_solar)</label><br/>
    <input type="date" name="death_date_solar" value="` + deathDateSolar + `" /><br/><br/>

    <label>Ngày mất âm lịch (death_date_lunar)</label><br/>
    <input type="date" name="death_date_lunar" value="` + deathDateLunar + `" /><br/><br/>

    <label>Nơi an táng (grave_location)</label><br/>
    <input name="grave_location" value="` + graveLocation + `" /><br/><br/>
  </div>

  <label>Số điện thoại</label><br/>
  <input name="phone" value="` + phone + `" /><br/><br/>

  <label>Nghề nghiệp</label><br/>
  <input name="occupation" value="` + occupation + `" /><br/><br/>

  <label>Ảnh đại diện (URL)</label><br/>
  <input name="avatar_url" value="` + avatarURL + `" /><br/><br/>

  <label>Địa chỉ</label><br/>
  <input name="address" value="` + address + `" /><br/><br/>

  <label>Ghi chú</label><br/>
  <input name="note" value="` + note + `" /><br/><br/>

  <button type="submit" class="btn btn-primary">Lưu</button>
  <a href="/admin/persons" style="
    display: inline-block;
    padding: 10px 20px;
    background: #64748b;
    color: white;
    text-decoration: none;
    border-radius: 6px;
    margin-left: 10px;
    font-weight: 600;
  ">← Quay về Danh sách</a>
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
