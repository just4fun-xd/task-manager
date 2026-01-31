package task

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type MockRepository struct {
	AddCalled        bool
	UpdateCalled     bool
	AddedTask        *Task
	TaskToReturn     *Task
	UpdatedTask      *Task
	GetAllCalledWith *int
	GetAllCalled     bool
}

type MockGroupRepository struct {
	AddCalled     bool
	AddedGroup    *Group
	GroupToReturn *Group
	ErrorToReturn error
}

func (m *MockRepository) Add(ctx context.Context, task *Task) error {
	m.AddCalled = true
	m.AddedTask = task
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

func (m *MockGroupRepository) Add(ctx context.Context, group *Group) error {
	m.AddCalled = true
	m.AddedGroup = group
	return m.ErrorToReturn
}
func (m *MockGroupRepository) GetAll(ctx context.Context) ([]Group, error) { return nil, nil }
func (m *MockGroupRepository) GetById(ctx context.Context, id int) (*Group, error) {
	return m.GroupToReturn, m.ErrorToReturn
}
func (m *MockGroupRepository) Update(ctx context.Context, group *Group) error { return nil }
func (m *MockGroupRepository) Delete(ctx context.Context, id int) error       { return nil }

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

func TestCreateGroup_DuplicateName(t *testing.T) {
	mockRepo := &MockRepository{}
	mockGroupRepo := &MockGroupRepository{
		ErrorToReturn: ErrNotUniqGroup,
	}
	service := NewService(mockRepo, mockGroupRepo)
	groupName := "DuplicateGroupName"
	_, err := service.CreateGroup(context.Background(), groupName)
	if !errors.Is(err, ErrNotUniqGroup) {
		t.Errorf("ожидалась ошибка %v, получена %v", ErrNotUniqGroup, err)
	}
	if !mockGroupRepo.AddCalled {
		t.Error("репозиторий групп не был вызван")
	}
	if mockGroupRepo.AddedGroup == nil {
		t.Fatalf("в репозиторий не была передана группа (nil)")
	}
	if mockGroupRepo.AddedGroup.Name != groupName {
		t.Errorf("ожидалось имя группы %q, получили %q", groupName, mockGroupRepo.AddedGroup.Name)
	}
}

func TestCreateGroup_Success(t *testing.T) {
	mockRepo := &MockRepository{}
	mockGroupRepo := &MockGroupRepository{}
	service := NewService(mockRepo, mockGroupRepo)

	groupName := "Работа"

	_, err := service.CreateGroup(context.Background(), groupName)
	if err != nil {
		t.Errorf("не ожидалось ошибки, получена: %v", err)
	}
	if !mockGroupRepo.AddCalled {
		t.Error("репозиторий групп не был вызван")
	}
	if mockGroupRepo.AddedGroup == nil {
		t.Fatalf("в репозиторий не была передана группа (nil)")
	}
	if mockGroupRepo.AddedGroup.Name != groupName {
		t.Errorf("ожидалось имя группы %q, получили %q", groupName, mockGroupRepo.AddedGroup.Name)
	}

}

func TestCreateTaskWithGroup(t *testing.T) {

	mockRepo := &MockRepository{}
	mockGroupRepo := &MockGroupRepository{GroupToReturn: &Group{ID: 1, Name: "Test"}}
	name := "Купить хлеб"
	desc := "Бородинский"
	groupID := 1
	service := NewService(mockRepo, mockGroupRepo)
	task, err := service.CreateTask(context.Background(), name, desc, &groupID)
	if err != nil {
		t.Fatalf("ожидалось error = nil, получено: %v", err)
	}
	if task.Name != name || task.Description != desc || task.GroupID == nil || *task.GroupID != groupID {
		t.Errorf("Данные задачи не совпадают!")
		t.Errorf("Ожидалось: Name=%q, Desc=%q, GroupID=%d", name, desc, groupID)

		gotID := 0
		if task.GroupID != nil {
			gotID = *task.GroupID
		}
		t.Errorf("Получено:  Name=%q, Desc=%q, GroupID=%d", task.Name, task.Description, gotID)
	}
	if !mockRepo.AddCalled {
		t.Fatal("Метод Add репозитория не был вызван!")
	}
	if mockRepo.AddedTask == nil {
		t.Fatal("репозиторий был вызван, но задача не была передана (nil)")
	}
	gotID := "nil"
	if mockRepo.AddedTask.GroupID != nil {
		gotID = fmt.Sprintf("%d", *mockRepo.AddedTask.GroupID)
	}
	if mockRepo.AddedTask.Name != name || mockRepo.AddedTask.GroupID == nil || *mockRepo.AddedTask.GroupID != groupID {
		t.Errorf("В репозиторий ушли неверные данные! Ожидали Name=%s, ID=%d. Получили Name=%s, ID=%s",
			name, groupID, mockRepo.AddedTask.Name, gotID)
	}
}

func TestGetAllTasks_GroupNotFound(t *testing.T) {
	mockRepo := &MockRepository{}
	mockGroupRepo := &MockGroupRepository{
		ErrorToReturn: errors.New("group not found"),
	}
	service := NewService(mockRepo, mockGroupRepo)
	id := 10
	tasks, err := service.GetAllTasks(context.Background(), &id)
	if !errors.Is(err, mockGroupRepo.ErrorToReturn) {
		t.Errorf("ожидалось error = %v, получена %v", mockGroupRepo.ErrorToReturn, err)
	}
	if len(tasks) != 0 {
		t.Errorf("ожидалось 0 задач, получено %d", len(tasks))
	}
	if mockRepo.GetAllCalled {
		t.Error("Метод GetAll был вызван")
	}

}

func TestGetAllTasks_Filtering(t *testing.T) {
	mockRepo := &MockRepository{}
	mockGroupRepo := &MockGroupRepository{}
	service := NewService(mockRepo, mockGroupRepo)
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
			const testDesc = "Описание задачи"
			_, err := service.UpdateTask(context.Background(), 1, tt.newName, testDesc, tt.newStatus, nil)
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
				if mockRepo.UpdatedTask.Description != testDesc {
					t.Errorf("в базу ушло описание: %v, а ожидали: %v", mockRepo.UpdatedTask.Description, testDesc)
				}
				if mockRepo.UpdatedTask.Status != tt.newStatus {
					t.Errorf("в базу ушел статус: %v, а ожидали: %v", mockRepo.UpdatedTask.Status, tt.newStatus)
				}
			}
		})
	}
}
