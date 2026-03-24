package repository

import (
	"context"
	"genealogy-be/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonRepository struct {
	db *pgxpool.Pool
}

func NewPersonRepository(db *pgxpool.Pool) *PersonRepository {
	return &PersonRepository{db: db}
}

// GetByID lấy thông tin person theo ID
func (r *PersonRepository) GetByID(ctx context.Context, id int64) (*PersonDetail, error) {
	var p PersonDetail
	err := r.db.QueryRow(ctx, `
		SELECT id, full_name, gender, birth_year, birth_date_solar, birth_date_lunar,
		       father_id, mother_id, clan_id,
		       is_alive, death_date_solar, death_date_lunar,
		       address, phone, occupation, avatar_url, grave_location, note
		FROM person
		WHERE id = $1
	`, id).Scan(
		&p.ID, &p.FullName, &p.Gender, &p.BirthYear, &p.BirthDateSolar, &p.BirthDateLunar,
		&p.FatherID, &p.MotherID, &p.ClanID,
		&p.IsAlive, &p.DeathDateSolar, &p.DeathDateLunar,
		&p.Address, &p.Phone, &p.Occupation, &p.AvatarURL, &p.GraveLocation, &p.Note,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetTreeFromRoot lấy cây gia phả từ root person
func (r *PersonRepository) GetTreeFromRoot(ctx context.Context, rootID int64) ([]model.PersonNode, error) {
	rows, err := r.db.Query(ctx, `
		WITH RECURSIVE family_tree AS (
			SELECT id, full_name, gender, father_id, mother_id, is_alive, clan_id, 0 AS depth
			FROM person
			WHERE id = $1
			UNION
			SELECT p.id, p.full_name, p.gender, p.father_id, p.mother_id, p.is_alive, p.clan_id, ft.depth + 1
			FROM person p
			JOIN family_tree ft ON p.father_id = ft.id OR p.mother_id = ft.id
		)
		SELECT DISTINCT id, full_name, gender, father_id, mother_id, is_alive, clan_id, 
		       MIN(depth) as depth
		FROM family_tree
		GROUP BY id, full_name, gender, father_id, mother_id, is_alive, clan_id
		ORDER BY depth, id
	`, rootID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []model.PersonNode
	for rows.Next() {
		var n model.PersonNode
		if err := rows.Scan(&n.ID, &n.FullName, &n.Gender, &n.FatherID, &n.MotherID, &n.IsAlive, &n.ClanID, &n.Depth); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, rows.Err()
}

// Create tạo person mới
func (r *PersonRepository) Create(ctx context.Context, in PersonInput) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO person (
			full_name, gender, birth_year, birth_date_solar, birth_date_lunar,
			father_id, mother_id, clan_id,
			is_alive, death_date_solar, death_date_lunar,
			address, phone, occupation, avatar_url, grave_location, note
		) VALUES (
			$1,$2,$3,$4::date,$5::date,
			$6,$7,$8,
			$9,$10::date,$11,
			$12,$13,$14,$15,$16,$17
		)
	`,
		in.FullName, in.Gender, in.BirthYear, in.BirthDateSolar, in.BirthDateLunar,
		in.FatherID, in.MotherID, in.ClanID,
		in.IsAlive, in.DeathDateSolar, in.DeathDateLunar,
		in.Address, in.Phone, in.Occupation, in.AvatarURL, in.GraveLocation, in.Note,
	)
	return err
}

// Update cập nhật thông tin person
func (r *PersonRepository) Update(ctx context.Context, id int64, in PersonInput) error {
	_, err := r.db.Exec(ctx, `
		UPDATE person SET
			full_name=$1, gender=$2, birth_year=$3, birth_date_solar=$4::date, birth_date_lunar=$5::date,
			father_id=$6, mother_id=$7, clan_id=$8,
			is_alive=$9, death_date_solar=$10::date, death_date_lunar=$11,
			address=$12, phone=$13, occupation=$14, avatar_url=$15, grave_location=$16, note=$17
		WHERE id=$18
	`,
		in.FullName, in.Gender, in.BirthYear, in.BirthDateSolar, in.BirthDateLunar,
		in.FatherID, in.MotherID, in.ClanID,
		in.IsAlive, in.DeathDateSolar, in.DeathDateLunar,
		in.Address, in.Phone, in.Occupation, in.AvatarURL, in.GraveLocation, in.Note,
		id,
	)
	return err
}

// ResolveID resolve person ID từ string (có thể là ID hoặc tên)
func (r *PersonRepository) ResolveID(ctx context.Context, value string) (*int64, error) {
	// Nếu value rỗng, trả về nil (không có cha/mẹ)
	if value == "" {
		return nil, nil
	}

	var id int64
	err := r.db.QueryRow(ctx, `
		SELECT id FROM person
		WHERE id::text = $1 OR full_name ILIKE $2
		LIMIT 1
	`, value, "%"+value+"%").Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// Delete xóa person theo ID
func (r *PersonRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM person WHERE id = $1`, id)
	return err
}

// GetAll lấy danh sách person với pagination
func (r *PersonRepository) GetAll(ctx context.Context, limit, offset int) ([]PersonListItem, int, error) {
	// Đếm tổng số
	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM person`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Lấy danh sách
	rows, err := r.db.Query(ctx, `
		SELECT id, full_name, gender, birth_year, is_alive, clan_id
		FROM person
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []PersonListItem
	for rows.Next() {
		var item PersonListItem
		if err := rows.Scan(&item.ID, &item.FullName, &item.Gender, &item.BirthYear, &item.IsAlive, &item.ClanID); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, rows.Err()
}

// PersonListItem struct cho danh sách
type PersonListItem struct {
	ID        int64
	FullName  string
	Gender    int
	BirthYear *int
	IsAlive   bool
	ClanID    *int64
}

// PersonInput struct cho create/update
type PersonInput struct {
	FullName       string
	Gender         int
	BirthYear      *int
	BirthDateSolar *string
	BirthDateLunar *string
	FatherID       *int64
	MotherID       *int64
	ClanID         *int64
	IsAlive        bool
	DeathDateSolar *string
	DeathDateLunar *string
	Address        *string
	Phone          *string
	Occupation     *string
	AvatarURL      *string
	GraveLocation  *string
	Note           *string
}

// PersonDetail struct cho query detail
type PersonDetail struct {
	ID             int64
	FullName       string
	Gender         int
	BirthYear      *int
	BirthDateSolar *string
	BirthDateLunar *string
	FatherID       *int64
	MotherID       *int64
	ClanID         *int64
	IsAlive        bool
	DeathDateSolar *string
	DeathDateLunar *string
	Address        *string
	Phone          *string
	Occupation     *string
	AvatarURL      *string
	GraveLocation  *string
	Note           *string
}
