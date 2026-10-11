package service

import (
	"context"
	"sort"
	"strings"

	"familytree/internal/domain"
)

// ---------------------------------------------------------------------------
// In-memory fakes implementing the store interfaces, so service tests need no
// database.
// ---------------------------------------------------------------------------

type fakePersonStore struct {
	people map[int64]domain.Person
	nextID int64
}

func newFakePersonStore() *fakePersonStore {
	return &fakePersonStore{people: map[int64]domain.Person{}, nextID: 1}
}

func (f *fakePersonStore) seed(p domain.Person) domain.Person {
	if p.ID == 0 {
		p.ID = f.nextID
	}
	f.people[p.ID] = p
	if p.ID >= f.nextID {
		f.nextID = p.ID + 1
	}
	return p
}

func brief(p domain.Person) domain.PersonBrief {
	return domain.PersonBrief{
		ID: p.ID, FullName: p.FullName, Gender: p.Gender,
		BirthYear: p.BirthYear, IsAlive: p.IsAlive, ClanID: p.ClanID,
	}
}

func (f *fakePersonStore) GetByID(_ context.Context, id int64) (*domain.Person, error) {
	p, ok := f.people[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := p
	return &cp, nil
}

func (f *fakePersonStore) GetBrief(_ context.Context, id int64) (*domain.PersonBrief, error) {
	p, ok := f.people[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	b := brief(p)
	return &b, nil
}

func (f *fakePersonStore) Exists(_ context.Context, id int64) (bool, error) {
	_, ok := f.people[id]
	return ok, nil
}

func (f *fakePersonStore) List(_ context.Context, limit, offset int) ([]domain.Person, int, error) {
	all := make([]domain.Person, 0, len(f.people))
	for _, p := range f.people {
		all = append(all, p)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].FullName < all[j].FullName })
	total := len(all)
	lo := min(offset, total)
	hi := min(lo+limit, total)
	return all[lo:hi], total, nil
}

func (f *fakePersonStore) ListBrief(_ context.Context) ([]domain.PersonBrief, error) {
	out := []domain.PersonBrief{}
	for _, p := range f.people {
		out = append(out, brief(p))
	}
	return out, nil
}

func (f *fakePersonStore) Search(_ context.Context, q string, _ int) ([]domain.PersonBrief, error) {
	out := []domain.PersonBrief{}
	for _, p := range f.people {
		if strings.Contains(strings.ToLower(p.FullName), strings.ToLower(q)) {
			out = append(out, brief(p))
		}
	}
	return out, nil
}

func (f *fakePersonStore) RootCandidate(_ context.Context) (*domain.Person, error) {
	var best *domain.Person
	for _, p := range f.people {
		if p.FatherID == nil && p.MotherID == nil {
			cp := p
			if best == nil || p.ID < best.ID {
				best = &cp
			}
		}
	}
	if best == nil {
		return nil, domain.ErrNotFound
	}
	return best, nil
}

func (f *fakePersonStore) ChildrenOf(_ context.Context, id int64) ([]domain.PersonBrief, error) {
	out := []domain.PersonBrief{}
	for _, p := range f.people {
		if isChildOf(p, id) {
			out = append(out, brief(p))
		}
	}
	return out, nil
}

func (f *fakePersonStore) CountChildren(_ context.Context, id int64) (int, error) {
	n := 0
	for _, p := range f.people {
		if isChildOf(p, id) {
			n++
		}
	}
	return n, nil
}

func (f *fakePersonStore) SubtreeContains(ctx context.Context, rootID, candidate int64) (bool, error) {
	if rootID == candidate {
		return true, nil
	}
	seen := map[int64]bool{}
	queue := []int64{rootID}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if seen[cur] {
			continue
		}
		seen[cur] = true
		kids, _ := f.ChildrenOf(ctx, cur)
		for _, k := range kids {
			if k.ID == candidate {
				return true, nil
			}
			queue = append(queue, k.ID)
		}
	}
	return false, nil
}

func (f *fakePersonStore) Create(_ context.Context, in domain.PersonInput) (int64, error) {
	p := domain.Person{
		FullName: in.FullName, Gender: in.Gender, BirthYear: in.BirthYear,
		FatherID: in.FatherID, MotherID: in.MotherID, IsAlive: in.IsAlive, ClanID: in.ClanID,
	}
	return f.seed(p).ID, nil
}

func (f *fakePersonStore) Update(_ context.Context, id int64, in domain.PersonInput) error {
	p, ok := f.people[id]
	if !ok {
		return domain.ErrNotFound
	}
	p.FullName, p.Gender, p.BirthYear = in.FullName, in.Gender, in.BirthYear
	p.FatherID, p.MotherID, p.IsAlive, p.ClanID = in.FatherID, in.MotherID, in.IsAlive, in.ClanID
	f.people[id] = p
	return nil
}

func (f *fakePersonStore) Delete(_ context.Context, id int64) error {
	if _, ok := f.people[id]; !ok {
		return domain.ErrNotFound
	}
	delete(f.people, id)
	return nil
}

