package domain

// Spouse is one partner as seen from a specific person.
type Spouse struct {
	ID               int64   `json:"id"`
	PersonID         int64   `json:"person_id"`
	PartnerID        int64   `json:"partner_id"`
	PartnerName      string  `json:"partner_name"`
	PartnerGender    Gender  `json:"partner_gender"`
	PartnerBirthYear *int    `json:"partner_birth_year,omitempty"`
	PartnerIsAlive   bool    `json:"partner_is_alive"`
	MarriageYear     *int    `json:"marriage_year,omitempty"`
	Note             *string `json:"note,omitempty"`
}

// SpouseInput links two people as a couple. The service stores each pair once.
type SpouseInput struct {
	PersonID     int64
	PartnerID    int64
	MarriageYear *int
	Note         *string
}
