package service

import (
	"strings"
	"time"

	"familytree/internal/domain"
)

// MinParentAge is the minimum number of years between a parent and their child.
const MinParentAge = 12

const (
	maxNameLen = 200
	maxTextLen = 2000
)

// normalizePersonInput trims and applies field-level rules to a person form.
// It returns a cleaned copy ready to persist.
func normalizePersonInput(in domain.PersonInput) (domain.PersonInput, error) {
	in.FullName = strings.TrimSpace(in.FullName)
	if in.FullName == "" {
		return in, domain.Invalid("Họ tên không được để trống")
	}
	if len([]rune(in.FullName)) > maxNameLen {
		return in, domain.Invalid("Họ tên quá dài (tối đa %d ký tự)", maxNameLen)
	}
	if in.Gender < domain.GenderUnknown || in.Gender > domain.GenderFemale {
		return in, domain.Invalid("Giới tính không hợp lệ")
	}

	var err error
	if in.BirthDateSolar, err = cleanDate(in.BirthDateSolar); err != nil {
		return in, err
	}
	if in.DeathDateSolar, err = cleanDate(in.DeathDateSolar); err != nil {
		return in, err
	}
	if in.BirthDateLunar, err = cleanText(in.BirthDateLunar, 100); err != nil {
		return in, err
	}
	if in.DeathDateLunar, err = cleanText(in.DeathDateLunar, 100); err != nil {
		return in, err
	}

	if in.BirthYear != nil && (*in.BirthYear < 1 || *in.BirthYear > 9999) {
		return in, domain.Invalid("Năm sinh không hợp lệ")
	}
	if in.DeathYear != nil && (*in.DeathYear < 1 || *in.DeathYear > 9999) {
		return in, domain.Invalid("Năm mất không hợp lệ")
	}

	// A living person must not carry death details.
	if in.IsAlive {
		in.DeathYear = nil
		in.DeathDateSolar = nil
		in.DeathDateLunar = nil
		in.GraveLocation = nil
	}

	if in.BirthDateSolar != nil && in.DeathDateSolar != nil && *in.DeathDateSolar < *in.BirthDateSolar {
		return in, domain.Invalid("Ngày mất không thể trước ngày sinh")
	}

	if in.Address, err = cleanText(in.Address, maxTextLen); err != nil {
		return in, err
	}
	if in.Phone, err = cleanText(in.Phone, 50); err != nil {
		return in, err
	}
	if in.Occupation, err = cleanText(in.Occupation, 200); err != nil {
		return in, err
	}
	if in.AvatarURL, err = cleanText(in.AvatarURL, 500); err != nil {
		return in, err
	}
	if in.GraveLocation, err = cleanText(in.GraveLocation, 300); err != nil {
		return in, err
	}
	if in.Note, err = cleanText(in.Note, maxTextLen); err != nil {
		return in, err
	}
	return in, nil
}

// normalizeClanInput trims a clan form.
func normalizeClanInput(in domain.ClanInput) (domain.ClanInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return in, domain.Invalid("Tên chi không được để trống")
	}
	if len([]rune(in.Name)) > maxNameLen {
		return in, domain.Invalid("Tên chi quá dài (tối đa %d ký tự)", maxNameLen)
	}
	var err error
	if in.Description, err = cleanText(in.Description, maxTextLen); err != nil {
		return in, err
	}
	return in, nil
}

// cleanDate validates an optional "YYYY-MM-DD" string.
func cleanDate(v *string) (*string, error) {
	if v == nil {
		return nil, nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil, nil
	}
	if _, err := time.Parse("2006-01-02", s); err != nil {
		return nil, domain.Invalid("Ngày không hợp lệ (định dạng YYYY-MM-DD)")
	}
	return &s, nil
}

// cleanText trims an optional free-text field and enforces a length limit.
func cleanText(v *string, max int) (*string, error) {
	if v == nil {
		return nil, nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil, nil
	}
	if len([]rune(s)) > max {
		return nil, domain.Invalid("Nội dung quá dài (tối đa %d ký tự)", max)
	}
	return &s, nil
}
