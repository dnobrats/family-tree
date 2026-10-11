package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"familytree/internal/domain"
)

// AdminRepo is the data access for back-office accounts.
type AdminRepo struct{ db *pgxpool.Pool }

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo { return &AdminRepo{db: db} }

func (r *AdminRepo) GetByUsername(ctx context.Context, username string) (*domain.Admin, error) {
	var a domain.Admin
	err := r.db.QueryRow(ctx, `
		SELECT id, username, password_hash, created_at
		FROM admin_user WHERE username = $1`, username,
	).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.CreatedAt)
	if err != nil {
		return nil, mapError(err)
	}
	return &a, nil
}

// Count returns how many admin accounts exist (used for first-run bootstrap).
func (r *AdminRepo) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `SELECT count(*) FROM admin_user`).Scan(&n)
	return n, err
}

// Upsert creates the account or replaces its password hash.
func (r *AdminRepo) Upsert(ctx context.Context, username, passwordHash string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO admin_user (username, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (username) DO UPDATE
		SET password_hash = EXCLUDED.password_hash, updated_at = now()`,
		username, passwordHash)
	return mapError(err)
}
