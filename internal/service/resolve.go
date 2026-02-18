package service

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ResolvePersonID(ctx context.Context, db *pgxpool.Pool, v string) (*int64, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}

	// Case 1: numeric ID
	if id, err := strconv.ParseInt(v, 10, 64); err == nil {
		return &id, nil
	}

	// Case 2: search by name (case-insensitive)
	rows, err := db.Query(ctx, `
		SELECT id FROM person
		WHERE full_name ILIKE $1
	`, v)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, errors.New("không tìm thấy người: " + v)
	}
	if len(ids) > 1 {
		return nil, errors.New("trùng tên, vui lòng dùng ID: " + v)
	}

	return &ids[0], nil
}

func ResolveClanID(ctx context.Context, db *pgxpool.Pool, v string) (*int64, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}

	// Case 1: numeric ID
	if id, err := strconv.ParseInt(v, 10, 64); err == nil {
		return &id, nil
	}

	// Case 2: find clan by root person name
	rows, err := db.Query(ctx, `
		SELECT c.id
		FROM clan c
		JOIN person p ON p.id = c.root_person_id
		WHERE p.full_name ILIKE $1
	`, v)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, errors.New("không tìm thấy chi theo trưởng chi: " + v)
	}
	if len(ids) > 1 {
		return nil, errors.New("trùng trưởng chi, vui lòng dùng ID")
	}

	return &ids[0], nil
}
