package task

import (
	"context"
	"fmt"
	"strings"
)

func (s *Service) CreateGroup(ctx context.Context, name string) (*Group, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyGroupName
	}
	group := &Group{
		Name: name,
	}
	err := s.groups.Add(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("failed to add group: %w", err)
	}
	return group, nil
}

func (s *Service) GetGroup(ctx context.Context, id int) (*Group, error) {
	if id <= 0 {
		return nil, fmt.Errorf("incorrect id: %d", id)
	}
	group, err := s.groups.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	return group, nil
}

func (s *Service) ListGroup(ctx context.Context) ([]Group, error) {
	groups, err := s.groups.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all groups: %w", err)
	}
	return groups, nil
}

func (s *Service) UpdateGroup(ctx context.Context, id int, name string) (*Group, error) {
	if id <= 0 {
		return nil, fmt.Errorf("incorrect id: %d", id)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyGroupName
	}

	group := &Group{
		ID:   id,
		Name: name,
	}
	err := s.groups.Update(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}
	return group, nil
}

func (s *Service) DeleteGroup(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("incorrect id: %d", id)
	}
	err := s.groups.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}
