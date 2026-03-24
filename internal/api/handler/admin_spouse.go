package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"genealogy-be/internal/repository"

	"github.com/go-chi/chi/v5"
)

// AdminSpouseManage hiển thị trang quản lý vợ/chồng của một person
func (h *Handler) AdminSpouseManage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		personID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid person ID", http.StatusBadRequest)
			return
		}

		// Lấy thông tin person
		person, err := h.service.Person.GetPersonForEdit(r.Context(), personID)
		if err != nil {
			http.Error(w, "Không tìm thấy người", http.StatusNotFound)
			return
		}

		// Lấy danh sách spouse hiện tại
		spouses, _ := h.service.Genealogy.GetPersonDetail(r.Context(), personID)
		csrfToken := csrfTokenFromRequest(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
<!doctype html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Quản lý Vợ/Chồng</title>
  <style>
    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      background: #f8fafc;
      margin: 0;
      padding: 20px;
    }
    .container {
      max-width: 800px;
      margin: 0 auto;
      background: white;
      padding: 30px;
      border-radius: 12px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }
    h1 {
      color: #1e293b;
      margin-top: 0;
      border-bottom: 3px solid #3b82f6;
      padding-bottom: 15px;
    }
    .spouse-list {
      margin: 20px 0;
    }
    .spouse-item {
      padding: 15px;
      background: #f8fafc;
      border-radius: 8px;
      margin-bottom: 10px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .btn {
      padding: 8px 16px;
      border: none;
      border-radius: 6px;
      cursor: pointer;
      text-decoration: none;
      display: inline-block;
      font-weight: 600;
      transition: all 0.2s;
    }
    .btn-primary {
      background: #3b82f6;
      color: white;
    }
    .btn-danger {
      background: #ef4444;
      color: white;
      font-size: 12px;
      padding: 6px 12px;
    }
    form {
      background: #f1f5f9;
      padding: 20px;
      border-radius: 8px;
      margin-top: 20px;
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
  </style>
</head>
<body>
  <div class="container">
    <h1>💑 Quản lý Vợ/Chồng: ` + person.FullName + `</h1>
    
    <a href="/admin/persons" class="btn btn-primary">← Quay lại danh sách</a>

    <div class="spouse-list">
      <h3>Danh sách hiện tại:</h3>
`))

		if spouses != nil && len(spouses.Spouses) > 0 {
			for _, sp := range spouses.Spouses {
				marriageYear := ""
				if sp.MarriageYear != nil {
					marriageYear = fmt.Sprintf(" - Kết hôn: %d", *sp.MarriageYear)
				}

				w.Write([]byte(fmt.Sprintf(`
      <div class="spouse-item">
        <div>
          <strong>%s</strong>%s
        </div>
        <button onclick="deleteSpouse(%d)" class="btn btn-danger">Xóa</button>
      </div>
				`, sp.SpouseName, marriageYear, sp.ID)))
			}
		} else {
			w.Write([]byte(`<p style="color: #64748b;">Chưa có thông tin vợ/chồng</p>`))
		}

		w.Write([]byte(fmt.Sprintf(`
    </div>

    <h3>Thêm Vợ/Chồng mới:</h3>
    <form method="post" action="/admin/persons/%d/spouses">
      <input type="hidden" name="csrf_token" value="%s"/>
      <label>Tên Vợ/Chồng (hoặc ID nếu đã có trong hệ thống):</label>
      <input name="spouse_name" required placeholder="Nhập tên hoặc ID"/>

      <label>Năm kết hôn (không bắt buộc):</label>
      <input type="number" name="marriage_year" placeholder="VD: 1990"/>

      <label>Ghi chú:</label>
      <input name="note" placeholder="Ghi chú thêm (nếu có)"/>

      <button type="submit" class="btn btn-primary">➕ Thêm</button>
      <a href="/admin" class="btn" style="background: #64748b; color: white; margin-left: 10px;">← Quay về Admin</a>
    </form>
  </div>

  <script>
    function getCookie(name) {
      const cookie = document.cookie
        .split('; ')
        .find(row => row.startsWith(name + '='));
      return cookie ? decodeURIComponent(cookie.split('=')[1]) : '';
    }

    function deleteSpouse(id) {
      if (confirm('Bạn có chắc muốn xóa quan hệ này?')) {
        const csrfToken = getCookie('csrf_token');
        fetch('/admin/spouses/' + id + '/delete', {
          method: 'POST',
          headers: { 'X-CSRF-Token': csrfToken }
        })
        .then(res => {
          if (res.ok) {
            alert('Đã xóa thành công!');
            window.location.reload();
          } else {
            alert('Lỗi khi xóa.');
          }
        });
      }
    }
  </script>
</body>
</html>
		`, personID, csrfToken)))
	}
}

// AdminSpouseAdd thêm spouse mới
func (h *Handler) AdminSpouseAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Lỗi đọc form", http.StatusBadRequest)
			return
		}

		idStr := chi.URLParam(r, "id")
		personID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid person ID", http.StatusBadRequest)
			return
		}

		spouseNameOrID := r.Form.Get("spouse_name")

		// Thử resolve spouse ID (có thể là ID hoặc tên)
		spouseID, err := h.service.Person.ResolvePersonID(r.Context(), spouseNameOrID)
		if err != nil || spouseID == nil {
			http.Error(w, "Không tìm thấy người với tên/ID: "+spouseNameOrID, http.StatusBadRequest)
			return
		}

		marriageYearStr := r.Form.Get("marriage_year")
		marriageYear, err := parseIntPtr(marriageYearStr)
		if err != nil {
			http.Error(w, "Năm kết hôn không hợp lệ", http.StatusBadRequest)
			return
		}

		note := r.Form.Get("note")
		var notePtr *string
		if note != "" {
			notePtr = &note
		}

		// Tạo spouse relationship
		input := repository.SpouseInput{
			PersonID:     personID,
			SpouseID:     *spouseID,
			MarriageYear: marriageYear,
			Note:         notePtr,
		}

		if err := h.service.Spouse.AddSpouse(r.Context(), input); err != nil {
			http.Error(w, "Lỗi thêm vợ/chồng: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/admin/persons/%d/spouses", personID), http.StatusSeeOther)
	}
}

// AdminSpouseDelete xóa spouse
func (h *Handler) AdminSpouseDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid spouse ID", http.StatusBadRequest)
			return
		}

		if err := h.service.Spouse.DeleteSpouse(r.Context(), id); err != nil {
			http.Error(w, "Lỗi xóa", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
