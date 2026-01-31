package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
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
		INSERT INTO tasks (name, description, created, status, group_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		task.Name,
		task.Description,
		task.Created,
		task.Status,
		task.GroupID,
	).Scan(&task.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("postgres.Add: insert task: %w", ErrGroupNotFound)
		}
		return fmt.Errorf("postgres.Add: insert task: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id int) (*Task, error) {
	var t Task
	query := `
		SELECT id, name, description, created, status, group_id 
		FROM tasks 
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Created,
		&t.Status,
		&t.GroupID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("postgres.GetById: scan task id=%d: %w", id, err)
	}
	return &t, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, groupId *int) ([]Task, error) {
	query := `
	SELECT 
		t.id, t.name, t.description, t.created, t.status, t.group_id,
		g.name as group_name
	FROM tasks t
	JOIN LEFT groups g ON t.group_id = g.id
	`
	var args []any
	conditions := []string{}
	if groupId != nil {
		args = append(args, *groupId)
		conditions = append(conditions, fmt.Sprintf("group_id = $%d", len(args)))

	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres.GetAll: query tasks: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Created, &t.Status, &t.GroupID, &t.GroupName)
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
		SET name = $1, description = $2, status = $3, group_id = $4 
		WHERE id = $5
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		task.Name,
		task.Description,
		task.Status,
		task.GroupID,
		task.ID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("postgres.Update: insert task: %w", ErrGroupNotFound)
		}
		return fmt.Errorf("failed to update task: %w", err)
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
