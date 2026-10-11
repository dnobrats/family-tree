package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"familytree/internal/domain"
)

// ClanRepo is the data access for clans (chi).
type ClanRepo struct{ db *pgxpool.Pool }

func NewClanRepo(db *pgxpool.Pool) *ClanRepo { return &ClanRepo{db: db} }

const clanColumns = `id, name, description, parent_clan_id, root_person_id, created_at, updated_at`

func scanClan(s scanner) (*domain.Clan, error) {
	var c domain.Clan
	if err := s.Scan(&c.ID, &c.Name, &c.Description, &c.ParentClanID, &c.RootPersonID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, mapError(err)
	}
	return &c, nil
}

func (r *ClanRepo) GetByID(ctx context.Context, id int64) (*domain.Clan, error) {
	row := r.db.QueryRow(ctx, `SELECT `+clanColumns+` FROM clan WHERE id = $1`, id)
	return scanClan(row)
}

func (r *ClanRepo) List(ctx context.Context) ([]domain.Clan, error) {
	rows, err := r.db.Query(ctx, `SELECT `+clanColumns+` FROM clan ORDER BY name, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Clan
	for rows.Next() {
		c, err := scanClan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

// NameExists reports whether another clan already uses this name.
func (r *ClanRepo) NameExists(ctx context.Context, name string, exceptID int64) (bool, error) {
	var ok bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM clan WHERE name = $1 AND id <> $2)`, name, exceptID).Scan(&ok)
	return ok, err
}

func (r *ClanRepo) Create(ctx context.Context, in domain.ClanInput) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO clan (name, description, parent_clan_id, root_person_id)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		in.Name, in.Description, in.ParentClanID, in.RootPersonID,
	).Scan(&id)
	return id, mapError(err)
}

func (r *ClanRepo) Update(ctx context.Context, id int64, in domain.ClanInput) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE clan SET name=$1, description=$2, parent_clan_id=$3, root_person_id=$4, updated_at=now()
		WHERE id=$5`,
		in.Name, in.Description, in.ParentClanID, in.RootPersonID, id)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ClanRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM clan WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
