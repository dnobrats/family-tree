// Package repository contains all SQL. Each store maps database rows to domain
// entities. No business rules live here — that is the service layer's job.
package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"familytree/internal/domain"
)

// Repository groups the individual stores.
type Repository struct {
	Person *PersonRepo
	Clan   *ClanRepo
	Spouse *SpouseRepo
	Admin  *AdminRepo
}

// New builds every store on top of one shared connection pool.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Person: NewPersonRepo(pool),
		Clan:   NewClanRepo(pool),
		Spouse: NewSpouseRepo(pool),
		Admin:  NewAdminRepo(pool),
	}
}

// scanner is satisfied by both pgx.Row and pgx.Rows.
type scanner interface{ Scan(dest ...any) error }

// mapError converts a PostgreSQL unique-violation into domain.ErrConflict.
func mapError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrConflict
	}
	return err
}
