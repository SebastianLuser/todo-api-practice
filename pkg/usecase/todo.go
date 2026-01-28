package usecase

import (
	"context"

	"todo-api/pkg/domain"
	"todo-api/pkg/service"
)

type (
	ListInput struct {
		Status   *domain.Status
		Priority *domain.Priority
	}

	ListOutput struct {
		Todos []domain.Todo
		Total int
	}

	Todo struct {
		service service.Todo
	}
)

func New(svc service.Todo) *Todo {
	return &Todo{
		service: svc,
	}
}

func (u *Todo) Get(ctx context.Context, input ListInput) (ListOutput, error) {
	filters := service.Filters{
		Status:   input.Status,
		Priority: input.Priority,
	}

	todos, err := u.service.Get(ctx, filters)
	if err != nil {
		return ListOutput{}, err
	}

	return ListOutput{
		Todos: todos,
		Total: len(todos),
	}, nil
}
