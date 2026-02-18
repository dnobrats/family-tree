package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Person struct {
	ID        int64
	FullName  string
	Gender    int
	BirthYear *int
	FatherID  *int64
	MotherID  *int64
	ClanID    *int64
	IsAlive   bool
	Address   *string
	Note      *string
}

// GetPersonByID lấy đầy đủ thông tin 1 người để phục vụ edit
func GetPersonByID(ctx context.Context, db *pgxpool.Pool, id int64) (*Person, error) {
	var p Person

	err := db.QueryRow(ctx, `
SELECT
  id,
  full_name,
  gender,
  birth_year,
  father_id,
  mother_id,
  clan_id,
  is_alive,
  address,
  note
FROM person
WHERE id = $1
`, id).Scan(
		&p.ID,
		&p.FullName,
		&p.Gender,
		&p.BirthYear,
		&p.FatherID,
		&p.MotherID,
		&p.ClanID,
		&p.IsAlive,
		&p.Address,
		&p.Note,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
