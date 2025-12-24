package task

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repo TaskRepository
}

func NewService(repo TaskRepository) *Service {
	return &Service{repo: repo}
}

var (
	ErrEmptyTaskName = errors.New("task name cannot be empty")
	ErrTaskNotFound  = errors.New("task not find")
	ErrNewTaskStatus = errors.New("cannot jump from New to Done: start working first")
)

func (s *Service) CreateTask(ctx context.Context, name, description string) (*Task, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyTaskName
	}

	task := &Task{
		Name:        name,
		Description: description,
		Created:     time.Now(),
		Status:      StatusNew,
	}
	err := s.repo.Add(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to add task: %w", err)
	}
	return task, nil
}

func (s *Service) GetTask(ctx context.Context, id int) (*Task, error) {
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, err
}

func (s *Service) GetAllTasks(ctx context.Context) ([]Task, error) {
	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tasks: %w", err)
	}
	return tasks, nil
}

func (s *Service) UpdateTask(ctx context.Context, id int, name, description string, status TaskStatus) (*Task, error) {
	// вначале достать задачу по id
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	// проверяем валидность имени задачи
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyTaskName
	}
	if task.Status == StatusNew && status == StatusDone {
		return nil, ErrNewTaskStatus
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
