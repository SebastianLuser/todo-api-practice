package service

import (
	"context"

	"todo-api/pkg/domain"
)

type (
	Filters struct {
		Status   *domain.Status
		Priority *domain.Priority
	}

	inMemoryService struct {
		todos map[string]domain.Todo
	}

	Service interface {
		List(ctx context.Context, filters Filters) ([]domain.Todo, error)
	}
)

func New() Service {
	return &inMemoryService{
		todos: make(map[string]domain.Todo),
	}
}

func (s *inMemoryService) List(ctx context.Context, filters Filters) ([]domain.Todo, error) {
	var result []domain.Todo

	for _, todo := range s.todos {
		if filters.Status != nil && todo.Status != *filters.Status {
			continue
		}
		if filters.Priority != nil && todo.Priority != *filters.Priority {
			continue
		}
		result = append(result, todo)
	}

	return result, nil
}
