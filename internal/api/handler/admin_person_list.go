package handler

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
)

// AdminPersonList hiển thị danh sách tất cả thành viên
func (h *Handler) AdminPersonList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Lấy page từ query param
		pageStr := r.URL.Query().Get("page")
		page, _ := strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}

		// Lấy danh sách từ service
		result, err := h.service.Person.ListPersons(r.Context(), page, 10)
		if err != nil {
			http.Error(w, "Lỗi tải danh sách", http.StatusInternalServerError)
			return
		}

		// Render HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
<!doctype html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Danh sách thành viên</title>
  <style>
    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      background: #f8fafc;
      margin: 0;
      padding: 20px;
    }
    .container {
      max-width: 1200px;
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
    .actions {
      margin-bottom: 20px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .btn {
      padding: 10px 20px;
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
    .btn-primary:hover {
      background: #2563eb;
    }
    .btn-danger {
      background: #ef4444;
      color: white;
      font-size: 12px;
      padding: 6px 12px;
    }
    .btn-danger:hover {
      background: #dc2626;
    }
    .btn-edit {
      background: #10b981;
      color: white;
      font-size: 12px;
      padding: 6px 12px;
      margin-right: 5px;
    }
    .btn-edit:hover {
      background: #059669;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      margin-bottom: 20px;
    }
    th {
      background: #f1f5f9;
      padding: 12px;
      text-align: left;
      font-weight: 600;
      color: #475569;
      border-bottom: 2px solid #e2e8f0;
    }
    td {
      padding: 12px;
      border-bottom: 1px solid #e2e8f0;
    }
    tr:hover {
      background: #f8fafc;
    }
    .pagination {
      display: flex;
      gap: 10px;
      justify-content: center;
      align-items: center;
    }
    .page-link {
      padding: 8px 12px;
      border: 1px solid #e2e8f0;
      border-radius: 6px;
      text-decoration: none;
      color: #475569;
      transition: all 0.2s;
    }
    .page-link:hover {
      background: #f1f5f9;
      border-color: #3b82f6;
    }
    .page-link.active {
      background: #3b82f6;
      color: white;
      border-color: #3b82f6;
    }
    .page-link.disabled {
      opacity: 0.5;
      pointer-events: none;
    }
    .badge {
      padding: 4px 8px;
      border-radius: 4px;
      font-size: 12px;
      font-weight: 600;
    }
    .badge-alive {
      background: #dcfce7;
      color: #166534;
    }
    .badge-dead {
      background: #fee2e2;
      color: #991b1b;
    }
    .badge-male {
      background: #dbeafe;
      color: #1e40af;
    }
    .badge-female {
      background: #fce7f3;
      color: #9f1239;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>📋 Danh sách thành viên</h1>
    
    <div class="actions">
      <div>
        <a href="/admin" class="btn btn-primary">← Quay lại Dashboard</a>
        <a href="/admin/persons/new" class="btn btn-primary">➕ Thêm người mới</a>
      </div>
      <div style="color: #64748b;">
        Tổng: <strong>` + fmt.Sprintf("%d", result.Total) + `</strong> người
      </div>
    </div>

    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Họ tên</th>
          <th>Giới tính</th>
          <th>Năm sinh</th>
          <th>Trạng thái</th>
          <th style="text-align: center;">Thao tác</th>
        </tr>
      </thead>
      <tbody>
`))

		// Render từng person
		for _, person := range result.Items {
			genderBadge := `<span class="badge badge-male">Nam</span>`
			if person.Gender == 2 || person.Gender == 0 {
				genderBadge = `<span class="badge badge-female">Nữ</span>`
			}

			statusBadge := `<span class="badge badge-alive">Còn sống</span>`
			if !person.IsAlive {
				statusBadge = `<span class="badge badge-dead">Đã mất</span>`
			}

			birthYear := "Không rõ"
			if person.BirthYear != nil {
				birthYear = fmt.Sprintf("%d", *person.BirthYear)
			}

			safeNameHTML := html.EscapeString(person.FullName)
			safeNameJS := strconv.Quote(person.FullName)

			w.Write([]byte(fmt.Sprintf(`
        <tr>
          <td>%d</td>
          <td><strong>%s</strong></td>
          <td>%s</td>
          <td>%s</td>
          <td>%s</td>
          <td style="text-align: center;">
            <a href="/admin/persons/%d/spouses" class="btn btn-edit" style="background: #8b5cf6;">💑 Vợ/Chồng</a>
            <a href="/admin/persons/%d" class="btn btn-edit">✏️ Sửa</a>
            <button onclick="confirmDelete(%d, %s)" class="btn btn-danger">🗑️ Xóa</button>
          </td>
        </tr>
			`, person.ID, safeNameHTML, genderBadge, birthYear, statusBadge, person.ID, person.ID, person.ID, safeNameJS)))
		}

		w.Write([]byte(`
      </tbody>
    </table>

    <div class="pagination">
`))

		// Pagination links
		if page > 1 {
			w.Write([]byte(fmt.Sprintf(`<a href="?page=%d" class="page-link">← Trước</a>`, page-1)))
		} else {
			w.Write([]byte(`<span class="page-link disabled">← Trước</span>`))
		}

		// Page numbers
		for i := 1; i <= result.TotalPages; i++ {
			if i == page {
				w.Write([]byte(fmt.Sprintf(`<span class="page-link active">%d</span>`, i)))
			} else {
				w.Write([]byte(fmt.Sprintf(`<a href="?page=%d" class="page-link">%d</a>`, i, i)))
			}
		}

		if page < result.TotalPages {
			w.Write([]byte(fmt.Sprintf(`<a href="?page=%d" class="page-link">Sau →</a>`, page+1)))
		} else {
			w.Write([]byte(`<span class="page-link disabled">Sau →</span>`))
		}

		w.Write([]byte(`
    </div>
  </div>

  <script>
    function getCookie(name) {
      const cookie = document.cookie
        .split('; ')
        .find(row => row.startsWith(name + '='));
      return cookie ? decodeURIComponent(cookie.split('=')[1]) : '';
    }

    function confirmDelete(id, name) {
      if (confirm('Bạn có chắc muốn xóa "' + name + '"?\n\nLưu ý: Thao tác này không thể hoàn tác!')) {
        const csrfToken = getCookie('csrf_token');
        fetch('/admin/persons/' + id + '/delete', {
          method: 'POST',
          headers: { 'X-CSRF-Token': csrfToken }
        })
        .then(res => {
          if (res.ok) {
            alert('Đã xóa thành công!');
            window.location.reload();
          } else {
            alert('Lỗi khi xóa. Có thể người này có con cháu.');
          }
        })
        .catch(err => {
          alert('Lỗi: ' + err.message);
        });
      }
    }
  </script>
</body>
</html>
		`))
	}
}