func isChildOf(p domain.Person, id int64) bool {
	return (p.FatherID != nil && *p.FatherID == id) || (p.MotherID != nil && *p.MotherID == id)
}

type fakeClanStore struct {
	clans  map[int64]domain.Clan
	nextID int64
}

func newFakeClanStore() *fakeClanStore {
	return &fakeClanStore{clans: map[int64]domain.Clan{}, nextID: 1}
}

func (f *fakeClanStore) seed(c domain.Clan) domain.Clan {
	if c.ID == 0 {
		c.ID = f.nextID
	}
	f.clans[c.ID] = c
	if c.ID >= f.nextID {
		f.nextID = c.ID + 1
	}
	return c
}

func (f *fakeClanStore) GetByID(_ context.Context, id int64) (*domain.Clan, error) {
	c, ok := f.clans[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := c
	return &cp, nil
}

func (f *fakeClanStore) List(_ context.Context) ([]domain.Clan, error) {
	out := []domain.Clan{}
	for _, c := range f.clans {
		out = append(out, c)
	}
	return out, nil
}

func (f *fakeClanStore) NameExists(_ context.Context, name string, exceptID int64) (bool, error) {
	for _, c := range f.clans {
		if c.Name == name && c.ID != exceptID {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeClanStore) Create(_ context.Context, in domain.ClanInput) (int64, error) {
	c := domain.Clan{Name: in.Name, Description: in.Description, ParentClanID: in.ParentClanID, RootPersonID: in.RootPersonID}
	return f.seed(c).ID, nil
}

func (f *fakeClanStore) Update(_ context.Context, id int64, in domain.ClanInput) error {
	c, ok := f.clans[id]
	if !ok {
		return domain.ErrNotFound
	}
	c.Name, c.Description, c.ParentClanID, c.RootPersonID = in.Name, in.Description, in.ParentClanID, in.RootPersonID
	f.clans[id] = c
	return nil
}

func (f *fakeClanStore) Delete(_ context.Context, id int64) error {
	if _, ok := f.clans[id]; !ok {
		return domain.ErrNotFound
	}
	delete(f.clans, id)
	return nil
}

type fakeSpouseStore struct {
	pairs  []domain.Spouse
	nextID int64
}

func newFakeSpouseStore() *fakeSpouseStore {
	return &fakeSpouseStore{nextID: 1}
}

func (f *fakeSpouseStore) ListByPerson(_ context.Context, personID int64) ([]domain.Spouse, error) {
	out := []domain.Spouse{}
	for _, s := range f.pairs {
		if s.PersonID == personID || s.PartnerID == personID {
			out = append(out, perspective(s, personID))
		}
	}
	return out, nil
}

func perspective(s domain.Spouse, viewer int64) domain.Spouse {
	if s.PersonID == viewer {
		return s
	}
	return domain.Spouse{
		ID: s.ID, PersonID: viewer, PartnerID: s.PersonID,
		PartnerName: s.PartnerName, PartnerGender: s.PartnerGender,
		MarriageYear: s.MarriageYear, Note: s.Note,
	}
}

func (f *fakeSpouseStore) Exists(_ context.Context, a, b int64) (bool, error) {
	lo, hi := orderIDs(a, b)
	for _, s := range f.pairs {
		slo, shi := orderIDs(s.PersonID, s.PartnerID)
		if slo == lo && shi == hi {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeSpouseStore) Create(_ context.Context, in domain.SpouseInput) (int64, error) {
	lo, hi := orderIDs(in.PersonID, in.PartnerID)
	s := domain.Spouse{ID: f.nextID, PersonID: lo, PartnerID: hi, MarriageYear: in.MarriageYear, Note: in.Note}
	f.nextID++
	f.pairs = append(f.pairs, s)
	return s.ID, nil
}

func (f *fakeSpouseStore) Delete(_ context.Context, id int64) error {
	for i, s := range f.pairs {
		if s.ID == id {
			f.pairs = append(f.pairs[:i], f.pairs[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func orderIDs(a, b int64) (int64, int64) {
	if a > b {
		return b, a
	}
	return a, b
}

type fakeAdminStore struct {
	byName map[string]domain.Admin
}

func newFakeAdminStore() *fakeAdminStore {
	return &fakeAdminStore{byName: map[string]domain.Admin{}}
}

func (f *fakeAdminStore) GetByUsername(_ context.Context, username string) (*domain.Admin, error) {
	a, ok := f.byName[username]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := a
	return &cp, nil
}

func (f *fakeAdminStore) Count(_ context.Context) (int, error) {
	return len(f.byName), nil
}

func (f *fakeAdminStore) Upsert(_ context.Context, username, passwordHash string) error {
	f.byName[username] = domain.Admin{Username: username, PasswordHash: passwordHash}
	return nil
}

// newTestService wires a Service with fresh fakes.
func newTestService() (*Service, *fakePersonStore, *fakeClanStore, *fakeSpouseStore) {
	ps, cs, ss, as := newFakePersonStore(), newFakeClanStore(), newFakeSpouseStore(), newFakeAdminStore()
	return New(ps, cs, ss, as), ps, cs, ss
}
