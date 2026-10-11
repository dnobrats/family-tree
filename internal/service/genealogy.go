package service

import (
	"context"

	"familytree/internal/domain"
)

// GetTree builds the family tree from a root person. If rootID <= 0 the
// earliest ancestor is used. maxGen caps the number of generations rendered.
func (s *Service) GetTree(ctx context.Context, rootID int64, maxGen int) (*domain.TreeNode, error) {
	if rootID <= 0 {
		root, err := s.persons.RootCandidate(ctx)
		if err != nil {
			return nil, err
		}
		rootID = root.ID
	}
	return s.buildNode(ctx, rootID, 1, clampGen(maxGen), map[int64]bool{})
}

// GetClanTree builds the tree rooted at the clan's head person.
func (s *Service) GetClanTree(ctx context.Context, clanID int64, maxGen int) (*domain.Clan, *domain.TreeNode, error) {
	clan, err := s.clans.GetByID(ctx, clanID)
	if err != nil {
		return nil, nil, err
	}
	if clan.RootPersonID == nil {
		return nil, nil, domain.Invalid("Chi này chưa có người đứng đầu")
	}
	node, err := s.buildNode(ctx, *clan.RootPersonID, 1, clampGen(maxGen), map[int64]bool{})
	if err != nil {
		return nil, nil, err
	}
	return clan, node, nil
}

// buildNode recursively builds one node. visited prevents infinite recursion if
// the data ever contains a cycle.
func (s *Service) buildNode(ctx context.Context, id int64, gen, maxGen int, visited map[int64]bool) (*domain.TreeNode, error) {
	if visited[id] {
		return nil, nil
	}
	visited[id] = true

	person, err := s.persons.GetBrief(ctx, id)
	if err != nil {
		return nil, err
	}

	node := &domain.TreeNode{
		Person:     *person,
		Generation: gen,
		Children:   []*domain.TreeNode{},
	}

	// Attach the first partner to display as a couple.
	if spouses, err := s.spouses.ListByPerson(ctx, id); err == nil && len(spouses) > 0 {
		if partner, err := s.persons.GetBrief(ctx, spouses[0].PartnerID); err == nil {
			node.Spouse = partner
		}
	}

	if gen < maxGen {
		children, err := s.persons.ChildrenOf(ctx, id)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			sub, err := s.buildNode(ctx, child.ID, gen+1, maxGen, visited)
			if err != nil {
				return nil, err
			}
			if sub != nil {
				node.Children = append(node.Children, sub)
			}
		}
	}
	return node, nil
}

func clampGen(n int) int {
	if n < 1 {
		return 5
	}
	if n > 15 {
		return 15
	}
	return n
}
