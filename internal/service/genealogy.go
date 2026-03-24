package service

import (
	"context"
	"errors"
	"fmt"

	"genealogy-be/internal/model"
	"genealogy-be/internal/repository"

	"github.com/jackc/pgx/v5"
)

type GenealogyService struct {
	repo *repository.Repository
}

func NewGenealogyService(repo *repository.Repository) *GenealogyService {
	return &GenealogyService{repo: repo}
}

// GetTree lấy cây gia phả từ root person
func (s *GenealogyService) GetTree(ctx context.Context, rootID int64) (*model.TreeResponse, error) {
	nodes, err := s.repo.Person.GetTreeFromRoot(ctx, rootID)
	if err != nil {
		return nil, fmt.Errorf("get tree from root: %w", err)
	}

	if len(nodes) == 0 {
		return nil, errors.New("root person not found")
	}

	return &model.TreeResponse{
		RootID: rootID,
		Nodes:  nodes,
	}, nil
}

// GetClanTree lấy cây gia phả theo chi
func (s *GenealogyService) GetClanTree(ctx context.Context, clanID int64) (*model.ClanTreeResponse, error) {
	clan, err := s.repo.Clan.GetByID(ctx, clanID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("clan not found")
		}
		return nil, fmt.Errorf("get clan: %w", err)
	}

	nodes, err := s.repo.Person.GetTreeFromRoot(ctx, clan.RootPersonID)
	if err != nil {
		return nil, fmt.Errorf("get clan tree: %w", err)
	}

	return &model.ClanTreeResponse{
		Clan:  *clan,
		Nodes: nodes,
	}, nil
}

// GetPersonDetail lấy chi tiết thông tin 1 người
func (s *GenealogyService) GetPersonDetail(ctx context.Context, personID int64) (*model.PersonDetail, error) {
	detail, err := s.repo.Person.GetByID(ctx, personID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("person not found")
		}
		return nil, fmt.Errorf("get person detail: %w", err)
	}

	// Chuyển đổi từ repository.PersonDetail sang model.PersonDetail
	result := &model.PersonDetail{
		ID:             detail.ID,
		FullName:       detail.FullName,
		Gender:         detail.Gender,
		BirthYear:      detail.BirthYear,
		BirthDateSolar: detail.BirthDateSolar,
		BirthDateLunar: detail.BirthDateLunar,
		IsAlive:        detail.IsAlive,
		DeathDateSolar: detail.DeathDateSolar,
		DeathDateLunar: detail.DeathDateLunar,
		FatherID:       detail.FatherID,
		MotherID:       detail.MotherID,
		Address:        detail.Address,
		Phone:          detail.Phone,
		Occupation:     detail.Occupation,
		AvatarURL:      detail.AvatarURL,
		GraveLocation:  detail.GraveLocation,
		Note:           detail.Note,
	}

	// Nếu có clan_id, lấy thông tin clan
	if detail.ClanID != nil {
		clan, err := s.repo.Clan.GetByID(ctx, *detail.ClanID)
		if err == nil {
			result.Clan = clan
		}
	}

	// Lấy danh sách vợ/chồng
	spouses, err := s.repo.Spouse.GetSpousesByPersonID(ctx, personID)
	if err == nil && len(spouses) > 0 {
		result.Spouses = spouses
	}

	return result, nil
}
