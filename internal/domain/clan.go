package domain

import "time"

// Clan is a branch (chi) of the family.
type Clan struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	ParentClanID *int64    `json:"parent_clan_id,omitempty"`
	RootPersonID *int64    `json:"root_person_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ClanInput carries form data for create/update.
type ClanInput struct {
	Name         string
	Description  *string
	ParentClanID *int64
	RootPersonID *int64
}
