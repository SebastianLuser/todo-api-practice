package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"todo-api/pkg/domain"
	"todo-api/test"
)

func TestDomainErrors_ErrorsIs(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"ErrTodoNotFound", domain.ErrTodoNotFound},
		{"ErrInvalidStatus", domain.ErrInvalidStatus},
		{"ErrInvalidPriority", domain.ErrInvalidPriority},
		{"ErrInvalidTitle", domain.ErrInvalidTitle},
		{"ErrInvalidDescription", domain.ErrInvalidDescription},
		{"ErrInvalidID", domain.ErrInvalidID},
		{"ErrEmptyUpdateRequest", domain.ErrEmptyUpdateRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			assert := test.NewAssert(t)

			// Act & Assert
			assert.ErrorIs(tc.err, tc.err)
		})
	}
}

func TestDomainErrors_WrappedErrorsIs(t *testing.T) {
	t.Run("wrapped ErrTodoNotFound should be identifiable", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		wrappedErr := fmt.Errorf("service error: %w", domain.ErrTodoNotFound)

		// Act & Assert
		assert.ErrorIs(wrappedErr, domain.ErrTodoNotFound)
	})

	t.Run("wrapped ErrInvalidID should be identifiable", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		wrappedErr := fmt.Errorf("validation failed: %w", domain.ErrInvalidID)

		// Act & Assert
		assert.ErrorIs(wrappedErr, domain.ErrInvalidID)
	})

	t.Run("double wrapped error should be identifiable", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		firstWrap := fmt.Errorf("layer 1: %w", domain.ErrInvalidStatus)
		secondWrap := fmt.Errorf("layer 2: %w", firstWrap)

		// Act & Assert
		assert.ErrorIs(secondWrap, domain.ErrInvalidStatus)
	})
}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func TestDomainErrors_ErrorsAs(t *testing.T) {
	t.Run("wrapped custom error should be extractable with ErrorAs", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		originalErr := &customError{msg: "custom error"}
		wrappedErr := fmt.Errorf("wrapped: %w", originalErr)
		var target *customError

		// Act
		result := errors.As(wrappedErr, &target)

		// Assert
		assert.True(result)
		assert.Equal("custom error", target.msg)
	})
}

func TestDomainErrors_Messages(t *testing.T) {
	testCases := []struct {
		name            string
		err             error
		expectedMessage string
	}{
		{
			name:            "ErrTodoNotFound message",
			err:             domain.ErrTodoNotFound,
			expectedMessage: "todo not found",
		},
		{
			name:            "ErrInvalidStatus message",
			err:             domain.ErrInvalidStatus,
			expectedMessage: "invalid status: must be pending, in_progress, or completed",
		},
		{
			name:            "ErrInvalidPriority message",
			err:             domain.ErrInvalidPriority,
			expectedMessage: "invalid priority: must be low, medium, or high",
		},
		{
			name:            "ErrInvalidTitle message",
			err:             domain.ErrInvalidTitle,
			expectedMessage: "invalid title: must be between 1 and 100 characters",
		},
		{
			name:            "ErrInvalidDescription message",
			err:             domain.ErrInvalidDescription,
			expectedMessage: "invalid description: must be at most 500 characters",
		},
		{
			name:            "ErrInvalidID message",
			err:             domain.ErrInvalidID,
			expectedMessage: "invalid id: must be a valid UUID",
		},
		{
			name:            "ErrEmptyUpdateRequest message",
			err:             domain.ErrEmptyUpdateRequest,
			expectedMessage: "update request must contain at least one field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			assert := test.NewAssert(t)

			// Act
			message := tc.err.Error()

			// Assert
			assert.Equal(tc.expectedMessage, message)
		})
	}
}

func TestDomainErrors_NotEqual(t *testing.T) {
	t.Run("different errors should not match with ErrorIs", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)

		// Act
		result := errors.Is(domain.ErrTodoNotFound, domain.ErrInvalidID)

		// Assert
		assert.False(result)
	})

	t.Run("ErrInvalidStatus should not match ErrInvalidPriority", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)

		// Act
		result := errors.Is(domain.ErrInvalidStatus, domain.ErrInvalidPriority)

		// Assert
		assert.False(result)
	})
}
