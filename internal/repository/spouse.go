package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"familytree/internal/domain"
)

// SpouseRepo is the data access for couples. Each pair is stored once with
// person_a_id < person_b_id, so it is symmetric.
type SpouseRepo struct{ db *pgxpool.Pool }

func NewSpouseRepo(db *pgxpool.Pool) *SpouseRepo { return &SpouseRepo{db: db} }

// ListByPerson returns every partner of the given person.
func (r *SpouseRepo) ListByPerson(ctx context.Context, personID int64) ([]domain.Spouse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.marriage_year, s.note,
		       p.id, p.full_name, p.gender, p.birth_year, p.is_alive
		FROM spouse s
		JOIN person p
		  ON p.id = CASE WHEN s.person_a_id = $1 THEN s.person_b_id ELSE s.person_a_id END
		WHERE s.person_a_id = $1 OR s.person_b_id = $1
		ORDER BY s.marriage_year NULLS LAST, s.id`, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Spouse
	for rows.Next() {
		var (
			sp     domain.Spouse
			gender int16
		)
		if err := rows.Scan(&sp.ID, &sp.MarriageYear, &sp.Note,
			&sp.PartnerID, &sp.PartnerName, &gender, &sp.PartnerBirthYear, &sp.PartnerIsAlive); err != nil {
			return nil, err
		}
		sp.PersonID = personID
		sp.PartnerGender = domain.Gender(gender)
		out = append(out, sp)
	}
	return out, rows.Err()
}

func (r *SpouseRepo) Exists(ctx context.Context, a, b int64) (bool, error) {
	lo, hi := order(a, b)
	var ok bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM spouse WHERE person_a_id = $1 AND person_b_id = $2)`,
		lo, hi).Scan(&ok)
	return ok, err
}

// Create stores a couple. The caller must pass the two IDs in any order; they
// are normalized to (min, max) here so the unique constraint is symmetric.
func (r *SpouseRepo) Create(ctx context.Context, in domain.SpouseInput) (int64, error) {
	lo, hi := order(in.PersonID, in.PartnerID)
	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO spouse (person_a_id, person_b_id, marriage_year, note)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		lo, hi, in.MarriageYear, in.Note,
	).Scan(&id)
	return id, mapError(err)
}

func (r *SpouseRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM spouse WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func order(a, b int64) (int64, int64) {
	if a > b {
		return b, a
	}
	return a, b
}
