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

	GetByIDOutput struct {
		Todo domain.Todo
	}

	CreateInput struct {
		Title       string
		Description *string
		Status      *domain.Status
		Priority    *domain.Priority
	}

	CreateOutput struct {
		Todo domain.Todo
	}

	UpdateInput struct {
		Title       *string
		Description *string
		Status      *domain.Status
		Priority    *domain.Priority
	}

	UpdateOutput struct {
		Todo domain.Todo
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

func (u *Todo) GetByID(ctx context.Context, id string) (GetByIDOutput, error) {
	todo, err := u.service.GetByID(ctx, id)
	if err != nil {
		return GetByIDOutput{}, err
	}

	return GetByIDOutput{Todo: todo}, nil
}

func (u *Todo) Create(ctx context.Context, input CreateInput) (CreateOutput, error) {
	status := domain.StatusPending
	if input.Status != nil {
		status = *input.Status
	}

	priority := domain.PriorityMedium
	if input.Priority != nil {
		priority = *input.Priority
	}

	svcInput := service.CreateInput{
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
		Priority:    priority,
	}

	todo, err := u.service.Create(ctx, svcInput)
	if err != nil {
		return CreateOutput{}, err
	}

	return CreateOutput{Todo: todo}, nil
}

func (u *Todo) Update(ctx context.Context, id string, input UpdateInput) (UpdateOutput, error) {
	svcInput := service.UpdateInput{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
	}

	todo, err := u.service.Update(ctx, id, svcInput)
	if err != nil {
		return UpdateOutput{}, err
	}

	return UpdateOutput{Todo: todo}, nil
}

func (u *Todo) Delete(ctx context.Context, id string) error {
	return u.service.Delete(ctx, id)
}
