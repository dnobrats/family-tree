package domain

import "time"

// Gender is stored as a small integer: 1 nam, 2 nữ, 0 khác/không rõ.
type Gender int

const (
	GenderUnknown Gender = 0
	GenderMale    Gender = 1
	GenderFemale  Gender = 2
)

// Label returns the Vietnamese display name.
func (g Gender) Label() string {
	switch g {
	case GenderMale:
		return "Nam"
	case GenderFemale:
		return "Nữ"
	default:
		return "Khác"
	}
}

// Person is a single member of the family tree.
type Person struct {
	ID             int64     `json:"id"`
	FullName       string    `json:"full_name"`
	Gender         Gender    `json:"gender"`
	BirthYear      *int      `json:"birth_year,omitempty"`
	BirthDateSolar *string   `json:"birth_date_solar,omitempty"`
	BirthDateLunar *string   `json:"birth_date_lunar,omitempty"`
	IsAlive        bool      `json:"is_alive"`
	DeathYear      *int      `json:"death_year,omitempty"`
	DeathDateSolar *string   `json:"death_date_solar,omitempty"`
	DeathDateLunar *string   `json:"death_date_lunar,omitempty"`
	FatherID       *int64    `json:"father_id,omitempty"`
	MotherID       *int64    `json:"mother_id,omitempty"`
	ClanID         *int64    `json:"clan_id,omitempty"`
	Address        *string   `json:"address,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	Occupation     *string   `json:"occupation,omitempty"`
	AvatarURL      *string   `json:"avatar_url,omitempty"`
	GraveLocation  *string   `json:"grave_location,omitempty"`
	Note           *string   `json:"note,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PersonInput carries raw form data for create/update, before validation.
type PersonInput struct {
	FullName       string
	Gender         Gender
	BirthYear      *int
	BirthDateSolar *string
	BirthDateLunar *string
	IsAlive        bool
	DeathYear      *int
	DeathDateSolar *string
	DeathDateLunar *string
	FatherID       *int64
	MotherID       *int64
	ClanID         *int64
	Address        *string
	Phone          *string
	Occupation     *string
	AvatarURL      *string
	GraveLocation  *string
	Note           *string
}
