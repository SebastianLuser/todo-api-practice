package service

import (
	"context"
	"database/sql"
	_ "embed"

	"todo-api/pkg/domain"
)

//go:embed sql/select/get_todos.sql
var getTodosQuery string

type (
	Filters struct {
		Status   *domain.Status
		Priority *domain.Priority
	}

	postgresService struct {
		db *sql.DB
	}

	Todo interface {
		Get(ctx context.Context, filters Filters) ([]domain.Todo, error)
	}
)

func New(db *sql.DB) Todo {
	return &postgresService{db: db}
}

func (s *postgresService) Get(ctx context.Context, filters Filters) ([]domain.Todo, error) {
	var statusFilter, priorityFilter *string

	if filters.Status != nil {
		v := string(*filters.Status)
		statusFilter = &v
	}
	if filters.Priority != nil {
		v := string(*filters.Priority)
		priorityFilter = &v
	}

	rows, err := s.db.QueryContext(ctx, getTodosQuery, statusFilter, priorityFilter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []domain.Todo
	for rows.Next() {
		var todo domain.Todo
		var description sql.NullString

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&description,
			&todo.Status,
			&todo.Priority,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			todo.Description = description.String
		}

		todos = append(todos, todo)
	}

	return todos, rows.Err()
}
