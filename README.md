# Gia phả dòng họ — Family Tree

Ứng dụng web nội bộ để **tra cứu và quản lý gia phả dòng họ**: xem cây gia phả,
tra cứu từng người, quản lý theo chi, và in ra giấy (A4/A3).

- **Backend:** Go (thư viện chuẩn `net/http` + `html/template`)
- **Database:** PostgreSQL
- **Frontend:** HTML server-rendered + CSS/JS nhúng trong binary (không cần build frontend)
- **Đóng gói:** một binary duy nhất; chạy bằng Docker rất nhẹ (distroless)

---

## Chạy nhanh bằng Docker

```bash
cp .env.example .env
# Sửa .env: đặt POSTGRES_PASSWORD và SESSION_SECRET (openssl rand -hex 32)
docker compose up -d --build
```

Mở http://localhost:8080 — migration **tự chạy** khi server khởi động.

Tạo tài khoản admin đầu tiên (chọn một cách):

```bash
# Cách 1: biến môi trường trong .env (tự tạo nếu chưa có admin), rồi docker compose up lại
ADMIN_USERNAME=admin
ADMIN_PASSWORD=mat-khau

# Cách 2: chạy lệnh trong container
docker compose exec app /app/createadmin -username admin -password "mat-khau"
```

## Chạy local (cần Go 1.24+ và PostgreSQL)

```bash
cp .env.example .env      # rồi điền DATABASE_URL và SESSION_SECRET
make tidy                 # tải dependencies
make migrate              # tạo bảng
make seed                 # (tuỳ chọn) nạp dữ liệu demo ~5 đời
make createadmin USER=admin PASS=mat-khau
make run                  # http://localhost:8080
```

---

## Lệnh thường dùng

| Lệnh | Việc |
|------|------|
| `make run` | Chạy server |
| `make build` | Build binary vào `./bin` |
| `make migrate` | Chạy migration rồi thoát |
| `make seed` | Nạp dữ liệu demo |
| `make createadmin USER=.. PASS=..` | Tạo/đổi mật khẩu admin |
| `make test` | Chạy toàn bộ test |
| `make fmt` / `make vet` / `make lint` | Định dạng / kiểm tra tĩnh |
| `docker compose up -d --build` | Dựng toàn bộ stack |

---

## Cấu trúc (tóm tắt)

```
cmd/            điểm khởi chạy (server, createadmin)
internal/       mã nguồn ứng dụng (domain, service, repository, http, auth, config, database)
migrations/     schema SQL, chạy tự động (chỉ thêm file mới, không sửa file cũ)
web/            templates + css + js (nhúng vào binary)
scripts/        script tiện ích (seed demo)
docs/           tài liệu kiến trúc & nghiệp vụ
```

> Muốn sửa nhanh một thứ gì đó? Xem bảng "Muốn đổi X → sửa file nào" trong
> [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md). Ý nghĩa dữ liệu và nghiệp vụ
> nằm ở [`docs/DOMAIN.md`](docs/DOMAIN.md).

---

## Tính năng

- Cây gia phả dạng sơ đồ, tra cứu tên, giới hạn số đời, chọn theo chi, thu gọn nhánh.
- Trang chi tiết từng người: thông tin, cha mẹ, vợ/chồng, con cái.
- Quản trị (cần đăng nhập):
  - Thêm / sửa / xóa thành viên (cha mẹ **không bắt buộc** → tạo được thuỷ tổ).
  - Ghi nhận quan hệ vợ/chồng (hai chiều, gồm dâu/rể).
  - Quản lý chi (tạo/sửa/xóa, chọn người đứng đầu).
- Giao diện responsive (điện thoại & máy tính), **in được A4/A3** (chọn khổ ở trang cây).

## Bảo mật (mức cơ bản, đủ cho nội bộ)

- Mật khẩu hash bằng **bcrypt**; session cookie ký **HMAC-SHA256** (HttpOnly, SameSite=Lax).
- **CSRF** double-submit cho mọi POST.
- **Rate limit theo IP** (chặt hơn ở `/admin/login`).
- Output tự escape (`html/template`) nên chống XSS; SQL dùng tham số.
- `SESSION_SECRET` là bắt buộc — server **không khởi động** nếu thiếu.

> Secrets (`SESSION_SECRET`, mật khẩu DB) nằm trong `.env` và **đã được gitignore**.
> Chỉ `.env.example` (không chứa giá trị thật) được commit.

## Ghi chú kỹ thuật

- Migration được **nhúng vào binary** và ghi vào bảng `schema_migrations`; thêm schema = thêm file `migrations/000N_*.sql`.
- Cây gia phả được **dựng lồng nhau ở Go** (`internal/service/genealogy.go`); JS chỉ render.
- Ảnh đại diện dùng URL (`avatar_url`), chưa có upload.
