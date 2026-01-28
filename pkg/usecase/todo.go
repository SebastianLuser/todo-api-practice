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

	usecase struct {
		service service.Service
	}

	Usecase interface {
		List(ctx context.Context, input ListInput) (ListOutput, error)
	}
)

func New(svc service.Service) Usecase {
	return &usecase{
		service: svc,
	}
}

func (u *usecase) List(ctx context.Context, input ListInput) (ListOutput, error) {
	filters := service.Filters{
		Status:   input.Status,
		Priority: input.Priority,
	}

	todos, err := u.service.List(ctx, filters)
	if err != nil {
		return ListOutput{}, err
	}

	return ListOutput{
		Todos: todos,
		Total: len(todos),
	}, nil
}
