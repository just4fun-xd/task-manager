package task

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repo   TaskRepository
	groups GroupRepository
}

func NewService(repo TaskRepository, groups GroupRepository) *Service {
	return &Service{
		repo:   repo,
		groups: groups,
	}
}

var (
	ErrEmptyTaskName    = errors.New("task name cannot be empty")
	ErrTaskNotFound     = errors.New("task not found")
	ErrNewTaskStatus    = errors.New("cannot jump from New to Done: start working first")
	ErrInProgressDelete = errors.New("cannot delete task with InProgress status")
	ErrDoneEdit         = errors.New("cannot edit done task")
	ErrGroupNotFound    = errors.New("group not found")
	ErrGroupHasTasks    = errors.New("group has tasks")
	ErrNotUniqGroup     = errors.New("group has not unique name")
	ErrEmptyGroupName   = errors.New("group name cannot be empty")
)

func (s *Service) CreateTask(ctx context.Context, name, description string, groupId *int) (*Task, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyTaskName
	}
	if groupId != nil {
		if _, err := s.groups.GetById(ctx, *groupId); err != nil {
			return nil, fmt.Errorf("failed to get group: %w", err)
		}

	}

	task := &Task{
		Name:        name,
		Description: description,
		Created:     time.Now(),
		Status:      StatusNew,
		GroupID:     groupId,
	}
	err := s.repo.Add(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to add task: %w", err)
	}
	return task, nil
}

func (s *Service) GetTask(ctx context.Context, id int) (*Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("incorrect id: %d", id)
	}
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

func (s *Service) GetAllTasks(ctx context.Context) ([]Task, error) {
	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tasks: %w", err)
	}
	return tasks, nil
}

func (s *Service) UpdateTask(ctx context.Context, id int, name, description string, status TaskStatus) (*Task, error) {
	task, err := s.GetTask(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task for update: %w", err)
	}
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyTaskName
	}
	if task.Status == StatusNew && status == StatusDone {
		return nil, ErrNewTaskStatus
	}
	if task.Status == StatusDone {
		return nil, ErrDoneEdit
	}

	task.Name = name
	task.Description = description
	task.Status = status

	err = s.repo.Update(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return task, err
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	task, err := s.GetTask(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get task for delete: %w", err)
	}
	if task.Status == StatusInProgress {
		return ErrInProgressDelete
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
