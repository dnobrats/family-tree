package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ValidateParentAge(
	ctx context.Context,
	db *pgxpool.Pool,
	childBirthYear *int,
	parentID *int64,
	parentRole string, // "cha" hoặc "mẹ"
) error {
	if childBirthYear == nil || parentID == nil {
		return nil
	}

	var parentBirthYear *int
	err := db.QueryRow(ctx,
		`SELECT birth_year FROM person WHERE id=$1`,
		*parentID,
	).Scan(&parentBirthYear)
	if err != nil {
		return err
	}

	if parentBirthYear == nil {
		return nil
	}

	if *parentBirthYear >= *childBirthYear {
		return fmt.Errorf(
			"%s sinh năm %d không thể có con sinh năm %d",
			parentRole,
			*parentBirthYear,
			*childBirthYear,
		)
	}

	return nil
}
