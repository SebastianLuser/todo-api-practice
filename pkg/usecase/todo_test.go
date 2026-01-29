package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"todo-api/pkg/domain"
	"todo-api/pkg/service"
	"todo-api/pkg/usecase"
	"todo-api/test"
)

const (
	validUUID     = "123e4567-e89b-12d3-a456-426614174000"
	nonExistentID = "00000000-0000-0000-0000-000000000000"
	validTitle    = "Test Todo"
	updatedTitle  = "Updated Title"
)

var fixedTime = time.Date(2026, 1, 28, 10, 30, 0, 0, time.UTC)

func buildValidTodo() domain.Todo {
	return domain.Todo{
		ID:          validUUID,
		Title:       validTitle,
		Description: "This is a test description",
		Status:      domain.StatusPending,
		Priority:    domain.PriorityMedium,
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}
}

func TestTodo_Get_ReturnsListSuccessfully(t *testing.T) {
	expectedTodo := buildValidTodo()
	expectedTodos := []domain.Todo{expectedTodo}
	mock := &test.MockTodoService{
		GetFn: func(ctx context.Context, filters service.Filters) ([]domain.Todo, error) {
			return expectedTodos, nil
		},
	}
	uc := usecase.New(mock)
	input := usecase.ListInput{}

	result, err := uc.Get(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
	if result.Todos[0].ID != expectedTodo.ID {
		t.Errorf("expected ID %s, got %s", expectedTodo.ID, result.Todos[0].ID)
	}
}

func TestTodo_Get_PassesStatusFilterToService(t *testing.T) {
	var capturedFilters service.Filters
	mock := &test.MockTodoService{
		GetFn: func(ctx context.Context, filters service.Filters) ([]domain.Todo, error) {
			capturedFilters = filters
			return []domain.Todo{}, nil
		},
	}
	uc := usecase.New(mock)
	status := domain.StatusCompleted
	input := usecase.ListInput{Status: &status}

	_, err := uc.Get(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if capturedFilters.Status == nil {
		t.Error("expected status filter to be set")
	}
	if *capturedFilters.Status != domain.StatusCompleted {
		t.Errorf("expected status %s, got %s", domain.StatusCompleted, *capturedFilters.Status)
	}
}

func TestTodo_Get_ReturnsErrorWhenServiceFails(t *testing.T) {
	expectedErr := errors.New("database error")
	mock := &test.MockTodoService{
		GetFn: func(ctx context.Context, filters service.Filters) ([]domain.Todo, error) {
			return nil, expectedErr
		},
	}
	uc := usecase.New(mock)

	_, err := uc.Get(context.Background(), usecase.ListInput{})

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestTodo_GetByID_ReturnsTodoSuccessfully(t *testing.T) {
	expectedTodo := buildValidTodo()
	mock := &test.MockTodoService{
		GetByIDFn: func(ctx context.Context, id string) (domain.Todo, error) {
			return expectedTodo, nil
		},
	}
	uc := usecase.New(mock)

	result, err := uc.GetByID(context.Background(), validUUID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Todo.ID != expectedTodo.ID {
		t.Errorf("expected ID %s, got %s", expectedTodo.ID, result.Todo.ID)
	}
	if result.Todo.Title != expectedTodo.Title {
		t.Errorf("expected title %s, got %s", expectedTodo.Title, result.Todo.Title)
	}
}

func TestTodo_GetByID_ReturnsErrTodoNotFound(t *testing.T) {
	mock := &test.MockTodoService{
		GetByIDFn: func(ctx context.Context, id string) (domain.Todo, error) {
			return domain.Todo{}, domain.ErrTodoNotFound
		},
	}
	uc := usecase.New(mock)

	_, err := uc.GetByID(context.Background(), nonExistentID)

	if !errors.Is(err, domain.ErrTodoNotFound) {
		t.Errorf("expected error %v, got %v", domain.ErrTodoNotFound, err)
	}
}

func TestTodo_Create_WithDefaultStatusAndPriority(t *testing.T) {
	var capturedInput service.CreateInput
	expectedTodo := buildValidTodo()
	mock := &test.MockTodoService{
		CreateFn: func(ctx context.Context, input service.CreateInput) (domain.Todo, error) {
			capturedInput = input
			return expectedTodo, nil
		},
	}
	uc := usecase.New(mock)
	input := usecase.CreateInput{Title: validTitle}

	result, err := uc.Create(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if capturedInput.Status != domain.StatusPending {
		t.Errorf("expected status %s, got %s", domain.StatusPending, capturedInput.Status)
	}
	if capturedInput.Priority != domain.PriorityMedium {
		t.Errorf("expected priority %s, got %s", domain.PriorityMedium, capturedInput.Priority)
	}
	if result.Todo.ID != expectedTodo.ID {
		t.Errorf("expected ID %s, got %s", expectedTodo.ID, result.Todo.ID)
	}
}

func TestTodo_Create_WithCustomStatusAndPriority(t *testing.T) {
	var capturedInput service.CreateInput
	mock := &test.MockTodoService{
		CreateFn: func(ctx context.Context, input service.CreateInput) (domain.Todo, error) {
			capturedInput = input
			return buildValidTodo(), nil
		},
	}
	uc := usecase.New(mock)
	status := domain.StatusInProgress
	priority := domain.PriorityHigh
	input := usecase.CreateInput{
		Title:    validTitle,
		Status:   &status,
		Priority: &priority,
	}

	_, err := uc.Create(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if capturedInput.Status != domain.StatusInProgress {
		t.Errorf("expected status %s, got %s", domain.StatusInProgress, capturedInput.Status)
	}
	if capturedInput.Priority != domain.PriorityHigh {
		t.Errorf("expected priority %s, got %s", domain.PriorityHigh, capturedInput.Priority)
	}
}

func TestTodo_Create_ReturnsErrorWhenServiceFails(t *testing.T) {
	expectedErr := errors.New("database error")
	mock := &test.MockTodoService{
		CreateFn: func(ctx context.Context, input service.CreateInput) (domain.Todo, error) {
			return domain.Todo{}, expectedErr
		},
	}
	uc := usecase.New(mock)

	_, err := uc.Create(context.Background(), usecase.CreateInput{Title: validTitle})

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestTodo_Update_Successfully(t *testing.T) {
	expectedTodo := buildValidTodo()
	expectedTodo.Title = updatedTitle
	mock := &test.MockTodoService{
		UpdateFn: func(ctx context.Context, id string, input service.UpdateInput) (domain.Todo, error) {
			return expectedTodo, nil
		},
	}
	uc := usecase.New(mock)
	title := updatedTitle
	input := usecase.UpdateInput{Title: &title}

	result, err := uc.Update(context.Background(), validUUID, input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Todo.Title != updatedTitle {
		t.Errorf("expected title %s, got %s", updatedTitle, result.Todo.Title)
	}
}

func TestTodo_Update_ReturnsErrTodoNotFound(t *testing.T) {
	mock := &test.MockTodoService{
		UpdateFn: func(ctx context.Context, id string, input service.UpdateInput) (domain.Todo, error) {
			return domain.Todo{}, domain.ErrTodoNotFound
		},
	}
	uc := usecase.New(mock)
	title := updatedTitle
	input := usecase.UpdateInput{Title: &title}

	_, err := uc.Update(context.Background(), nonExistentID, input)

	if !errors.Is(err, domain.ErrTodoNotFound) {
		t.Errorf("expected error %v, got %v", domain.ErrTodoNotFound, err)
	}
}

func TestTodo_Delete_Successfully(t *testing.T) {
	var capturedID string
	mock := &test.MockTodoService{
		DeleteFn: func(ctx context.Context, id string) error {
			capturedID = id
			return nil
		},
	}
	uc := usecase.New(mock)

	err := uc.Delete(context.Background(), validUUID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if capturedID != validUUID {
		t.Errorf("expected ID %s, got %s", validUUID, capturedID)
	}
}

func TestTodo_Delete_ReturnsErrTodoNotFound(t *testing.T) {
	mock := &test.MockTodoService{
		DeleteFn: func(ctx context.Context, id string) error {
			return domain.ErrTodoNotFound
		},
	}
	uc := usecase.New(mock)

	err := uc.Delete(context.Background(), nonExistentID)

	if !errors.Is(err, domain.ErrTodoNotFound) {
		t.Errorf("expected error %v, got %v", domain.ErrTodoNotFound, err)
	}
}
