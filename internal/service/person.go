package service

import (
	"context"
	"errors"

	"familytree/internal/domain"
)

// PersonPage is a page of the admin list.
type PersonPage struct {
	Items      []domain.Person
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}

// CreatePerson validates and stores a new person. Parents are optional, so a
// founding ancestor (thủy tổ) can be created with no father/mother.
func (s *Service) CreatePerson(ctx context.Context, in domain.PersonInput) (int64, error) {
	in, err := normalizePersonInput(in)
	if err != nil {
		return 0, err
	}
	if err := s.validateRelations(ctx, 0, in); err != nil {
		return 0, err
	}
	return s.persons.Create(ctx, in)
}

// UpdatePerson validates and updates an existing person.
func (s *Service) UpdatePerson(ctx context.Context, id int64, in domain.PersonInput) error {
	in, err := normalizePersonInput(in)
	if err != nil {
		return err
	}
	if _, err := s.persons.GetByID(ctx, id); err != nil {
		return err
	}
	if err := s.validateRelations(ctx, id, in); err != nil {
		return err
	}
	return s.persons.Update(ctx, id, in)
}

// DeletePerson refuses to delete someone who still has children in the tree.
func (s *Service) DeletePerson(ctx context.Context, id int64) error {
	n, err := s.persons.CountChildren(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return domain.Invalid("Không thể xóa: người này còn %d người con. Hãy xử lý con trước.", n)
	}
	return s.persons.Delete(ctx, id)
}

// GetPersonDetail returns the full profile with relatives resolved.
func (s *Service) GetPersonDetail(ctx context.Context, id int64) (*domain.PersonDetail, error) {
	person, err := s.persons.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	detail := &domain.PersonDetail{
		Person:   *person,
		Spouses:  []domain.Spouse{},
		Children: []domain.PersonBrief{},
	}

	if person.FatherID != nil {
		if f, err := s.persons.GetBrief(ctx, *person.FatherID); err == nil {
			detail.Father = f
		}
	}
	if person.MotherID != nil {
		if m, err := s.persons.GetBrief(ctx, *person.MotherID); err == nil {
			detail.Mother = m
		}
	}
	if person.ClanID != nil {
		if c, err := s.clans.GetByID(ctx, *person.ClanID); err == nil {
			detail.Clan = c
		}
	}
	if sp, err := s.spouses.ListByPerson(ctx, id); err == nil {
		detail.Spouses = sp
	}
	if ch, err := s.persons.ChildrenOf(ctx, id); err == nil {
		detail.Children = ch
	}
	return detail, nil
}

// ListPersons returns a page of people for the admin list.
func (s *Service) ListPersons(ctx context.Context, page, pageSize int) (*PersonPage, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	items, total, err := s.persons.List(ctx, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	return &PersonPage{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
	}, nil
}

// SearchPersons powers the admin autocomplete (father/mother/spouse pickers).
func (s *Service) SearchPersons(ctx context.Context, q string) ([]domain.PersonBrief, error) {
	return s.persons.Search(ctx, q, 20)
}

// ListAllPersons returns every person as a brief (for dropdowns).
func (s *Service) ListAllPersons(ctx context.Context) ([]domain.PersonBrief, error) {
	return s.persons.ListBrief(ctx)
}

// validateRelations enforces parent rules: correct gender, no self-parent, no
// ancestor cycles, and a minimum age gap.
func (s *Service) validateRelations(ctx context.Context, selfID int64, in domain.PersonInput) error {
	if in.FatherID != nil && in.MotherID != nil && *in.FatherID == *in.MotherID {
		return domain.Invalid("Cha và mẹ không thể là cùng một người")
	}
	if err := s.validateParent(ctx, selfID, in.FatherID, in.BirthYear, domain.GenderMale, "Cha"); err != nil {
		return err
	}
	if err := s.validateParent(ctx, selfID, in.MotherID, in.BirthYear, domain.GenderFemale, "Mẹ"); err != nil {
		return err
	}
	return nil
}

func (s *Service) validateParent(ctx context.Context, selfID int64, parentID *int64, childBirthYear *int, want domain.Gender, role string) error {
	if parentID == nil {
		return nil
	}
	if selfID != 0 && *parentID == selfID {
		return domain.Invalid("%s không thể là chính người này", role)
	}

	parent, err := s.persons.GetByID(ctx, *parentID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Invalid("Không tìm thấy %s", role)
		}
		return err
	}
	if parent.Gender != want {
		return domain.Invalid("%s phải có giới tính %s", role, want.Label())
	}

	// The candidate parent must not be a descendant of this person.
	if selfID != 0 {
		inside, err := s.persons.SubtreeContains(ctx, selfID, *parentID)
		if err != nil {
			return err
		}
		if inside {
			return domain.Invalid("%s không thể là con cháu của người này", role)
		}
	}

	if parent.BirthYear != nil && childBirthYear != nil && *childBirthYear-*parent.BirthYear < MinParentAge {
		return domain.Invalid("%s sinh năm %d, phải lớn hơn con ít nhất %d tuổi", role, *parent.BirthYear, MinParentAge)
	}
	return nil
}
