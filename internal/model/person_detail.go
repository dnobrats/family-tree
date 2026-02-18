package model

type PersonDetail struct {
	ID             int64   `json:"id"`
	FullName       string  `json:"full_name"`
	Gender         int     `json:"gender"`
	BirthYear      *int    `json:"birth_year"`
	BirthDateSolar *string `json:"birth_date_solar"`
	IsAlive        bool    `json:"is_alive"`

	DeathDateSolar *string `json:"death_date_solar"`
	DeathDateLunar *string `json:"death_date_lunar"`

	FatherID *int64 `json:"father_id"`
	MotherID *int64 `json:"mother_id"`

	Clan *Clan `json:"clan"`

	Address       *string `json:"address"`
	GraveLocation *string `json:"grave_location"`
	Note          *string `json:"note"`
}
