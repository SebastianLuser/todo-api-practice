package domain

import "errors"

var (
	ErrTodoNotFound    = errors.New("todo not found")
	ErrInvalidStatus   = errors.New("invalid status: must be pending, in_progress, or completed")
	ErrInvalidPriority = errors.New("invalid priority: must be low, medium, or high")
	ErrInvalidTitle    = errors.New("invalid title: must be between 1 and 100 characters")
)
