package model

// Response cho /api/clans/{id}/tree
type ClanTreeResponse struct {
	Clan  Clan         `json:"clan"`
	Nodes []PersonNode `json:"nodes"`
}

