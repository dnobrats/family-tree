package repository

import (
	"context"
	"genealogy-be/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClanRepository struct {
	db *pgxpool.Pool
}

func NewClanRepository(db *pgxpool.Pool) *ClanRepository {
	return &ClanRepository{db: db}
}

// GetByID lấy thông tin clan theo ID
func (r *ClanRepository) GetByID(ctx context.Context, id int64) (*model.Clan, error) {
	var c model.Clan
	err := r.db.QueryRow(ctx, `
		SELECT id, name, parent_clan_id, root_person_id
		FROM clan
		WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.ParentClanID, &c.RootPersonID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetAll lấy tất cả clans
func (r *ClanRepository) GetAll(ctx context.Context) ([]model.Clan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, parent_clan_id, root_person_id
		FROM clan
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clans []model.Clan
	for rows.Next() {
		var c model.Clan
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentClanID, &c.RootPersonID); err != nil {
			return nil, err
		}
		clans = append(clans, c)
	}
	return clans, rows.Err()
}

// ResolveID resolve clan ID từ string (có thể là ID hoặc tên)
func (r *ClanRepository) ResolveID(ctx context.Context, value string) (*int64, error) {
	// Nếu value rỗng, trả về nil (không thuộc chi nào)
	if value == "" {
		return nil, nil
	}

	var id int64
	err := r.db.QueryRow(ctx, `
		SELECT id FROM clan
		WHERE id::text = $1 OR name ILIKE $2
		LIMIT 1
	`, value, "%"+value+"%").Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
