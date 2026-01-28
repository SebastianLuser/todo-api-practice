package service

import (
	"context"
	"database/sql"
	_ "embed"

	"todo-api/pkg/domain"
)

//go:embed sql/select/get_todos.sql
var getTodosQuery string

//go:embed sql/select/get_todo_by_id.sql
var getTodoByIDQuery string

//go:embed sql/insert/create_todo.sql
var createTodoQuery string

//go:embed sql/update/update_todo.sql
var updateTodoQuery string

//go:embed sql/delete/delete_todo.sql
var deleteTodoQuery string

type (
	Filters struct {
		Status   *domain.Status
		Priority *domain.Priority
	}

	postgresService struct {
		db *sql.DB
	}

	CreateInput struct {
		Title       string
		Description *string
		Status      domain.Status
		Priority    domain.Priority
	}

	UpdateInput struct {
		Title       *string
		Description *string
		Status      *domain.Status
		Priority    *domain.Priority
	}

	Todo interface {
		Get(ctx context.Context, filters Filters) ([]domain.Todo, error)
		GetByID(ctx context.Context, id string) (domain.Todo, error)
		Create(ctx context.Context, input CreateInput) (domain.Todo, error)
		Update(ctx context.Context, id string, input UpdateInput) (domain.Todo, error)
		Delete(ctx context.Context, id string) error
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

func (s *postgresService) GetByID(ctx context.Context, id string) (domain.Todo, error) {
	row := s.db.QueryRowContext(ctx, getTodoByIDQuery, id)

	var todo domain.Todo
	var description sql.NullString

	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&description,
		&todo.Status,
		&todo.Priority,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Todo{}, domain.ErrTodoNotFound
		}
		return domain.Todo{}, err
	}

	if description.Valid {
		todo.Description = description.String
	}

	return todo, nil
}

func (s *postgresService) Create(ctx context.Context, input CreateInput) (domain.Todo, error) {
	var todo domain.Todo
	var description sql.NullString

	if input.Description != nil {
		description = sql.NullString{String: *input.Description, Valid: true}
	}

	row := s.db.QueryRowContext(
		ctx,
		createTodoQuery,
		input.Title,
		description,
		input.Status,
		input.Priority,
	)

	var descResult sql.NullString
	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&descResult,
		&todo.Status,
		&todo.Priority,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return domain.Todo{}, err
	}

	if descResult.Valid {
		todo.Description = descResult.String
	}

	return todo, nil
}

func (s *postgresService) Update(ctx context.Context, id string, input UpdateInput) (domain.Todo, error) {
	var title, description, status, priority *string

	if input.Title != nil {
		title = input.Title
	}
	if input.Description != nil {
		description = input.Description
	}
	if input.Status != nil {
		v := string(*input.Status)
		status = &v
	}
	if input.Priority != nil {
		v := string(*input.Priority)
		priority = &v
	}

	row := s.db.QueryRowContext(
		ctx,
		updateTodoQuery,
		id,
		title,
		description,
		status,
		priority,
	)

	var todo domain.Todo
	var descResult sql.NullString

	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&descResult,
		&todo.Status,
		&todo.Priority,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Todo{}, domain.ErrTodoNotFound
		}
		return domain.Todo{}, err
	}

	if descResult.Valid {
		todo.Description = descResult.String
	}

	return todo, nil
}

func (s *postgresService) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, deleteTodoQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrTodoNotFound
	}

	return nil
}
