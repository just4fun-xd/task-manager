package task

import (
	"context"
	"errors"
	"testing"
)

type MockRepository struct {
	AddCalled        bool
	UpdateCalled     bool
	TaskToReturn     *Task
	UpdatedTask      *Task
	GetAllCalledWith *int
	GetAllCalled     bool
}

type MockGroupRepository struct {
	ErrorToReturn error
}

func (m *MockRepository) Add(ctx context.Context, task *Task) error {
	m.AddCalled = true
	return nil
}

func (m *MockRepository) GetAll(ctx context.Context, groupId *int) ([]Task, error) {
	m.GetAllCalled = true
	m.GetAllCalledWith = groupId
	return nil, nil
}
func (m *MockRepository) GetById(ctx context.Context, id int) (*Task, error) {
	return m.TaskToReturn, nil
}
func (m *MockRepository) Update(ctx context.Context, task *Task) error {
	m.UpdatedTask = task
	m.UpdateCalled = true
	return nil
}
func (m *MockRepository) Delete(ctx context.Context, id int) error { return nil }

func (m *MockGroupRepository) Add(ctx context.Context, group *Group) error { return nil }
func (m *MockGroupRepository) GetAll(ctx context.Context) ([]Group, error) { return nil, nil }
func (m *MockGroupRepository) GetById(ctx context.Context, id int) (*Group, error) {
	return nil, m.ErrorToReturn
}
func (m *MockGroupRepository) Update(ctx context.Context, group *Group) error { return nil }
func (m *MockGroupRepository) Delete(ctx context.Context, id int) error       { return nil }

/*
type GroupRepository interface {
	Add(ctx context.Context, group *Group) error +
	GetAll(ctx context.Context) ([]Group, error) +
	GetById(ctx context.Context, id int) (*Group, error) +
	Update(ctx context.Context, group *Group) error +
	Delete(ctx context.Context, id int) error
}
*/

func TestCreateTask_EmptyName(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo, nil)
	_, err := service.CreateTask(context.Background(), "", "Описание", nil)

	if !errors.Is(err, ErrEmptyTaskName) {
		t.Errorf("ожидалась ошибка %v, получена %v", ErrEmptyTaskName, err)
	}
	if mockRepo.AddCalled {
		t.Error("репозиторий не должен был вызваться при пустом имени")
	}
}

func TestGetAllTasks_GroupNotFound(t *testing.T) {
	mockRepo := &MockRepository{}
	mockGroupRepository := &MockGroupRepository{
		ErrorToReturn: errors.New("group not found"),
	}
	service := NewService(mockRepo, mockGroupRepository)
	id := 10
	tasks, err := service.GetAllTasks(context.Background(), &id)
	if !errors.Is(err, mockGroupRepository.ErrorToReturn) {
		t.Errorf("ожидалось error = %v, получена %v", mockGroupRepository.ErrorToReturn, err)
	}
	if len(tasks) != 0 {
		t.Errorf("ожидалось 0 задач, получено %d", len(tasks))
	}
	if mockRepo.GetAllCalled == true {
		t.Error("Метод GetAll был вызван")
	}

}

func TestGetAllTasks_Filtering(t *testing.T) {
	mockRepo := &MockRepository{}
	mockGroupRepository := &MockGroupRepository{}
	service := NewService(mockRepo, mockGroupRepository)
	groupId := 5
	_, _ = service.GetAllTasks(context.Background(), &groupId)
	if mockRepo.GetAllCalledWith == nil {
		t.Error("ожидалось groupId != nil ")
	}
	if *mockRepo.GetAllCalledWith != groupId {
		t.Errorf("ожидалось groupId = %d, получена %d", groupId, *mockRepo.GetAllCalledWith)
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
			_, err := service.UpdateTask(context.Background(), 1, tt.newName, "Описание", tt.newStatus, nil)
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
