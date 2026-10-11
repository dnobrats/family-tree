package service

import (
	"context"

	"familytree/internal/domain"
)

// ListClans returns all clans.
func (s *Service) ListClans(ctx context.Context) ([]domain.Clan, error) {
	return s.clans.List(ctx)
}

// GetClan returns one clan.
func (s *Service) GetClan(ctx context.Context, id int64) (*domain.Clan, error) {
	return s.clans.GetByID(ctx, id)
}

// CreateClan validates and stores a new clan.
func (s *Service) CreateClan(ctx context.Context, in domain.ClanInput) (int64, error) {
	in, err := normalizeClanInput(in)
	if err != nil {
		return 0, err
	}
	if err := s.validateClanRefs(ctx, in); err != nil {
		return 0, err
	}
	return s.clans.Create(ctx, in)
}

// UpdateClan validates and updates a clan.
func (s *Service) UpdateClan(ctx context.Context, id int64, in domain.ClanInput) error {
	in, err := normalizeClanInput(in)
	if err != nil {
		return err
	}
	if err := s.validateClanRefs(ctx, in); err != nil {
		return err
	}
	return s.clans.Update(ctx, id, in)
}

// DeleteClan removes a clan. Members are detached (clan_id set to NULL) by the
// database, so this is safe.
func (s *Service) DeleteClan(ctx context.Context, id int64) error {
	return s.clans.Delete(ctx, id)
}

func (s *Service) validateClanRefs(ctx context.Context, in domain.ClanInput) error {
	if in.ParentClanID != nil {
		if _, err := s.clans.GetByID(ctx, *in.ParentClanID); err != nil {
			return domain.Invalid("Không tìm thấy chi cha")
		}
	}
	if in.RootPersonID != nil {
		if _, err := s.persons.GetBrief(ctx, *in.RootPersonID); err != nil {
			return domain.Invalid("Không tìm thấy người đứng đầu chi")
		}
	}
	return nil
}
