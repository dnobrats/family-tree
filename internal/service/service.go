// Package service holds ALL business rules. Handlers and repositories know
// nothing about them; if you need to change how the family tree behaves, this
// is the only package you should need to read.
//
// The store interfaces are declared here (consumer side) so the service can be
// tested with in-memory fakes and never depends on SQL details.
package service

import (
	"context"

	"familytree/internal/domain"
)

// PersonStore is the persistence contract for people.
type PersonStore interface {
	GetByID(ctx context.Context, id int64) (*domain.Person, error)
	GetBrief(ctx context.Context, id int64) (*domain.PersonBrief, error)
	Exists(ctx context.Context, id int64) (bool, error)
	List(ctx context.Context, limit, offset int) ([]domain.Person, int, error)
	ListBrief(ctx context.Context) ([]domain.PersonBrief, error)
	Search(ctx context.Context, q string, limit int) ([]domain.PersonBrief, error)
	RootCandidate(ctx context.Context) (*domain.Person, error)
	ChildrenOf(ctx context.Context, id int64) ([]domain.PersonBrief, error)
	CountChildren(ctx context.Context, id int64) (int, error)
	SubtreeContains(ctx context.Context, rootID, candidate int64) (bool, error)
	Create(ctx context.Context, in domain.PersonInput) (int64, error)
	Update(ctx context.Context, id int64, in domain.PersonInput) error
	Delete(ctx context.Context, id int64) error
}

// ClanStore is the persistence contract for clans (chi).
type ClanStore interface {
	GetByID(ctx context.Context, id int64) (*domain.Clan, error)
	List(ctx context.Context) ([]domain.Clan, error)
	NameExists(ctx context.Context, name string, exceptID int64) (bool, error)
	Create(ctx context.Context, in domain.ClanInput) (int64, error)
	Update(ctx context.Context, id int64, in domain.ClanInput) error
	Delete(ctx context.Context, id int64) error
}

// SpouseStore is the persistence contract for couples.
type SpouseStore interface {
	ListByPerson(ctx context.Context, personID int64) ([]domain.Spouse, error)
	Exists(ctx context.Context, a, b int64) (bool, error)
	Create(ctx context.Context, in domain.SpouseInput) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// AdminStore is the persistence contract for admin accounts.
type AdminStore interface {
	GetByUsername(ctx context.Context, username string) (*domain.Admin, error)
	Count(ctx context.Context) (int, error)
	Upsert(ctx context.Context, username, passwordHash string) error
}

// Service exposes every use case of the app.
type Service struct {
	persons PersonStore
	clans   ClanStore
	spouses SpouseStore
	admins  AdminStore
}

// New wires the service with its stores.
func New(persons PersonStore, clans ClanStore, spouses SpouseStore, admins AdminStore) *Service {
	return &Service{persons: persons, clans: clans, spouses: spouses, admins: admins}
}
