package task

import (
	"context"
	"errors"
	"testing"
)

type MockRepository struct {
	AddCalled    bool
	UpdateCalled bool
	TaskToReturn *Task
	UpdatedTask  *Task
}

func (m *MockRepository) Add(ctx context.Context, task *Task) error {
	m.AddCalled = true
	return nil
}

func (m *MockRepository) GetAll(ctx context.Context) ([]Task, error) { return nil, nil }
func (m *MockRepository) GetById(ctx context.Context, id int) (*Task, error) {
	return m.TaskToReturn, nil
}
func (m *MockRepository) Update(ctx context.Context, task *Task) error {
	m.UpdatedTask = task
	m.UpdateCalled = true
	return nil
}
func (m *MockRepository) Delete(ctx context.Context, id int) error { return nil }

func TestCreateTask_EmptyName(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo, nil)
	_, err := service.CreateTask(context.Background(), "", "Описание", nil) // добавить

	if !errors.Is(err, ErrEmptyTaskName) {
		t.Errorf("ожидалась ошибка %v, получена %v", ErrEmptyTaskName, err)
	}
	if mockRepo.AddCalled {
		t.Error("репозиторий не должен был вызваться при пустом имени")
	}
}

func TestUpdateTask_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		existingStatus TaskStatus
		newName        string
		newStatus      TaskStatus
		expectedErr    error
		wantUpdate     bool
	}{
		{
			name:           "Успешное обновление",
			existingStatus: StatusInProgress,
			newName:        "Новое имя",
			newStatus:      StatusDone,
			expectedErr:    nil,
			wantUpdate:     true,
		},
		{
			name:           "Ошибка: пустой заголовок",
			existingStatus: StatusInProgress,
			newName:        "",
			newStatus:      StatusDone,
			expectedErr:    ErrEmptyTaskName,
			wantUpdate:     false,
		},
		{
			name:           "Ошибка: прыжок через статус",
			existingStatus: StatusNew,
			newName:        "Новое имя",
			newStatus:      StatusDone,
			expectedErr:    ErrNewTaskStatus,
			wantUpdate:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				TaskToReturn: &Task{Status: tt.existingStatus},
			}
			service := NewService(mockRepo, nil)
			_, err := service.UpdateTask(context.Background(), 1, tt.newName, "Описание", tt.newStatus)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("ожидалась ошибка %v, получена %v", tt.expectedErr, err)
			}
			if mockRepo.UpdateCalled != tt.wantUpdate {
				t.Errorf("UpdateCalled = %v, а ожидалось %v", mockRepo.UpdateCalled, tt.wantUpdate)
			}
			if tt.wantUpdate && mockRepo.UpdatedTask != nil {
				if mockRepo.UpdatedTask.Name != tt.newName {
					t.Errorf("в базу ушло имя %v, а ожидали %v", mockRepo.UpdatedTask.Name, tt.newName)
				}
			}
		})
	}
}
