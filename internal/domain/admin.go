package domain

import "time"

// Admin is a back-office account. PasswordHash is never serialized.
type Admin struct {
	ID           int64
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
