package handler

import "net/http"

// AdminHome hiển thị trang admin dashboard
func (h *Handler) AdminHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
<!doctype html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Admin Dashboard</title>
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
      padding: 40px;
      border-radius: 12px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }
    h1 {
      color: #1e293b;
      margin-top: 0;
      border-bottom: 3px solid #3b82f6;
      padding-bottom: 15px;
    }
    .menu {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 20px;
      margin-top: 30px;
    }
    .menu-item {
      display: block;
      padding: 30px 20px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      text-decoration: none;
      border-radius: 12px;
      text-align: center;
      font-size: 18px;
      font-weight: 600;
      transition: transform 0.2s, box-shadow 0.2s;
      box-shadow: 0 4px 6px rgba(0,0,0,0.1);
    }
    .menu-item:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 12px rgba(0,0,0,0.15);
    }
    .menu-item.green {
      background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
    }
    .menu-item.orange {
      background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
    }
    .menu-item.blue {
      background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    }
    .icon {
      font-size: 32px;
      display: block;
      margin-bottom: 10px;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>🎛️ Admin Dashboard</h1>
    <p style="color: #64748b; margin-bottom: 30px;">Chào mừng đến với trang quản trị gia phả</p>
    
    <div class="menu">
      <a href="/admin/persons" class="menu-item">
        <span class="icon">📋</span>
        Danh sách thành viên
      </a>
      <a href="/admin/persons/new" class="menu-item green">
        <span class="icon">➕</span>
        Thêm người mới
      </a>
      <a href="/tree" class="menu-item orange">
        <span class="icon">🌳</span>
        Xem cây gia phả
      </a>
      <a href="/docs" class="menu-item blue">
        <span class="icon">📖</span>
        API Documentation
      </a>
    </div>
  </div>
</body>
</html>
		`))
	}
}
