package service

import (
	"context"
	"errors"
	"fmt"

	"genealogy-be/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	repo *repository.Repository
}

func NewAdminService(repo *repository.Repository) *AdminService {
	return &AdminService{repo: repo}
}

// Login xác thực admin user
func (s *AdminService) Login(ctx context.Context, username, password string) error {
	hash, err := s.repo.Admin.GetPasswordHash(ctx, username)
	if err != nil {
		return errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New("invalid username or password")
	}

	return nil
}

// ValidateSession kiểm tra session có hợp lệ không
func (s *AdminService) ValidateSession(sessionToken string) error {
	// TODO: Implement proper session validation
	// Hiện tại chỉ check token có tồn tại không
	if sessionToken == "" {
		return fmt.Errorf("session token is empty")
	}
	return nil
}
