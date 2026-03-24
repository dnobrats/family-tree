# Migration Guide - Từ Code Cũ sang Code Mới

## Tóm tắt thay đổi

Code đã được refactor hoàn toàn theo Clean Architecture với 3 layer rõ ràng:
- **Repository Layer**: Xử lý database queries
- **Service Layer**: Business logic
- **Handler Layer**: HTTP handling

## Thay đổi chi tiết

### 1. Handler Functions → Handler Methods

**Trước:**
```go
func TreeHandler(db *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        resp, err := service.GetTree(r.Context(), db, rootID)
        // ...
    }
}
```

**Sau:**
```go
func (h *Handler) TreeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        resp, err := h.service.Genealogy.GetTree(r.Context(), rootID)
        // ...
    }
}
```

### 2. Service Functions → Service Methods

**Trước:**
```go
func GetTree(ctx context.Context, db *pgxpool.Pool, rootID int64) (*model.TreeResponse, error) {
    rows, err := db.Query(ctx, `SELECT ...`, rootID)
    // ... SQL query logic
}
```

**Sau:**
```go
func (s *GenealogyService) GetTree(ctx context.Context, rootID int64) (*model.TreeResponse, error) {
    nodes, err := s.repo.Person.GetTreeFromRoot(ctx, rootID)
    // ... business logic only
}
```

### 3. SQL Queries → Repository Methods

**Trước:** SQL queries nằm rải rác trong service

**Sau:** Tất cả SQL queries tập trung trong repository
```go
func (r *PersonRepository) GetTreeFromRoot(ctx context.Context, rootID int64) ([]model.PersonNode, error) {
    rows, err := r.db.Query(ctx, `WITH RECURSIVE ...`, rootID)
    // ... query logic
}
```

### 4. Router Registration

**Trước:**
```go
r.Get("/api/tree", handler.TreeHandler(db))
r.Get("/api/persons/{id}", handler.PersonHandler(db))
```

**Sau:**
```go
h := handler.New(svc)
r.Get("/api/tree", h.TreeHandler())
r.Get("/api/persons/{id}", h.PersonHandler())
```

### 5. Configuration

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
// Type-safe với validation
```

## Dependency Flow

```
main.go
  ↓
config.LoadFromEnv()
  ↓
db.NewPostgres(config)
  ↓
repository.New(db)
  ↓
service.New(repo)
  ↓
handler.New(service)
  ↓
router.NewRouter(handler)
```

## Nếu bạn muốn thêm feature mới

### Ví dụ: Thêm endpoint tìm kiếm person theo tên

**1. Repository (internal/repository/person.go):**
```go
func (r *PersonRepository) SearchByName(ctx context.Context, name string) ([]PersonDetail, error) {
    rows, err := r.db.Query(ctx, `
        SELECT id, full_name, gender, birth_year, ...
        FROM person
        WHERE full_name ILIKE $1
        ORDER BY full_name
    `, "%"+name+"%")
    // ... scan results
}
```

**2. Service (internal/service/person.go):**
```go
func (s *PersonService) SearchByName(ctx context.Context, name string) ([]PersonDetail, error) {
    if name == "" {
        return nil, fmt.Errorf("search name cannot be empty")
    }
    return s.repo.Person.SearchByName(ctx, name)
}
```

**3. Handler (internal/api/handler/person.go):**
```go
func (h *Handler) SearchPersonHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        name := r.URL.Query().Get("name")
        results, err := h.service.Person.SearchByName(r.Context(), name)
        if err != nil {
            respondError(w, http.StatusBadRequest, err.Error())
            return
        }
        respondJSON(w, http.StatusOK, results)
    }
}
```

**4. Router (internal/api/router.go):**
```go
r.Get("/api/persons/search", h.SearchPersonHandler())
```

## Testing Strategy

### Repository Tests
Mock `*pgxpool.Pool` hoặc dùng testcontainers

### Service Tests
Mock `*repository.Repository`
```go
type mockPersonRepo struct {
    mock.Mock
}

func (m *mockPersonRepo) GetByID(ctx context.Context, id int64) (*PersonDetail, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*PersonDetail), args.Error(1)
}
```

### Handler Tests
Mock `*service.Service`
```go
type mockGenealogyService struct {
    mock.Mock
}

func (m *mockGenealogyService) GetTree(ctx context.Context, rootID int64) (*model.TreeResponse, error) {
    args := m.Called(ctx, rootID)
    return args.Get(0).(*model.TreeResponse), args.Error(1)
}
```

## Breaking Changes

### Không còn export các function cũ:
- ❌ `service.GetTree(ctx, db, rootID)`
- ✅ `genealogyService.GetTree(ctx, rootID)`

### Handler signature thay đổi:
- ❌ `func TreeHandler(db *pgxpool.Pool) http.HandlerFunc`
- ✅ `func (h *Handler) TreeHandler() http.HandlerFunc`

### Import paths không đổi:
- ✅ `genealogy-be/internal/service`
- ✅ `genealogy-be/internal/model`
- ✅ `genealogy-be/internal/api/handler`

## Rollback Plan

Nếu cần rollback về code cũ:
```bash
git checkout <commit-before-refactor>
```

Hoặc giữ lại branch cũ:
```bash
git branch backup-before-refactor
```

## Checklist sau khi deploy

- [ ] Database connection hoạt động
- [ ] Tất cả API endpoints trả về đúng
- [ ] Admin login hoạt động
- [ ] CRUD operations hoạt động
- [ ] Rate limiting hoạt động
- [ ] Logs hiển thị đúng
- [ ] Performance không giảm

## Support

Nếu gặp vấn đề, check:
1. CHECKLIST.md - Danh sách các file đã thay đổi
2. REFACTORING.md - Chi tiết về refactoring
3. README.md - Hướng dẫn sử dụng
