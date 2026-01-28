package domain

import "errors"

var (
	ErrTodoNotFound       = errors.New("todo not found")
	ErrInvalidStatus      = errors.New("invalid status: must be pending, in_progress, or completed")
	ErrInvalidPriority    = errors.New("invalid priority: must be low, medium, or high")
	ErrInvalidTitle       = errors.New("invalid title: must be between 1 and 100 characters")
	ErrInvalidDescription = errors.New("invalid description: must be at most 500 characters")
	ErrInvalidID          = errors.New("invalid id: must be a valid UUID")
	ErrEmptyUpdateRequest = errors.New("update request must contain at least one field")
)
