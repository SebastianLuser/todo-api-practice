package usecase_test

import (
	"context"
	"errors"
	"testing"

	"todo-api/pkg/domain"
	"todo-api/pkg/service"
	"todo-api/pkg/usecase"
	"todo-api/test"
)

// mockTodoService implements service.Todo for testing
type mockTodoService struct {
	getFn      func(ctx context.Context, filters service.Filters) ([]domain.Todo, error)
	getByIDFn  func(ctx context.Context, id string) (domain.Todo, error)
	createFn   func(ctx context.Context, input service.CreateInput) (domain.Todo, error)
	updateFn   func(ctx context.Context, id string, input service.UpdateInput) (domain.Todo, error)
	deleteFn   func(ctx context.Context, id string) error
}

func (m *mockTodoService) Get(ctx context.Context, filters service.Filters) ([]domain.Todo, error) {
	return m.getFn(ctx, filters)
}

func (m *mockTodoService) GetByID(ctx context.Context, id string) (domain.Todo, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockTodoService) Create(ctx context.Context, input service.CreateInput) (domain.Todo, error) {
	return m.createFn(ctx, input)
}

func (m *mockTodoService) Update(ctx context.Context, id string, input service.UpdateInput) (domain.Todo, error) {
	return m.updateFn(ctx, id, input)
}

func (m *mockTodoService) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

func TestTodo_Get(t *testing.T) {
	t.Run("should return todos list successfully", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		expectedTodo := test.BuildValidTodo()
		expectedTodos := []domain.Todo{expectedTodo}

		mock := &mockTodoService{
			getFn: func(ctx context.Context, filters service.Filters) ([]domain.Todo, error) {
				return expectedTodos, nil
			},
		}
		uc := usecase.New(mock)
		input := usecase.ListInput{}

		// Act
		result, err := uc.Get(context.Background(), input)

		// Assert
		assert.NoError(err)
		assert.Equal(1, result.Total)
		assert.Equal(expectedTodo.ID, result.Todos[0].ID)
	})

	t.Run("should pass status filter to service", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		var capturedFilters service.Filters

		mock := &mockTodoService{
			getFn: func(ctx context.Context, filters service.Filters) ([]domain.Todo, error) {
				capturedFilters = filters
				return []domain.Todo{}, nil
			},
		}
		uc := usecase.New(mock)
		status := domain.StatusCompleted
		input := usecase.ListInput{Status: &status}

		// Act
		_, err := uc.Get(context.Background(), input)

		// Assert
		assert.NoError(err)
		assert.NotNil(capturedFilters.Status)
		assert.Equal(domain.StatusCompleted, *capturedFilters.Status)
	})

	t.Run("should return error when service fails", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		expectedErr := errors.New("database error")

		mock := &mockTodoService{
			getFn: func(ctx context.Context, filters service.Filters) ([]domain.Todo, error) {
				return nil, expectedErr
			},
		}
		uc := usecase.New(mock)

		// Act
		_, err := uc.Get(context.Background(), usecase.ListInput{})

		// Assert
		assert.ErrorIs(err, expectedErr)
	})
}

func TestTodo_GetByID(t *testing.T) {
	t.Run("should return todo by id successfully", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		expectedTodo := test.BuildValidTodo()

		mock := &mockTodoService{
			getByIDFn: func(ctx context.Context, id string) (domain.Todo, error) {
				return expectedTodo, nil
			},
		}
		uc := usecase.New(mock)

		// Act
		result, err := uc.GetByID(context.Background(), test.ValidUUID)

		// Assert
		assert.NoError(err)
		assert.Equal(expectedTodo.ID, result.Todo.ID)
		assert.Equal(expectedTodo.Title, result.Todo.Title)
	})

	t.Run("should return ErrTodoNotFound when todo does not exist", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)

		mock := &mockTodoService{
			getByIDFn: func(ctx context.Context, id string) (domain.Todo, error) {
				return domain.Todo{}, domain.ErrTodoNotFound
			},
		}
		uc := usecase.New(mock)

		// Act
		_, err := uc.GetByID(context.Background(), test.NonExistentID)

		// Assert
		assert.ErrorIs(err, domain.ErrTodoNotFound)
	})
}

