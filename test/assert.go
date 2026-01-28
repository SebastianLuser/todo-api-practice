package test

import (
	"errors"
	"testing"
)

// Assert provides test assertion helpers
type Assert struct {
	t *testing.T
}

// NewAssert creates a new Assert instance
func NewAssert(t *testing.T) *Assert {
	t.Helper()
	return &Assert{t: t}
}

// Equal asserts that two values are equal
func (a *Assert) Equal(expected, actual any) {
	a.t.Helper()
	if expected != actual {
		a.t.Errorf("expected %v, got %v", expected, actual)
	}
}

// NotEqual asserts that two values are not equal
func (a *Assert) NotEqual(expected, actual any) {
	a.t.Helper()
	if expected == actual {
		a.t.Errorf("expected values to be different, got %v", actual)
	}
}

// Nil asserts that a value is nil
func (a *Assert) Nil(actual any) {
	a.t.Helper()
	if actual != nil {
		a.t.Errorf("expected nil, got %v", actual)
	}
}

// NotNil asserts that a value is not nil
func (a *Assert) NotNil(actual any) {
	a.t.Helper()
	if actual == nil {
		a.t.Error("expected not nil, got nil")
	}
}

// NoError asserts that error is nil
func (a *Assert) NoError(err error) {
	a.t.Helper()
	if err != nil {
		a.t.Errorf("expected no error, got %v", err)
	}
}

// Error asserts that error is not nil
func (a *Assert) Error(err error) {
	a.t.Helper()
	if err == nil {
		a.t.Error("expected error, got nil")
	}
}

// ErrorIs asserts that error matches target using errors.Is
func (a *Assert) ErrorIs(err, target error) {
	a.t.Helper()
	if !errors.Is(err, target) {
		a.t.Errorf("expected error %v, got %v", target, err)
	}
}

// ErrorAs asserts that error can be assigned to target using errors.As
func (a *Assert) ErrorAs(err error, target any) {
	a.t.Helper()
	if !errors.As(err, target) {
		a.t.Errorf("expected error to be assignable to %T, got %v", target, err)
	}
}

// True asserts that value is true
func (a *Assert) True(actual bool) {
	a.t.Helper()
	if !actual {
		a.t.Error("expected true, got false")
	}
}

// False asserts that value is false
func (a *Assert) False(actual bool) {
	a.t.Helper()
	if actual {
		a.t.Error("expected false, got true")
	}
}

// Len asserts that a slice has expected length
func (a *Assert) Len(actual any, expected int) {
	a.t.Helper()
	switch v := actual.(type) {
	case string:
		if len(v) != expected {
			a.t.Errorf("expected length %d, got %d", expected, len(v))
		}
	case []any:
		if len(v) != expected {
			a.t.Errorf("expected length %d, got %d", expected, len(v))
		}
	default:
		a.t.Errorf("unsupported type for Len assertion: %T", actual)
	}
}

// StatusCode asserts HTTP status code
func (a *Assert) StatusCode(expected, actual int) {
	a.t.Helper()
	if expected != actual {
		a.t.Errorf("expected status code %d, got %d", expected, actual)
	}
}
