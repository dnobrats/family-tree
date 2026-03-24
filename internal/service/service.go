package service

import "genealogy-be/internal/repository"

// Service chứa tất cả business logic services
type Service struct {
	Genealogy *GenealogyService
	Person    *PersonService
	Admin     *AdminService
	Spouse    *SpouseService
}

// New tạo service mới với repository
func New(repo *repository.Repository) *Service {
	return &Service{
		Genealogy: NewGenealogyService(repo),
		Person:    NewPersonService(repo),
		Admin:     NewAdminService(repo),
		Spouse:    NewSpouseService(repo),
	}
}
