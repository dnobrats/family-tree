package repository

import (
	"context"
	"genealogy-be/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SpouseRepository struct {
	db *pgxpool.Pool
}

func NewSpouseRepository(db *pgxpool.Pool) *SpouseRepository {
	return &SpouseRepository{db: db}
}

// GetSpousesByPersonID lấy danh sách vợ/chồng của một người
func (r *SpouseRepository) GetSpousesByPersonID(ctx context.Context, personID int64) ([]model.Spouse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			s.id, s.person_id, s.spouse_id, s.marriage_year, s.note,
			p.full_name, p.gender, p.birth_year, p.is_alive
		FROM spouse s
		JOIN person p ON p.id = s.spouse_id
		WHERE s.person_id = $1
		ORDER BY s.marriage_year ASC NULLS LAST
	`, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spouses []model.Spouse
	for rows.Next() {
		var sp model.Spouse
		if err := rows.Scan(
			&sp.ID, &sp.PersonID, &sp.SpouseID, &sp.MarriageYear, &sp.Note,
			&sp.SpouseName, &sp.SpouseGender, &sp.SpouseBirthYear, &sp.SpouseIsAlive,
		); err != nil {
			return nil, err
		}
		spouses = append(spouses, sp)
	}
	return spouses, rows.Err()
}

// Create tạo quan hệ vợ/chồng mới
func (r *SpouseRepository) Create(ctx context.Context, in SpouseInput) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO spouse (person_id, spouse_id, marriage_year, note)
		VALUES ($1, $2, $3, $4)
	`, in.PersonID, in.SpouseID, in.MarriageYear, in.Note)
	return err
}

// Update cập nhật thông tin quan hệ vợ/chồng
func (r *SpouseRepository) Update(ctx context.Context, id int64, in SpouseInput) error {
	_, err := r.db.Exec(ctx, `
		UPDATE spouse SET
			marriage_year = $1,
			note = $2
		WHERE id = $3
	`, in.MarriageYear, in.Note, id)
	return err
}

// Delete xóa quan hệ vợ/chồng
func (r *SpouseRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM spouse WHERE id = $1`, id)
	return err
}

// SpouseInput struct cho create/update
type SpouseInput struct {
	PersonID     int64
	SpouseID     int64
	MarriageYear *int
	Note         *string
}
