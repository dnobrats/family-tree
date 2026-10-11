package service

import (
	"context"
	"errors"
	"testing"

	"familytree/internal/domain"
)

func ptrInt(v int) *int     { return &v }
func ptrI64(v int64) *int64 { return &v }
func strptr(v string) *string { return &v }

func TestNormalizePersonInput_RequiresName(t *testing.T) {
	_, err := normalizePersonInput(domain.PersonInput{FullName: "   ", Gender: domain.GenderMale})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestNormalizePersonInput_ClearsDeathWhenAlive(t *testing.T) {
	got, err := normalizePersonInput(domain.PersonInput{
		FullName: "  Nguyễn Văn A ", Gender: domain.GenderMale, IsAlive: true,
		DeathYear: ptrInt(2000), GraveLocation: strptr("nghĩa trang"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.FullName != "Nguyễn Văn A" {
		t.Errorf("name not trimmed: %q", got.FullName)
	}
	if got.DeathYear != nil || got.GraveLocation != nil {
		t.Errorf("death fields must be cleared for a living person")
	}
}

func TestCreatePerson_RootAncestorAllowed(t *testing.T) {
	svc, ps, _, _ := newTestService()
	id, err := svc.CreatePerson(context.Background(), domain.PersonInput{
		FullName: "Trần Văn Tổ", Gender: domain.GenderMale,
	})
	if err != nil {
		t.Fatalf("a founding ancestor (no parents) must be creatable: %v", err)
	}
	if _, ok := ps.people[id]; !ok {
		t.Fatalf("person was not stored")
	}
}

func TestCreatePerson_FatherMustBeMale(t *testing.T) {
	svc, ps, _, _ := newTestService()
	mom := ps.seed(domain.Person{FullName: "Mẹ", Gender: domain.GenderFemale, BirthYear: ptrInt(1950)})
	_, err := svc.CreatePerson(context.Background(), domain.PersonInput{
		FullName: "Con", Gender: domain.GenderMale, BirthYear: ptrInt(1980), FatherID: ptrI64(mom.ID),
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestCreatePerson_ParentMinAge(t *testing.T) {
	svc, ps, _, _ := newTestService()
	dad := ps.seed(domain.Person{FullName: "Cha", Gender: domain.GenderMale, BirthYear: ptrInt(1975)})
	_, err := svc.CreatePerson(context.Background(), domain.PersonInput{
		FullName: "Con", Gender: domain.GenderMale, BirthYear: ptrInt(1980), FatherID: ptrI64(dad.ID),
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected min-age validation error, got %v", err)
	}
}

func TestUpdatePerson_BlocksAncestorCycle(t *testing.T) {
	svc, ps, _, _ := newTestService()
	grand := ps.seed(domain.Person{FullName: "Ông", Gender: domain.GenderMale, BirthYear: ptrInt(1930)})
	dad := ps.seed(domain.Person{FullName: "Cha", Gender: domain.GenderMale, BirthYear: ptrInt(1955), FatherID: ptrI64(grand.ID)})

	// Making grand the child of dad would create a cycle (dad is grand's child).
	err := svc.UpdatePerson(context.Background(), grand.ID, domain.PersonInput{
		FullName: "Ông", Gender: domain.GenderMale, BirthYear: ptrInt(1930), FatherID: ptrI64(dad.ID),
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected cycle validation error, got %v", err)
	}
}

func TestDeletePerson_RefusedWhenHasChildren(t *testing.T) {
	svc, ps, _, _ := newTestService()
	dad := ps.seed(domain.Person{FullName: "Cha", Gender: domain.GenderMale})
	ps.seed(domain.Person{FullName: "Con", Gender: domain.GenderMale, FatherID: ptrI64(dad.ID)})

	if err := svc.DeletePerson(context.Background(), dad.ID); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected refusal, got %v", err)
	}
}

func TestAddSpouse_IsBidirectional(t *testing.T) {
	svc, ps, _, _ := newTestService()
	ctx := context.Background()
	a := ps.seed(domain.Person{FullName: "A", Gender: domain.GenderMale, BirthYear: ptrInt(1950)})
	b := ps.seed(domain.Person{FullName: "B", Gender: domain.GenderFemale, BirthYear: ptrInt(1952)})

	if err := svc.AddSpouse(ctx, domain.SpouseInput{PersonID: a.ID, PartnerID: b.ID, MarriageYear: ptrInt(1975)}); err != nil {
		t.Fatalf("add spouse: %v", err)
	}

	fromA, _ := svc.ListSpouses(ctx, a.ID)
	fromB, _ := svc.ListSpouses(ctx, b.ID)
	if len(fromA) != 1 || len(fromB) != 1 {
		t.Fatalf("both sides must see the couple: A=%d B=%d", len(fromA), len(fromB))
	}
	if fromB[0].PartnerID != a.ID {
		t.Errorf("B should see partner A, got %d", fromB[0].PartnerID)
	}
}

func TestAddSpouse_RejectsDuplicateAndSelf(t *testing.T) {
	svc, ps, _, _ := newTestService()
	ctx := context.Background()
	a := ps.seed(domain.Person{FullName: "A", Gender: domain.GenderMale})
	b := ps.seed(domain.Person{FullName: "B", Gender: domain.GenderFemale})

	if err := svc.AddSpouse(ctx, domain.SpouseInput{PersonID: a.ID, PartnerID: b.ID}); err != nil {
		t.Fatalf("first add: %v", err)
	}
	if err := svc.AddSpouse(ctx, domain.SpouseInput{PersonID: b.ID, PartnerID: a.ID}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("duplicate couple should be rejected, got %v", err)
	}
	if err := svc.AddSpouse(ctx, domain.SpouseInput{PersonID: a.ID, PartnerID: a.ID}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("self-spouse should be rejected, got %v", err)
	}
}

func TestGetTree_NestedWithSpouse(t *testing.T) {
	svc, ps, _, _ := newTestService()
	ctx := context.Background()

	root := ps.seed(domain.Person{FullName: "Tổ", Gender: domain.GenderMale, BirthYear: ptrInt(1900)})
	mom := ps.seed(domain.Person{FullName: "Bà", Gender: domain.GenderFemale, BirthYear: ptrInt(1902)})
	child := ps.seed(domain.Person{FullName: "Con", Gender: domain.GenderMale, BirthYear: ptrInt(1930), FatherID: ptrI64(root.ID), MotherID: ptrI64(mom.ID)})
	if err := svc.AddSpouse(ctx, domain.SpouseInput{PersonID: root.ID, PartnerID: mom.ID}); err != nil {
		t.Fatalf("spouse: %v", err)
	}

	tree, err := svc.GetTree(ctx, 0, 5)
	if err != nil {
		t.Fatalf("tree: %v", err)
	}
	if tree.Person.ID != root.ID {
		t.Fatalf("root should be %d, got %d", root.ID, tree.Person.ID)
	}
	if tree.Spouse == nil || tree.Spouse.ID != mom.ID {
		t.Fatalf("spouse not attached to root")
	}
	if len(tree.Children) != 1 || tree.Children[0].Person.ID != child.ID {
		t.Fatalf("children wrong: %+v", tree.Children)
	}
	if tree.Children[0].Generation != 2 {
		t.Errorf("child generation should be 2, got %d", tree.Children[0].Generation)
	}
}
