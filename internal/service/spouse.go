package service

import (
	"context"

	"familytree/internal/domain"
)

// ListSpouses returns every partner of a person.
func (s *Service) ListSpouses(ctx context.Context, personID int64) ([]domain.Spouse, error) {
	return s.spouses.ListByPerson(ctx, personID)
}

// AddSpouse links two people as a couple. It works for a married-in partner
// (dâu/rể) as long as that person exists in the system.
func (s *Service) AddSpouse(ctx context.Context, in domain.SpouseInput) error {
	if in.PersonID == in.PartnerID {
		return domain.Invalid("Không thể ghép một người với chính họ")
	}

	person, err := s.persons.GetBrief(ctx, in.PersonID)
	if err != nil {
		return domain.Invalid("Không tìm thấy người")
	}
	partner, err := s.persons.GetBrief(ctx, in.PartnerID)
	if err != nil {
		return domain.Invalid("Không tìm thấy vợ/chồng")
	}

	exists, err := s.spouses.Exists(ctx, in.PersonID, in.PartnerID)
	if err != nil {
		return err
	}
	if exists {
		return domain.Invalid("Hai người này đã được ghi nhận là vợ/chồng")
	}

	if in.MarriageYear != nil {
		if person.BirthYear != nil && *in.MarriageYear < *person.BirthYear {
			return domain.Invalid("Năm kết hôn sớm hơn năm sinh của %s", person.FullName)
		}
		if partner.BirthYear != nil && *in.MarriageYear < *partner.BirthYear {
			return domain.Invalid("Năm kết hôn sớm hơn năm sinh của %s", partner.FullName)
		}
	}

	if _, err := s.spouses.Create(ctx, in); err != nil {
		return err
	}
	return nil
}

// DeleteSpouse removes a couple link.
func (s *Service) DeleteSpouse(ctx context.Context, id int64) error {
	return s.spouses.Delete(ctx, id)
}
