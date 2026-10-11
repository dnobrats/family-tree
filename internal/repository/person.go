package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"familytree/internal/domain"
)

// PersonRepo is the data access for people.
type PersonRepo struct{ db *pgxpool.Pool }

func NewPersonRepo(db *pgxpool.Pool) *PersonRepo { return &PersonRepo{db: db} }

// personColumns is the canonical SELECT list. to_char keeps dates as plain
// strings so the app never deals with time zones.
const personColumns = `id, full_name, gender, birth_year,
	to_char(birth_date_solar, 'YYYY-MM-DD'), birth_date_lunar,
	is_alive, death_year, to_char(death_date_solar, 'YYYY-MM-DD'), death_date_lunar,
	father_id, mother_id, clan_id,
	address, phone, occupation, avatar_url, grave_location, note,
	created_at, updated_at`

const briefColumns = `id, full_name, gender, birth_year, is_alive, clan_id`

func scanPerson(s scanner) (*domain.Person, error) {
	var (
		p      domain.Person
		gender int16
	)
	err := s.Scan(
		&p.ID, &p.FullName, &gender, &p.BirthYear,
		&p.BirthDateSolar, &p.BirthDateLunar,
		&p.IsAlive, &p.DeathYear, &p.DeathDateSolar, &p.DeathDateLunar,
		&p.FatherID, &p.MotherID, &p.ClanID,
		&p.Address, &p.Phone, &p.Occupation, &p.AvatarURL, &p.GraveLocation, &p.Note,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	p.Gender = domain.Gender(gender)
	return &p, nil
}

func scanBrief(s scanner) (*domain.PersonBrief, error) {
	var (
		b      domain.PersonBrief
		gender int16
	)
	if err := s.Scan(&b.ID, &b.FullName, &gender, &b.BirthYear, &b.IsAlive, &b.ClanID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	b.Gender = domain.Gender(gender)
	return &b, nil
}

func (r *PersonRepo) GetByID(ctx context.Context, id int64) (*domain.Person, error) {
	row := r.db.QueryRow(ctx, `SELECT `+personColumns+` FROM person WHERE id = $1`, id)
	return scanPerson(row)
}

func (r *PersonRepo) GetBrief(ctx context.Context, id int64) (*domain.PersonBrief, error) {
	row := r.db.QueryRow(ctx, `SELECT `+briefColumns+` FROM person WHERE id = $1`, id)
	return scanBrief(row)
}

func (r *PersonRepo) Exists(ctx context.Context, id int64) (bool, error) {
	var ok bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM person WHERE id = $1)`, id).Scan(&ok)
	return ok, err
}

func (r *PersonRepo) List(ctx context.Context, limit, offset int) ([]domain.Person, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM person`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT `+personColumns+` FROM person ORDER BY full_name, id LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	people := make([]domain.Person, 0, limit)
	for rows.Next() {
		p, err := scanPerson(rows)
		if err != nil {
			return nil, 0, err
		}
		people = append(people, *p)
	}
	return people, total, rows.Err()
}

func (r *PersonRepo) ListBrief(ctx context.Context) ([]domain.PersonBrief, error) {
	rows, err := r.db.Query(ctx, `SELECT `+briefColumns+` FROM person ORDER BY full_name, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.PersonBrief
	for rows.Next() {
		b, err := scanBrief(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *b)
	}
	return out, rows.Err()
}

func (r *PersonRepo) Search(ctx context.Context, q string, limit int) ([]domain.PersonBrief, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+briefColumns+` FROM person
		 WHERE full_name ILIKE '%' || $1 || '%'
		 ORDER BY full_name, id LIMIT $2`, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.PersonBrief
	for rows.Next() {
		b, err := scanBrief(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *b)
	}
	return out, rows.Err()
}

// RootCandidate returns the earliest ancestor (no parents, lowest id).
func (r *PersonRepo) RootCandidate(ctx context.Context) (*domain.Person, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+personColumns+` FROM person
		 WHERE father_id IS NULL AND mother_id IS NULL
		 ORDER BY id LIMIT 1`)
	return scanPerson(row)
}

func (r *PersonRepo) ChildrenOf(ctx context.Context, id int64) ([]domain.PersonBrief, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+briefColumns+` FROM person
		 WHERE father_id = $1 OR mother_id = $1
		 ORDER BY birth_year NULLS LAST, id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.PersonBrief
	for rows.Next() {
		b, err := scanBrief(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *b)
	}
	return out, rows.Err()
}

func (r *PersonRepo) CountChildren(ctx context.Context, id int64) (int, error) {
	var n int
	err := r.db.QueryRow(ctx,
		`SELECT count(*) FROM person WHERE father_id = $1 OR mother_id = $1`, id).Scan(&n)
	return n, err
}

// SubtreeContains reports whether candidate is anywhere in the subtree rooted
// at rootID (candidate is a descendant of rootID). Used to block ancestor cycles.
func (r *PersonRepo) SubtreeContains(ctx context.Context, rootID, candidate int64) (bool, error) {
	var found bool
	err := r.db.QueryRow(ctx, `
		WITH RECURSIVE subtree AS (
			SELECT id FROM person WHERE id = $1
			UNION
			SELECT p.id FROM person p
			JOIN subtree s ON p.father_id = s.id OR p.mother_id = s.id
		)
		SELECT EXISTS(SELECT 1 FROM subtree WHERE id = $2)`,
		rootID, candidate,
	).Scan(&found)
	return found, err
}

func (r *PersonRepo) Create(ctx context.Context, in domain.PersonInput) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO person (
			full_name, gender, birth_year, birth_date_solar, birth_date_lunar,
			is_alive, death_year, death_date_solar, death_date_lunar,
			father_id, mother_id, clan_id, address, phone, occupation, avatar_url, grave_location, note
		) VALUES (
			$1,$2,$3,$4::date,$5,
			$6,$7,$8::date,$9,
			$10,$11,$12,$13,$14,$15,$16,$17,$18
		) RETURNING id`,
		in.FullName, int16(in.Gender), in.BirthYear, in.BirthDateSolar, in.BirthDateLunar,
		in.IsAlive, in.DeathYear, in.DeathDateSolar, in.DeathDateLunar,
		in.FatherID, in.MotherID, in.ClanID,
		in.Address, in.Phone, in.Occupation, in.AvatarURL, in.GraveLocation, in.Note,
	).Scan(&id)
	return id, err
}

func (r *PersonRepo) Update(ctx context.Context, id int64, in domain.PersonInput) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE person SET
			full_name=$1, gender=$2, birth_year=$3, birth_date_solar=$4::date, birth_date_lunar=$5,
			is_alive=$6, death_year=$7, death_date_solar=$8::date, death_date_lunar=$9,
			father_id=$10, mother_id=$11, clan_id=$12,
			address=$13, phone=$14, occupation=$15, avatar_url=$16, grave_location=$17, note=$18,
			updated_at = now()
		WHERE id=$19`,
		in.FullName, int16(in.Gender), in.BirthYear, in.BirthDateSolar, in.BirthDateLunar,
		in.IsAlive, in.DeathYear, in.DeathDateSolar, in.DeathDateLunar,
		in.FatherID, in.MotherID, in.ClanID,
		in.Address, in.Phone, in.Occupation, in.AvatarURL, in.GraveLocation, in.Note,
		id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PersonRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM person WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
