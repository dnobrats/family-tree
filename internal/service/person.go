package service

import (
	"context"
	"fmt"

	"genealogy-be/internal/repository"
)

type PersonService struct {
	repo *repository.Repository
}

func NewPersonService(repo *repository.Repository) *PersonService {
	return &PersonService{repo: repo}
}

// CreatePerson tạo person mới
func (s *PersonService) CreatePerson(ctx context.Context, in repository.PersonInput) error {
	// Validate parent age
	if err := ValidateParentAge(ctx, s.repo, in.BirthYear, in.FatherID, "Cha"); err != nil {
		return err
	}
	if err := ValidateParentAge(ctx, s.repo, in.BirthYear, in.MotherID, "Mẹ"); err != nil {
		return err
	}

	if err := s.repo.Person.Create(ctx, in); err != nil {
		return fmt.Errorf("create person: %w", err)
	}
	return nil
}

// UpdatePerson cập nhật thông tin person
func (s *PersonService) UpdatePerson(ctx context.Context, id int64, in repository.PersonInput) error {
	// Validate parent age
	if err := ValidateParentAge(ctx, s.repo, in.BirthYear, in.FatherID, "Cha"); err != nil {
		return err
	}
	if err := ValidateParentAge(ctx, s.repo, in.BirthYear, in.MotherID, "Mẹ"); err != nil {
		return err
	}

	if err := s.repo.Person.Update(ctx, id, in); err != nil {
		return fmt.Errorf("update person: %w", err)
	}
	return nil
}

// ResolvePersonID resolve person ID từ string
func (s *PersonService) ResolvePersonID(ctx context.Context, value string) (*int64, error) {
	return s.repo.Person.ResolveID(ctx, value)
}

// ResolveClanID resolve clan ID từ string
func (s *PersonService) ResolveClanID(ctx context.Context, value string) (*int64, error) {
	return s.repo.Clan.ResolveID(ctx, value)
}

// GetPersonForEdit lấy thông tin person để edit
func (s *PersonService) GetPersonForEdit(ctx context.Context, id int64) (*PersonForEdit, error) {
	detail, err := s.repo.Person.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get person: %w", err)
	}

	return &PersonForEdit{
		ID:             detail.ID,
		FullName:       detail.FullName,
		Gender:         detail.Gender,
		BirthYear:      detail.BirthYear,
		BirthDateSolar: detail.BirthDateSolar,
		BirthDateLunar: detail.BirthDateLunar,
		FatherID:       detail.FatherID,
		MotherID:       detail.MotherID,
		ClanID:         detail.ClanID,
		IsAlive:        detail.IsAlive,
		DeathDateSolar: detail.DeathDateSolar,
		DeathDateLunar: detail.DeathDateLunar,
		Address:        detail.Address,
		Phone:          detail.Phone,
		Occupation:     detail.Occupation,
		AvatarURL:      detail.AvatarURL,
		GraveLocation:  detail.GraveLocation,
		Note:           detail.Note,
	}, nil
}

// DeletePerson xóa person
func (s *PersonService) DeletePerson(ctx context.Context, id int64) error {
	// TODO: Kiểm tra xem person có con không trước khi xóa
	if err := s.repo.Person.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete person: %w", err)
	}
	return nil
}

// ListPersons lấy danh sách person với pagination
func (s *PersonService) ListPersons(ctx context.Context, page, pageSize int) (*PersonListResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	items, total, err := s.repo.Person.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("list persons: %w", err)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &PersonListResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// PersonListResult kết quả danh sách với pagination
type PersonListResult struct {
	Items      []repository.PersonListItem
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}

// PersonForEdit struct cho form edit
type PersonForEdit struct {
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
