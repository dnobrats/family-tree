package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository chứa tất cả các repository con
type Repository struct {
	Person *PersonRepository
	Clan   *ClanRepository
	Admin  *AdminRepository
	Spouse *SpouseRepository
}

// New tạo repository mới với database pool
func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		Person: NewPersonRepository(db),
		Clan:   NewClanRepository(db),
		Admin:  NewAdminRepository(db),
		Spouse: NewSpouseRepository(db),
	}
}

// Querier interface cho testing (không cần thiết lúc này, có thể xóa hoặc comment)
// type Querier interface {
// 	Exec(ctx context.Context, sql string, arguments ...interface{}) error
// 	Query(ctx context.Context, sql string, args ...interface{}) error
// 	QueryRow(ctx context.Context, sql string, args ...interface{}) error
// }
