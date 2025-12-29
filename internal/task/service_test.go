package task

import (
	"context"
	"errors"
	"testing"
)

type MockRepository struct {
	AddCalled bool
}

func (m *MockRepository) Add(ctx context.Context, task *Task) error {
	m.AddCalled = true
	return nil
}

func (m *MockRepository) GetAll(ctx context.Context) ([]Task, error)         { return nil, nil }
func (m *MockRepository) GetById(ctx context.Context, id int) (*Task, error) { return nil, nil }
func (m *MockRepository) Update(ctx context.Context, task *Task) error       { return nil }
func (m *MockRepository) Delete(ctx context.Context, id int) error           { return nil }

func TestCreateTask_EmptyName(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	_, err := service.CreateTask(context.Background(), "", "Описание")

	if !errors.Is(err, ErrEmptyTaskName) {
		t.Errorf("ожидалась ошибка %v, получена %v", ErrEmptyTaskName, err)
	}
	if mockRepo.AddCalled {
		t.Error("репозиторий не должен был вызваться при пустом имени")
	}
}
