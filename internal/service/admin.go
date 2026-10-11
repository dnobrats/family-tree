package service

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"familytree/internal/domain"
)

// Authenticate checks a username/password pair. The same generic message is
// returned for both "unknown user" and "wrong password".
func (s *Service) Authenticate(ctx context.Context, username, password string) error {
	admin, err := s.admins.GetByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return domain.Invalid("Tên đăng nhập hoặc mật khẩu không đúng")
	}
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return domain.Invalid("Tên đăng nhập hoặc mật khẩu không đúng")
	}
	return nil
}

// SetAdminPassword creates the account or replaces its password.
func (s *Service) SetAdminPassword(ctx context.Context, username, password string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return domain.Invalid("Tên đăng nhập không được để trống")
	}
	if len(password) < 6 {
		return domain.Invalid("Mật khẩu phải có ít nhất 6 ký tự")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.admins.Upsert(ctx, username, string(hash))
}

// EnsureBootstrapAdmin creates the first admin from env config, but only when
// no admin exists yet. Safe to call on every startup.
func (s *Service) EnsureBootstrapAdmin(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return nil
	}
	n, err := s.admins.Count(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return s.SetAdminPassword(ctx, username, password)
}
