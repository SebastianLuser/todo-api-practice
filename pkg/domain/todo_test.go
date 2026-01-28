package domain_test

import (
	"testing"

	"todo-api/pkg/domain"
	"todo-api/test"
)

func TestStatus_IsValid(t *testing.T) {
	t.Run("valid statuses should return true", func(t *testing.T) {
		assert := test.NewAssert(t)
		validStatuses := []domain.Status{
			domain.StatusPending,
			domain.StatusInProgress,
			domain.StatusCompleted,
		}

		for _, status := range validStatuses {
			result := status.IsValid()

			assert.True(result)
		}
	})

	t.Run("invalid status should return false", func(t *testing.T) {
		assert := test.NewAssert(t)
		invalidStatus := domain.Status(test.InvalidStatus)

		result := invalidStatus.IsValid()

		assert.False(result)
	})

	t.Run("empty status should return false", func(t *testing.T) {
		assert := test.NewAssert(t)
		emptyStatus := domain.Status("")

		result := emptyStatus.IsValid()

		assert.False(result)
	})
}

func TestPriority_IsValid(t *testing.T) {
	t.Run("valid priorities should return true", func(t *testing.T) {
		assert := test.NewAssert(t)
		validPriorities := []domain.Priority{
			domain.PriorityLow,
			domain.PriorityMedium,
			domain.PriorityHigh,
		}

		for _, priority := range validPriorities {
			result := priority.IsValid()

			assert.True(result)
		}
	})

	t.Run("invalid priority should return false", func(t *testing.T) {
		assert := test.NewAssert(t)
		invalidPriority := domain.Priority(test.InvalidPriority)

		result := invalidPriority.IsValid()

		assert.False(result)
	})

	t.Run("empty priority should return false", func(t *testing.T) {
		assert := test.NewAssert(t)
		emptyPriority := domain.Priority("")

		result := emptyPriority.IsValid()

		assert.False(result)
	})
}

func TestValidateUUID(t *testing.T) {
	t.Run("valid UUIDs should return no error", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		validUUIDs := []string{
			test.ValidUUID,
			test.ValidUUID2,
			"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			"AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE",
		}

		for _, uuid := range validUUIDs {
			err := domain.ValidateUUID(uuid)

			assert.NoError(err)
		}
	})

	t.Run("invalid UUID should return ErrInvalidID", func(t *testing.T) {
		assert := test.NewAssert(t)

		err := domain.ValidateUUID(test.InvalidUUID)

		assert.ErrorIs(err, domain.ErrInvalidID)
	})

	t.Run("empty string should return ErrInvalidID", func(t *testing.T) {
		assert := test.NewAssert(t)

		err := domain.ValidateUUID("")

		assert.ErrorIs(err, domain.ErrInvalidID)
	})

	t.Run("UUID without dashes should return ErrInvalidID", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		uuidNoDashes := "123e4567e89b12d3a456426614174000"

		err := domain.ValidateUUID(uuidNoDashes)

		assert.ErrorIs(err, domain.ErrInvalidID)
	})

	t.Run("partial UUID should return ErrInvalidID", func(t *testing.T) {
		// Arrange
		assert := test.NewAssert(t)
		partialUUID := "123e4567-e89b-12d3-a456"

		err := domain.ValidateUUID(partialUUID)

		assert.ErrorIs(err, domain.ErrInvalidID)
	})
}
