# Kiến trúc

Mục tiêu của tài liệu này: giúp bạn (hoặc một AI) **biết ngay cần đọc/sửa file nào**
mà không phải lần theo toàn bộ codebase.

## Sơ đồ tầng

```
        ┌───────────────────────────────────────────────┐
        │  cmd/server            cmd/createadmin         │  điểm khởi chạy
        └───────────────┬───────────────────────────────┘
                        │  wire dependencies
        ┌───────────────▼───────────────────────────────┐
        │  internal/http  (package httpx)                │  routes + handlers
        │    server.go  router & middleware              │
        │    api.go     JSON API                          │
        │    pages_*.go trang HTML                        │
        │    render.go  templates + view models           │
        │    helpers.go parse form / error / json         │
        └───────────────┬───────────────────────────────┘
                        │  calls
        ┌───────────────▼───────────────────────────────┐
        │  internal/service                              │  TOÀN BỘ nghiệp vụ
        │    person / clan / spouse / genealogy / admin   │
        │    validate.go  quy tắc dữ liệu                 │
        └───────────────┬───────────────────────────────┘
                        │  gọi qua interface (PersonStore,...)
        ┌───────────────▼───────────────────────────────┐
        │  internal/repository                           │  chỉ SQL
        └───────────────┬───────────────────────────────┘
                        │
        ┌───────────────▼───────────────────────────────┐
        │  PostgreSQL  (schema trong migrations/)         │
        └───────────────────────────────────────────────┘
```

Các gói phụ trợ: `internal/domain` (thực thể + lỗi), `internal/auth` (session),
`internal/http/middleware` (security/CSRF/rate-limit/auth), `internal/config`,
`internal/database` (pool + migration runner).

## Luồng một request

1. `middleware.Chain` bọc: **SecurityHeaders → RateLimit(IP) → EnsureCSRFCookie**.
2. Route khớp trong `internal/http/server.go`.
3. Với trang admin: thêm **RequireLogin → RequireCSRF**.
4. Handler (`api.go` / `pages_*.go`) parse input (helpers) và gọi `service`.
5. `service` kiểm tra nghiệp vụ rồi gọi `repository`.
6. Kết quả: JSON (`writeJSON`) hoặc HTML (`Renderer.Render`) — dữ liệu vào template
   đã được `html/template` tự escape.

## Muốn đổi X → sửa file nào

| Muốn thay đổi | Sửa ở |
|---|---|
| Thêm/bớt cột trong bảng | `migrations/000N_*.sql` (thêm file mới) + `internal/domain/*.go` + `internal/repository/*.go` |
| Quy tắc nghiệp vụ (tuổi cha mẹ, chu trình, xóa...) | `internal/service/validate.go` và `internal/service/person.go` |
| Cách dựng cây gia phả | `internal/service/genealogy.go` |
| Thêm/đổi route | `internal/http/server.go` (`Routes`) |
| Thêm endpoint API JSON | `internal/http/api.go` |
| Sửa trang HTML | `web/templates/**` |
| Sửa giao diện/màu sắc | `web/static/css/app.css` |
| Sửa hành vi in A4/A3 | `web/static/css/print.css` |
| Hành vi JS (tree, xóa, toggle) | `web/static/js/tree.js`, `web/static/js/admin.js` |
| Xác thực / mật khẩu | `internal/service/admin.go`, `internal/auth/session.go` |
| Bảo mật header / CSRF / rate limit | `internal/http/middleware/*` |
| Config & biến môi trường | `internal/config/config.go`, `.env.example` |

## Thêm một trường cho "person" (ví dụ: số CMND)

1. `migrations/0002_add_person_idcard.sql`: `ALTER TABLE person ADD COLUMN id_card TEXT;`
2. `internal/domain/person.go`: thêm `IDCard *string` vào `Person` và `PersonInput`.
3. `internal/repository/person.go`: thêm cột vào `personColumns`, `scanPerson`, và câu `INSERT/UPDATE`.
4. `internal/service/validate.go`: `cleanText` cho trường mới (nếu cần).
5. `web/templates/admin/person_form.html`: thêm `<input name="id_card">`.
6. `internal/http/helpers.go`: đọc `formOptional(r, "id_card")` trong `personInputFromForm`
   và `formValuesFromPerson` / `formValuesFromRequest`.
7. Hiển thị ở `web/templates/public/person.html` (dùng `{{deref .Detail.Person.IDCard}}`).

## Thêm một trang HTML mới

1. Tạo `web/templates/<khu vực>/<tên>.html` với `{{define "content"}}...{{end}}`.
2. Đăng ký tên → đường dẫn trong `pageTemplates` (`internal/http/render.go`).
3. Thêm handler render trang trong `internal/http/pages_*.go`.
4. Thêm route trong `Routes`.

## Kiểm thử

- Nghiệp vụ: `internal/service/*_test.go` dùng **store giả trong bộ nhớ** (`fake_test.go`),
  không cần DB. Chạy `make test`.
- Auth & middleware: `internal/auth/session_test.go`, `internal/http/middleware/middleware_test.go`.
- Thêm test mới cho nghiệp vụ: dùng `newTestService()` rồi seed dữ liệu qua `ps.seed(...)`.
