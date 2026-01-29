package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"todo-api/pkg/domain"
)

func TestErrTodoNotFound_ErrorsIs(t *testing.T) {
	if !errors.Is(domain.ErrTodoNotFound, domain.ErrTodoNotFound) {
		t.Error("expected ErrTodoNotFound to match itself")
	}
}

func TestErrInvalidStatus_ErrorsIs(t *testing.T) {
	if !errors.Is(domain.ErrInvalidStatus, domain.ErrInvalidStatus) {
		t.Error("expected ErrInvalidStatus to match itself")
	}
}

func TestErrInvalidPriority_ErrorsIs(t *testing.T) {
	if !errors.Is(domain.ErrInvalidPriority, domain.ErrInvalidPriority) {
		t.Error("expected ErrInvalidPriority to match itself")
	}
}

func TestErrInvalidTitle_ErrorsIs(t *testing.T) {
	if !errors.Is(domain.ErrInvalidTitle, domain.ErrInvalidTitle) {
		t.Error("expected ErrInvalidTitle to match itself")
	}
}

func TestErrInvalidDescription_ErrorsIs(t *testing.T) {
	if !errors.Is(domain.ErrInvalidDescription, domain.ErrInvalidDescription) {
		t.Error("expected ErrInvalidDescription to match itself")
	}
}

func TestErrInvalidID_ErrorsIs(t *testing.T) {
	if !errors.Is(domain.ErrInvalidID, domain.ErrInvalidID) {
		t.Error("expected ErrInvalidID to match itself")
	}
}

func TestErrEmptyUpdateRequest_ErrorsIs(t *testing.T) {
	if !errors.Is(domain.ErrEmptyUpdateRequest, domain.ErrEmptyUpdateRequest) {
		t.Error("expected ErrEmptyUpdateRequest to match itself")
	}
}

func TestWrappedErrTodoNotFound_IsIdentifiable(t *testing.T) {
	wrappedErr := fmt.Errorf("service error: %w", domain.ErrTodoNotFound)

	if !errors.Is(wrappedErr, domain.ErrTodoNotFound) {
		t.Error("expected wrapped ErrTodoNotFound to be identifiable")
	}
}

func TestWrappedErrInvalidID_IsIdentifiable(t *testing.T) {
	wrappedErr := fmt.Errorf("validation failed: %w", domain.ErrInvalidID)

	if !errors.Is(wrappedErr, domain.ErrInvalidID) {
		t.Error("expected wrapped ErrInvalidID to be identifiable")
	}
}

func TestDoubleWrappedError_IsIdentifiable(t *testing.T) {
	firstWrap := fmt.Errorf("layer 1: %w", domain.ErrInvalidStatus)
	secondWrap := fmt.Errorf("layer 2: %w", firstWrap)

	if !errors.Is(secondWrap, domain.ErrInvalidStatus) {
		t.Error("expected double wrapped error to be identifiable")
	}
}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func TestWrappedCustomError_ExtractableWithErrorAs(t *testing.T) {
	originalErr := &customError{msg: "custom error"}
	wrappedErr := fmt.Errorf("wrapped: %w", originalErr)
	var target *customError

	result := errors.As(wrappedErr, &target)

	if !result {
		t.Error("expected errors.As to return true")
	}
	if target.msg != "custom error" {
		t.Errorf("expected message %s, got %s", "custom error", target.msg)
	}
}

func TestErrTodoNotFound_Message(t *testing.T) {
	expected := "todo not found"
	if domain.ErrTodoNotFound.Error() != expected {
		t.Errorf("expected message %s, got %s", expected, domain.ErrTodoNotFound.Error())
	}
}

func TestErrInvalidStatus_Message(t *testing.T) {
	expected := "invalid status: must be pending, in_progress, or completed"
	if domain.ErrInvalidStatus.Error() != expected {
		t.Errorf("expected message %s, got %s", expected, domain.ErrInvalidStatus.Error())
	}
}

func TestErrInvalidPriority_Message(t *testing.T) {
	expected := "invalid priority: must be low, medium, or high"
	if domain.ErrInvalidPriority.Error() != expected {
		t.Errorf("expected message %s, got %s", expected, domain.ErrInvalidPriority.Error())
	}
}

func TestErrInvalidTitle_Message(t *testing.T) {
	expected := "invalid title: must be between 1 and 100 characters"
	if domain.ErrInvalidTitle.Error() != expected {
		t.Errorf("expected message %s, got %s", expected, domain.ErrInvalidTitle.Error())
	}
}

func TestErrInvalidDescription_Message(t *testing.T) {
	expected := "invalid description: must be at most 500 characters"
	if domain.ErrInvalidDescription.Error() != expected {
		t.Errorf("expected message %s, got %s", expected, domain.ErrInvalidDescription.Error())
	}
}

func TestErrInvalidID_Message(t *testing.T) {
	expected := "invalid id: must be a valid UUID"
	if domain.ErrInvalidID.Error() != expected {
		t.Errorf("expected message %s, got %s", expected, domain.ErrInvalidID.Error())
	}
}

func TestErrEmptyUpdateRequest_Message(t *testing.T) {
	expected := "update request must contain at least one field"
	if domain.ErrEmptyUpdateRequest.Error() != expected {
		t.Errorf("expected message %s, got %s", expected, domain.ErrEmptyUpdateRequest.Error())
	}
}

func TestDifferentErrors_ShouldNotMatch(t *testing.T) {
	if errors.Is(domain.ErrTodoNotFound, domain.ErrInvalidID) {
		t.Error("expected ErrTodoNotFound to not match ErrInvalidID")
	}
}

func TestErrInvalidStatus_ShouldNotMatchErrInvalidPriority(t *testing.T) {
	if errors.Is(domain.ErrInvalidStatus, domain.ErrInvalidPriority) {
		t.Error("expected ErrInvalidStatus to not match ErrInvalidPriority")
	}
}
