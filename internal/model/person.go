package model

type PersonNode struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Gender   int    `json:"gender"`
	FatherID *int64 `json:"father_id"`
	MotherID *int64 `json:"mother_id"`
	IsAlive  bool   `json:"is_alive"`
	ClanID   *int64 `json:"clan_id"`
	Depth    int    `json:"depth"`
}

type TreeResponse struct {
	RootID int64        `json:"root_id"`
	Nodes  []PersonNode `json:"nodes"`
}
