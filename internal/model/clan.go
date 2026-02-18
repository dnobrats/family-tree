package model

type Clan struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ParentClanID *int64 `json:"parent_clan_id"`
	RootPersonID int64  `json:"root_person_id"`
}

