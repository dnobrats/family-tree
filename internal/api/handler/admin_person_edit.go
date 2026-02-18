package handler

import (
	"net/http"
	"strconv"

	"genealogy-be/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminEditPerson(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
			return
		}

		// Lấy dữ liệu người
		person, err := service.GetPersonByID(r.Context(), db, id)
		if err != nil {
			http.Error(w, "Không tìm thấy người", http.StatusNotFound)
			return
		}

		// Render form edit (đơn giản, không template engine)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
<!doctype html>
<html>
<head><title>Sửa người</title></head>
<body>
<h2>Sửa thông tin</h2>

<form method="post" action="/admin/persons/` + idStr + `">
  <label>Họ tên</label><br/>
  <input name="full_name" value="` + person.FullName + `" /><br/><br/>

  <label>Giới tính</label><br/>
  <select name="gender">
    <option value="1" ` + selected(person.Gender == 1) + `>Nam</option>
    <option value="2" ` + selected(person.Gender == 2) + `>Nữ</option>
  </select><br/><br/>

  <label>Năm sinh</label><br/>
  <input name="birth_year" value="` + intPtrToStr(person.BirthYear) + `" /><br/><br/>

  <label>Cha (ID hoặc tên)</label><br/>
  <input name="father" /><br/><br/>

  <label>Mẹ (ID hoặc tên)</label><br/>
  <input name="mother" /><br/><br/>

  <label>Chi (ID hoặc trưởng chi)</label><br/>
  <input name="clan" /><br/><br/>

  <label>
    <input type="checkbox" name="is_alive" value="1" ` + checked(person.IsAlive) + ` />
    Còn sống
  </label><br/><br/>

  <button>Lưu</button>
</form>

</body>
</html>
		`))
	}
}
