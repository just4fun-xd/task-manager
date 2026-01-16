package task

import "context"

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GroupRepository interface {
	Add(ctx context.Context, group *Group) error
	GetAll(ctx context.Context) ([]Group, error)
	GetById(ctx context.Context, id int) (*Group, error)
	Update(ctx context.Context, group *Group) error
	Delete(ctx context.Context, id int) error
}
