package service

import (
	"context"
	"fmt"

	"genealogy-be/internal/repository"
)

type SpouseService struct {
	repo *repository.Repository
}

func NewSpouseService(repo *repository.Repository) *SpouseService {
	return &SpouseService{repo: repo}
}

// AddSpouse thêm quan hệ vợ/chồng
func (s *SpouseService) AddSpouse(ctx context.Context, in repository.SpouseInput) error {
	if err := s.repo.Spouse.Create(ctx, in); err != nil {
		return fmt.Errorf("add spouse: %w", err)
	}
	return nil
}

// DeleteSpouse xóa quan hệ vợ/chồng
func (s *SpouseService) DeleteSpouse(ctx context.Context, id int64) error {
	if err := s.repo.Spouse.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete spouse: %w", err)
	}
	return nil
}
