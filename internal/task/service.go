package task

type Service struct {
	repo TaskRepository
}

func NewService(repo TaskRepository) *Service {
	return &Service{repo: repo}
}
