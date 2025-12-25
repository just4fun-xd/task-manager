package task

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Add(ctx context.Context, task *Task) error {
	query := `
		INSERT INTO tasks (name, description, created, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		task.Name,
		task.Description,
		task.Created,
		task.Status,
	).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("postgres.Add: insert task: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id int) (*Task, error) {
	var t Task
	query := `
		SELECT id, name, description, created, status 
		FROM tasks 
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Created,
		&t.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("postgres.GetById: scan task id=%d: %w", id, err)
	}
	return &t, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]Task, error) {
	query := `SELECT id, name, description, created, status FROM tasks`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres.GetAll: query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Created, &t.Status)
		if err != nil {
			return nil, fmt.Errorf("postgres.GetAll: scan task row: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.GetAll: rows iteration: %w", err)
	}
	return tasks, nil
}

func (r *PostgresRepository) Update(ctx context.Context, task *Task) error {
	query := `
		UPDATE tasks
		SET name = $1, description = $2, status = $3 
		WHERE id = $4
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		task.Name,
		task.Description,
		task.Status,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("tasks not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int) error {
	query := `
	DELETE FROM tasks WHERE id = $1`
	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}
