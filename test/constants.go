package test

import (
	"time"

	"todo-api/pkg/domain"
)

// Test UUIDs
const (
	ValidUUID     = "123e4567-e89b-12d3-a456-426614174000"
	ValidUUID2    = "550e8400-e29b-41d4-a716-446655440000"
	InvalidUUID   = "not-a-valid-uuid"
	NonExistentID = "00000000-0000-0000-0000-000000000000"
)

// Test Titles
const (
	ValidTitle   = "Test Todo"
	ValidTitle2  = "Another Todo"
	EmptyTitle   = ""
	UpdatedTitle = "Updated Title"
)

// Test Descriptions
const (
	ValidDescription = "This is a test description"
	EmptyDescription = ""
)

// MaxLengths for validation testing
const (
	MaxTitleLength       = 100
	MaxDescriptionLength = 500
)

// Test Status and Priority values
var (
	StatusPending    = domain.StatusPending
	StatusInProgress = domain.StatusInProgress
	StatusCompleted  = domain.StatusCompleted

	PriorityLow    = domain.PriorityLow
	PriorityMedium = domain.PriorityMedium
	PriorityHigh   = domain.PriorityHigh
)

// Invalid values
const (
	InvalidStatus   = "invalid_status"
	InvalidPriority = "invalid_priority"
)

// Test timestamps
var (
	FixedTime    = time.Date(2026, 1, 28, 10, 30, 0, 0, time.UTC)
	FixedTimeStr = "2026-01-28T10:30:00Z"
)

// BuildValidTodo creates a valid Todo for testing
func BuildValidTodo() domain.Todo {
	return domain.Todo{
		ID:          ValidUUID,
		Title:       ValidTitle,
		Description: ValidDescription,
		Status:      domain.StatusPending,
		Priority:    domain.PriorityMedium,
		CreatedAt:   FixedTime,
		UpdatedAt:   FixedTime,
	}
}

// BuildValidTodoWithID creates a valid Todo with custom ID
func BuildValidTodoWithID(id string) domain.Todo {
	todo := BuildValidTodo()
	todo.ID = id
	return todo
}

// BuildLongString creates a string of specified length
func BuildLongString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}

// TitleTooLong returns a title that exceeds max length
func TitleTooLong() string {
	return BuildLongString(MaxTitleLength + 1)
}

// DescriptionTooLong returns a description that exceeds max length
func DescriptionTooLong() string {
	return BuildLongString(MaxDescriptionLength + 1)
}
