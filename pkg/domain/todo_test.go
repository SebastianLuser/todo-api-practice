package domain_test

import (
	"errors"
	"testing"

	"todo-api/pkg/domain"
)

const (
	validUUID      = "123e4567-e89b-12d3-a456-426614174000"
	validUUID2     = "550e8400-e29b-41d4-a716-446655440000"
	invalidUUID    = "not-a-valid-uuid"
	invalidStatus  = "invalid_status"
	invalidPriority = "invalid_priority"
)

func TestStatus_IsValid_Pending(t *testing.T) {
	if !domain.StatusPending.IsValid() {
		t.Error("expected StatusPending to be valid")
	}
}

func TestStatus_IsValid_InProgress(t *testing.T) {
	if !domain.StatusInProgress.IsValid() {
		t.Error("expected StatusInProgress to be valid")
	}
}

func TestStatus_IsValid_Completed(t *testing.T) {
	if !domain.StatusCompleted.IsValid() {
		t.Error("expected StatusCompleted to be valid")
	}
}

func TestStatus_IsValid_Invalid(t *testing.T) {
	status := domain.Status(invalidStatus)
	if status.IsValid() {
		t.Error("expected invalid status to return false")
	}
}

func TestStatus_IsValid_Empty(t *testing.T) {
	emptyStatus := domain.Status("")
	if emptyStatus.IsValid() {
		t.Error("expected empty status to return false")
	}
}

func TestPriority_IsValid_Low(t *testing.T) {
	if !domain.PriorityLow.IsValid() {
		t.Error("expected PriorityLow to be valid")
	}
}

func TestPriority_IsValid_Medium(t *testing.T) {
	if !domain.PriorityMedium.IsValid() {
		t.Error("expected PriorityMedium to be valid")
	}
}

func TestPriority_IsValid_High(t *testing.T) {
	if !domain.PriorityHigh.IsValid() {
		t.Error("expected PriorityHigh to be valid")
	}
}

func TestPriority_IsValid_Invalid(t *testing.T) {
	priority := domain.Priority(invalidPriority)
	if priority.IsValid() {
		t.Error("expected invalid priority to return false")
	}
}

func TestPriority_IsValid_Empty(t *testing.T) {
	emptyPriority := domain.Priority("")
	if emptyPriority.IsValid() {
		t.Error("expected empty priority to return false")
	}
}

func TestValidateUUID_Valid(t *testing.T) {
	err := domain.ValidateUUID(validUUID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateUUID_Valid2(t *testing.T) {
	err := domain.ValidateUUID(validUUID2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateUUID_ValidFormat(t *testing.T) {
	err := domain.ValidateUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateUUID_ValidUppercase(t *testing.T) {
	err := domain.ValidateUUID("AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateUUID_Invalid(t *testing.T) {
	err := domain.ValidateUUID(invalidUUID)
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Errorf("expected ErrInvalidID, got %v", err)
	}
}

func TestValidateUUID_Empty(t *testing.T) {
	err := domain.ValidateUUID("")
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Errorf("expected ErrInvalidID, got %v", err)
	}
}

func TestValidateUUID_NoDashes(t *testing.T) {
	uuidNoDashes := "123e4567e89b12d3a456426614174000"
	err := domain.ValidateUUID(uuidNoDashes)
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Errorf("expected ErrInvalidID, got %v", err)
	}
}

func TestValidateUUID_Partial(t *testing.T) {
	partialUUID := "123e4567-e89b-12d3-a456"
	err := domain.ValidateUUID(partialUUID)
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Errorf("expected ErrInvalidID, got %v", err)
	}
}
