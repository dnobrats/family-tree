package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonInput struct {
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

func CreatePerson(ctx context.Context, db *pgxpool.Pool, in PersonInput) error {
	_, err := db.Exec(ctx, `
INSERT INTO person (
  full_name, gender, birth_year,
  father_id, mother_id, clan_id,
  is_alive, address, note
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
`,
		in.FullName,
		in.Gender,
		in.BirthYear,
		in.FatherID,
		in.MotherID,
		in.ClanID,
		in.IsAlive,
		in.Address,
		in.Note,
	)
	return err
}

func UpdatePerson(ctx context.Context, db *pgxpool.Pool, id int64, in PersonInput) error {
	_, err := db.Exec(ctx, `
UPDATE person SET
  full_name=$1,
  gender=$2,
  birth_year=$3,
  father_id=$4,
  mother_id=$5,
  clan_id=$6,
  is_alive=$7,
  address=$8,
  note=$9
WHERE id=$10
`,
		in.FullName,
		in.Gender,
		in.BirthYear,
		in.FatherID,
		in.MotherID,
		in.ClanID,
		in.IsAlive,
		in.Address,
		in.Note,
		id,
	)
	return err
}
