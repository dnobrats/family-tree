package service

import (
	"context"
	"errors"

	"genealogy-be/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
====================================================
TREE: CÂY GIA PHẢ CHUNG
GET /api/tree?root={id}
====================================================
*/
func GetTree(ctx context.Context, db *pgxpool.Pool, rootID int64) (*model.TreeResponse, error) {
	rows, err := db.Query(ctx, `
WITH RECURSIVE family_tree AS (
    SELECT
        id,
        full_name,
        gender,
        father_id,
        mother_id,
        is_alive,
        clan_id,
        0 AS depth
    FROM person
    WHERE id = $1

    UNION ALL

    SELECT
        p.id,
        p.full_name,
        p.gender,
        p.father_id,
        p.mother_id,
        p.is_alive,
        p.clan_id,
        ft.depth + 1
    FROM person p
    JOIN family_tree ft
      ON p.father_id = ft.id
      OR p.mother_id = ft.id
)
SELECT
    id,
    full_name,
    gender,
    father_id,
    mother_id,
    is_alive,
    clan_id,
    depth
FROM family_tree
ORDER BY depth, id
`, rootID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := make([]model.PersonNode, 0)
	for rows.Next() {
		var n model.PersonNode
		if err := rows.Scan(
			&n.ID,
			&n.FullName,
			&n.Gender,
			&n.FatherID,
			&n.MotherID,
			&n.IsAlive,
			&n.ClanID,
			&n.Depth,
		); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}

	if len(nodes) == 0 {
		return nil, errors.New("root person not found")
	}

	return &model.TreeResponse{
		RootID: rootID,
		Nodes:  nodes,
	}, nil
}

/*
====================================================
TREE THEO CHI
GET /api/clans/{id}/tree
====================================================
*/
func GetClanTree(ctx context.Context, db *pgxpool.Pool, clanID int64) (*model.ClanTreeResponse, error) {
	var clan model.Clan

	err := db.QueryRow(ctx, `
SELECT
    id,
    name,
    parent_clan_id,
    root_person_id
FROM clan
WHERE id = $1
`, clanID).Scan(
		&clan.ID,
		&clan.Name,
		&clan.ParentClanID,
		&clan.RootPersonID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("clan not found")
		}
		return nil, err
	}

	rows, err := db.Query(ctx, `
WITH RECURSIVE clan_tree AS (
    SELECT
        id,
        full_name,
        gender,
        father_id,
        mother_id,
        is_alive,
        clan_id,
        0 AS depth
    FROM person
    WHERE id = $1

    UNION ALL

    SELECT
        p.id,
        p.full_name,
        p.gender,
        p.father_id,
        p.mother_id,
        p.is_alive,
        p.clan_id,
        ct.depth + 1
    FROM person p
    JOIN clan_tree ct
      ON p.father_id = ct.id
      OR p.mother_id = ct.id
)
SELECT
    id,
    full_name,
    gender,
    father_id,
    mother_id,
    is_alive,
    clan_id,
    depth
FROM clan_tree
ORDER BY depth, id
`, clan.RootPersonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := make([]model.PersonNode, 0)
	for rows.Next() {
		var n model.PersonNode
		if err := rows.Scan(
			&n.ID,
			&n.FullName,
			&n.Gender,
			&n.FatherID,
			&n.MotherID,
			&n.IsAlive,
			&n.ClanID,
			&n.Depth,
		); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}

	return &model.ClanTreeResponse{
		Clan:  clan,
		Nodes: nodes,
	}, nil
}

/*
====================================================
CHI TIẾT 1 NGƯỜI
GET /api/persons/{id}
====================================================
*/
func GetPersonDetail(ctx context.Context, db *pgxpool.Pool, personID int64) (*model.PersonDetail, error) {
	var p model.PersonDetail

	// Clan nullable fields
	var clanID *int64
	var clanName *string
	var parentClanID *int64
	var rootPersonID *int64

	err := db.QueryRow(ctx, `
SELECT
    p.id,
    p.full_name,
    p.gender,
    p.birth_year,
    p.birth_date_solar,
    p.is_alive,
    p.death_date_solar,
    p.death_date_lunar,
    p.father_id,
    p.mother_id,
    p.address,
    p.grave_location,
    p.note,
    c.id,
    c.name,
    c.parent_clan_id,
    c.root_person_id
FROM person p
LEFT JOIN clan c ON c.id = p.clan_id
WHERE p.id = $1
`, personID).Scan(
		&p.ID,
		&p.FullName,
		&p.Gender,
		&p.BirthYear,
		&p.BirthDateSolar,
		&p.IsAlive,
		&p.DeathDateSolar,
		&p.DeathDateLunar,
		&p.FatherID,
		&p.MotherID,
		&p.Address,
		&p.GraveLocation,
		&p.Note,
		&clanID,
		&clanName,
		&parentClanID,
		&rootPersonID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("person not found")
		}
		return nil, err
	}

	if clanID != nil {
		p.Clan = &model.Clan{
			ID:           *clanID,
			Name:         *clanName,
			ParentClanID: parentClanID,
			RootPersonID: *rootPersonID,
		}
	}

	return &p, nil
}
