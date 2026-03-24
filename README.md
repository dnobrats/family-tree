# Genealogy Backend - Hệ thống quản lý gia phả

Backend API cho hệ thống tra cứu và quản lý gia phả dòng họ, được xây dựng với Go và PostgreSQL.

## Tính năng

- 🌳 Xem cây gia phả từ root person
- 👥 Tra cứu thông tin chi tiết từng người
- 🏛️ Quản lý theo chi (clan)
- 🔐 Admin panel để thêm/sửa thông tin
- ✅ Validation tuổi cha/mẹ hợp lệ

## Kiến trúc

Dự án tuân theo **Clean Architecture** với các layer rõ ràng:

```
┌─────────────────────────────────────────────┐
│           HTTP Handlers                      │
│  (Nhận request, trả response)               │
└──────────────┬──────────────────────────────┘
               │
┌──────────────▼──────────────────────────────┐
│           Service Layer                      │
│  (Business logic, validation)               │
└──────────────┬──────────────────────────────┘
               │
┌──────────────▼──────────────────────────────┐
│         Repository Layer                     │
│  (Data access, SQL queries)                 │
└──────────────┬──────────────────────────────┘
               │
┌──────────────▼──────────────────────────────┐
│           Database (PostgreSQL)              │
└─────────────────────────────────────────────┘
```

## Cấu trúc thư mục

```
.
├── cmd/server/              # Entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── db/                 # Database connection
│   ├── repository/         # Data access layer
│   ├── service/            # Business logic layer
│   ├── api/                # HTTP routing
│   │   └── handler/        # HTTP handlers
│   ├── middleware/         # HTTP middlewares
│   └── model/              # Domain models
├── config/                 # Config files
├── .env.example            # Environment variables template
└── Makefile               # Build commands
```

Chi tiết về refactoring xem tại [REFACTORING.md](./REFACTORING.md)

## Yêu cầu

- Go 1.24+
- PostgreSQL 12+

## Cài đặt

1. Clone repository:
```bash
git clone <repo-url>
cd genealogy-be
```

2. Copy và cấu hình environment variables:
```bash
cp .env.example .env
# Sửa file .env với thông tin database của bạn
```

3. Khởi tạo database:
```bash
psql -U your_user -d your_database -f migrations/init.sql
```

4. Build và chạy:
```bash
make build
make run
```

Hoặc development mode với hot reload:
```bash
make dev
```

## API Endpoints

### Public API (Read-only)

- `GET /api/tree?root={id}` - Lấy cây gia phả từ root person
- `GET /api/clans/{id}/tree` - Lấy cây gia phả theo chi
- `GET /api/persons/{id}` - Lấy chi tiết thông tin 1 người

### Admin API (Protected)

- `GET /admin/login` - Trang đăng nhập
- `POST /admin/login` - Xử lý đăng nhập
- `GET /admin` - Dashboard
- `GET /admin/persons/new` - Form thêm người mới
- `POST /admin/persons/new` - Tạo người mới
- `GET /admin/persons/{id}` - Form sửa thông tin
- `POST /admin/persons/{id}` - Cập nhật thông tin

### Public Pages

- `GET /tree` - Trang visualization cây gia phả
- `GET /docs` - API documentation

## Development

### Tạo admin user

```bash
# Generate password hash
make genpass
# Nhập password, copy hash

# Insert vào database
psql -U your_user -d your_database
INSERT INTO admin_user (username, password_hash) VALUES ('admin', 'hash_từ_genpass');
```

### Reset ID sequences (sau khi xóa người)

Nếu bạn xóa nhiều người và muốn ID tiếp tục từ số gần nhất:

```bash
# Cách 1: Dùng SQL script
psql -U your_user -d your_database -f migrations/reset_all_sequences.sql

# Cách 2: Dùng Makefile (cần set env vars trước)
make reset-sequences
```

Xem chi tiết tại [RESET_SEQUENCES.md](./RESET_SEQUENCES.md)

### Run tests

```bash
make test
```

### Format code

```bash
make fmt
```

### Lint code

```bash
make lint
```

## Docker

Build và chạy với Docker:

```bash
docker-compose up -d
```

## Roadmap

- [ ] Unit tests và integration tests
- [ ] JWT authentication thay vì cookie đơn giản
- [ ] API documentation với Swagger
- [ ] Caching với Redis
- [ ] Structured logging
- [ ] Metrics và monitoring
- [ ] Database migrations tool
- [ ] GraphQL API (optional)

## License

MIT

## Đóng góp

Pull requests are welcome! Đối với thay đổi lớn, vui lòng mở issue trước để thảo luận.
