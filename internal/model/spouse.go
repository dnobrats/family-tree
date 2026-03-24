package model

// Spouse thông tin vợ/chồng
type Spouse struct {
	ID              int64   `json:"id"`
	PersonID        int64   `json:"person_id"`
	SpouseID        int64   `json:"spouse_id"`
	SpouseName      string  `json:"spouse_name"`
	SpouseGender    int     `json:"spouse_gender"`
	MarriageYear    *int    `json:"marriage_year,omitempty"`
	Note            *string `json:"note,omitempty"`
	SpouseBirthYear *int    `json:"spouse_birth_year,omitempty"`
	SpouseIsAlive   bool    `json:"spouse_is_alive"`
}
