package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresGroupRepository struct {
	db *sql.DB
}

func NewPostgresGroupRepository(db *sql.DB) *PostgresGroupRepository {
	return &PostgresGroupRepository{
		db: db,
	}
}

func (r *PostgresGroupRepository) Add(ctx context.Context, group *Group) error {
	query := `
		INSERT INTO groups (name)
		VALUES ($1)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, group.Name).Scan(&group.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("postgres.Add group: %w: ", ErrNotUniqGroup)
		}
		return fmt.Errorf("postgres.Add into groups: %w", err)
	}
	return nil
}

func (r *PostgresGroupRepository) GetById(ctx context.Context, id int) (*Group, error) {
	var group Group
	query := `SELECT id, name FROM groups WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&group.ID, &group.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrGroupNotFound
		}
		return nil, fmt.Errorf("postgres.GetById scan group id=%d: %w", id, err)
	}
	return &group, err
}

func (r *PostgresGroupRepository) GetAll(ctx context.Context) ([]Group, error) {
	query := `SELECT id, name FROM groups`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres.GetAll query group %w", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.ID, &group.Name)
		if err != nil {
			return nil, fmt.Errorf("postgres.GetAll row group %w", err)
		}
		groups = append(groups, group)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.GetAll row iteration %w", err)
	}
	return groups, nil
}

func (r *PostgresGroupRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM groups WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("postgres.Delete group: %w", ErrGroupHasTasks)
		}
		return fmt.Errorf("failed to delete group: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get group rows affected %w", err)
	}
	if rows == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func (r *PostgresGroupRepository) Update(ctx context.Context, group *Group) error {
	query := `
		UPDATE groups
		SET name = $1
		WHERE id = $2
	`
	result, err := r.db.ExecContext(ctx, query, group.Name, group.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("postgres.Update group: %w", ErrNotUniqGroup)
		}
		return fmt.Errorf("failed to update group: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get group rows affected %w", err)
	}
	if rows == 0 {
		return ErrGroupNotFound
	}
	return nil
}
