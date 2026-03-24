package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

// GetPasswordHash lấy password hash của admin theo username
func (r *AdminRepository) GetPasswordHash(ctx context.Context, username string) (string, error) {
	var hash string
	err := r.db.QueryRow(ctx, `
		SELECT password_hash
		FROM admin_user
		WHERE username = $1
	`, username).Scan(&hash)
	return hash, err
}
