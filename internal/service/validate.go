package service

import (
	"context"
	"fmt"

	"genealogy-be/internal/repository"
)

// ValidateParentAge kiểm tra tuổi cha/mẹ hợp lệ
func ValidateParentAge(
	ctx context.Context,
	repo *repository.Repository,
	childBirthYear *int,
	parentID *int64,
	parentRole string, // "Cha" hoặc "Mẹ"
) error {
	if childBirthYear == nil || parentID == nil {
		return nil
	}

	parent, err := repo.Person.GetByID(ctx, *parentID)
	if err != nil {
		return fmt.Errorf("không tìm thấy %s", parentRole)
	}

	if parent.BirthYear == nil {
		return nil // Không có năm sinh của cha/mẹ thì bỏ qua
	}

	if *parent.BirthYear >= *childBirthYear {
		return fmt.Errorf(
			"%s sinh năm %d không thể có con sinh năm %d",
			parentRole,
			*parent.BirthYear,
			*childBirthYear,
		)
	}

	return nil
}