func TestTodo_Create(t *testing.T) {
	t.Run("should create todo with default status and priority", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		var capturedInput service.CreateInput
		expectedTodo := test.BuildValidTodo()

		mock := &mockTodoService{
			createFn: func(ctx context.Context, input service.CreateInput) (domain.Todo, error) {
				capturedInput = input
				return expectedTodo, nil
			},
		}
		uc := usecase.New(mock)
		input := usecase.CreateInput{Title: test.ValidTitle}

		// Act
		result, err := uc.Create(context.Background(), input)

		// Assert
		assert.NoError(err)
		assert.Equal(domain.StatusPending, capturedInput.Status)
		assert.Equal(domain.PriorityMedium, capturedInput.Priority)
		assert.Equal(expectedTodo.ID, result.Todo.ID)
	})

	t.Run("should create todo with custom status and priority", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		var capturedInput service.CreateInput

		mock := &mockTodoService{
			createFn: func(ctx context.Context, input service.CreateInput) (domain.Todo, error) {
				capturedInput = input
				return test.BuildValidTodo(), nil
			},
		}
		uc := usecase.New(mock)
		status := domain.StatusInProgress
		priority := domain.PriorityHigh
		input := usecase.CreateInput{
			Title:    test.ValidTitle,
			Status:   &status,
			Priority: &priority,
		}

		// Act
		_, err := uc.Create(context.Background(), input)

		// Assert
		assert.NoError(err)
		assert.Equal(domain.StatusInProgress, capturedInput.Status)
		assert.Equal(domain.PriorityHigh, capturedInput.Priority)
	})

	t.Run("should return error when service fails", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		expectedErr := errors.New("database error")

		mock := &mockTodoService{
			createFn: func(ctx context.Context, input service.CreateInput) (domain.Todo, error) {
				return domain.Todo{}, expectedErr
			},
		}
		uc := usecase.New(mock)

		// Act
		_, err := uc.Create(context.Background(), usecase.CreateInput{Title: test.ValidTitle})

		// Assert
		assert.ErrorIs(err, expectedErr)
	})
}

func TestTodo_Update(t *testing.T) {
	t.Run("should update todo successfully", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		expectedTodo := test.BuildValidTodo()
		expectedTodo.Title = test.UpdatedTitle

		mock := &mockTodoService{
			updateFn: func(ctx context.Context, id string, input service.UpdateInput) (domain.Todo, error) {
				return expectedTodo, nil
			},
		}
		uc := usecase.New(mock)
		title := test.UpdatedTitle
		input := usecase.UpdateInput{Title: &title}

		// Act
		result, err := uc.Update(context.Background(), test.ValidUUID, input)

		// Assert
		assert.NoError(err)
		assert.Equal(test.UpdatedTitle, result.Todo.Title)
	})

	t.Run("should return ErrTodoNotFound when todo does not exist", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)

		mock := &mockTodoService{
			updateFn: func(ctx context.Context, id string, input service.UpdateInput) (domain.Todo, error) {
				return domain.Todo{}, domain.ErrTodoNotFound
			},
		}
		uc := usecase.New(mock)
		title := test.UpdatedTitle
		input := usecase.UpdateInput{Title: &title}

		// Act
		_, err := uc.Update(context.Background(), test.NonExistentID, input)

		// Assert
		assert.ErrorIs(err, domain.ErrTodoNotFound)
	})
}

func TestTodo_Delete(t *testing.T) {
	t.Run("should delete todo successfully", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		var capturedID string

		mock := &mockTodoService{
			deleteFn: func(ctx context.Context, id string) error {
				capturedID = id
				return nil
			},
		}
		uc := usecase.New(mock)

		// Act
		err := uc.Delete(context.Background(), test.ValidUUID)

		// Assert
		assert.NoError(err)
		assert.Equal(test.ValidUUID, capturedID)
	})

	t.Run("should return ErrTodoNotFound when todo does not exist", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)

		mock := &mockTodoService{
			deleteFn: func(ctx context.Context, id string) error {
				return domain.ErrTodoNotFound
			},
		}
		uc := usecase.New(mock)

		// Act
		err := uc.Delete(context.Background(), test.NonExistentID)

		// Assert
		assert.ErrorIs(err, domain.ErrTodoNotFound)
	})
}
