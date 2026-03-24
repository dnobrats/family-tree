# Refactoring Summary - Genealogy Backend

## Cấu trúc mới (Clean Architecture)

```
genealogy-be/
├── cmd/server/              # Entry point
│   └── main.go             # Khởi tạo app với dependency injection
├── internal/
│   ├── config/             # ✨ MỚI: Configuration management
│   │   └── config.go       # Load & validate config từ env
│   ├── db/                 # Database connection
│   │   └── postgres.go     # Connection pool với proper config
│   ├── repository/         # ✨ MỚI: Data access layer
│   │   ├── repository.go   # Repository container
│   │   ├── person.go       # Person queries
│   │   ├── clan.go         # Clan queries
│   │   └── admin.go        # Admin queries
│   ├── service/            # ✨ REFACTORED: Business logic layer
│   │   ├── service.go      # Service container
│   │   ├── genealogy.go    # Genealogy business logic
│   │   ├── person.go       # Person CRUD logic
│   │   ├── admin.go        # Admin authentication
│   │   └── validate.go     # Validation logic
│   ├── api/
│   │   ├── router.go       # ✨ REFACTORED: Router với DI
│   │   └── handler/        # ✨ REFACTORED: HTTP handlers
│   │       ├── handler.go  # Handler struct với service injection
│   │       ├── tree.go     # Tree endpoints
│   │       ├── clan.go     # Clan endpoints
│   │       ├── person.go   # Person endpoints
│   │       └── admin_*.go  # Admin endpoints
│   ├── middleware/         # HTTP middlewares
│   └── model/              # Domain models
```

## Những gì đã thay đổi

### 1. Config Management (internal/config/)
**Trước:**
```go
cfg := map[string]string{
    "host": os.Getenv("DB_HOST"),
    // ...
}
```

**Sau:**
```go
cfg, err := config.LoadFromEnv()
// Có validation, type-safe, default values
```

### 2. Repository Layer (internal/repository/)
**Mới thêm** - Tách biệt SQL queries khỏi business logic:
- `PersonRepository`: Tất cả queries liên quan đến person
- `ClanRepository`: Tất cả queries liên quan đến clan
- `AdminRepository`: Queries cho authentication

**Lợi ích:**
- Dễ test (mock repository)
- SQL queries tập trung một chỗ
- Tái sử dụng queries

### 3. Service Layer (internal/service/)
**Trước:**
```go
func GetTree(ctx, db, rootID) // Function với DB injection
```

**Sau:**
```go
type GenealogyService struct {
    repo *repository.Repository
}
func (s *GenealogyService) GetTree(ctx, rootID) // Method với repo injection
```

**Lợi ích:**
- Dependency injection rõ ràng
- Dễ test với mock repository
- Business logic tách biệt khỏi data access

### 4. Handler Layer (internal/api/handler/)
**Trước:**
```go
func TreeHandler(db *pgxpool.Pool) http.HandlerFunc {
    // Gọi service.GetTree(ctx, db, rootID)
}
```

**Sau:**
```go
type Handler struct {
    service *service.Service
}
func (h *Handler) TreeHandler() http.HandlerFunc {
    // Gọi h.service.Genealogy.GetTree(ctx, rootID)
}
```

**Lợi ích:**
- Handler không biết về database
- Dễ test với mock service
- Consistent error handling với `respondJSON` và `respondError`

### 5. Router (internal/api/router.go)
**Dependency Injection Flow:**
```go
DB Pool → Repository → Service → Handler → Router
```

```go
repo := repository.New(db)
svc := service.New(repo)
h := handler.New(svc)
// Sử dụng h.TreeHandler(), h.PersonHandler(), etc.
```

### 6. Database Connection (internal/db/postgres.go)
**Cải thiện:**
- Connection pool configuration (MaxConns, MinConns, timeouts)
- Proper error handling với context timeout
- Ping để verify connection

## Migration Guide

### Nếu muốn thêm endpoint mới:

1. **Thêm query vào Repository:**
```go
// internal/repository/person.go
func (r *PersonRepository) GetByName(ctx context.Context, name string) (*Person, error) {
    // SQL query
}
```

2. **Thêm business logic vào Service:**
```go
// internal/service/person.go
func (s *PersonService) SearchByName(ctx context.Context, name string) (*Person, error) {
    return s.repo.Person.GetByName(ctx, name)
}
```

3. **Thêm handler:**
```go
// internal/api/handler/person.go
func (h *Handler) SearchPersonHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        name := r.URL.Query().Get("name")
        person, err := h.service.Person.SearchByName(r.Context(), name)
        // ...
        respondJSON(w, http.StatusOK, person)
    }
}
```

4. **Đăng ký route:**
```go
// internal/api/router.go
r.Get("/api/persons/search", h.SearchPersonHandler())
```

## Environment Variables

```bash
# Server
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SCHEMA=public
```

## Testing Strategy

### Unit Tests
- **Repository**: Mock pgxpool.Pool
- **Service**: Mock repository.Repository
- **Handler**: Mock service.Service

### Integration Tests
- Test với real database (testcontainers)

## Các cải tiến tiếp theo có thể làm

1. **Session Management**: Thay cookie đơn giản bằng JWT hoặc session store
2. **Validation**: Thêm struct validation với tags (go-playground/validator)
3. **Logging**: Structured logging với zerolog hoặc zap
4. **Metrics**: Prometheus metrics cho monitoring
5. **API Documentation**: Swagger/OpenAPI spec
6. **Testing**: Thêm unit tests và integration tests
7. **Migration**: Database migration tool (golang-migrate)
8. **Caching**: Redis cache cho queries thường dùng

## Kết luận

Code giờ đã:
- ✅ Tuân theo Clean Architecture
- ✅ Dễ test với dependency injection
- ✅ Tách biệt concerns (data access, business logic, HTTP handling)
- ✅ Type-safe configuration
- ✅ Consistent error handling
- ✅ Dễ mở rộng và maintain
