package domain

// PersonBrief is a lightweight person used in lists and in the tree.
type PersonBrief struct {
	ID        int64  `json:"id"`
	FullName  string `json:"full_name"`
	Gender    Gender `json:"gender"`
	BirthYear *int   `json:"birth_year,omitempty"`
	IsAlive   bool   `json:"is_alive"`
	ClanID    *int64 `json:"clan_id,omitempty"`
}

// TreeNode is a person plus partner (if any) and their children, nested.
// The whole tree is built in Go so the frontend only has to render it.
type TreeNode struct {
	Person     PersonBrief  `json:"person"`
	Spouse     *PersonBrief `json:"spouse,omitempty"`
	Children   []*TreeNode  `json:"children"`
	Generation int          `json:"generation"`
}

// TreeResponse is returned by GET /api/tree.
type TreeResponse struct {
	Root *TreeNode `json:"root"`
}

// ClanTreeResponse is returned by GET /api/clans/{id}/tree.
type ClanTreeResponse struct {
	Clan *Clan     `json:"clan"`
	Root *TreeNode `json:"root"`
}

// PersonDetail is a full profile with relatives resolved.
type PersonDetail struct {
	Person   Person        `json:"person"`
	Clan     *Clan         `json:"clan,omitempty"`
	Father   *PersonBrief  `json:"father,omitempty"`
	Mother   *PersonBrief  `json:"mother,omitempty"`
	Spouses  []Spouse      `json:"spouses"`
	Children []PersonBrief `json:"children"`
}
