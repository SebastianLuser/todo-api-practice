package controller_test

import (
	"net/http"
	"testing"
	"time"

	"todo-api/pkg/controller"
	"todo-api/pkg/domain"
	"todo-api/test"
	"todo-api/web"
)

const (
	validUUID          = "123e4567-e89b-12d3-a456-426614174000"
	invalidUUID        = "not-a-valid-uuid"
	invalidStatus      = "invalid_status"
	invalidPriority    = "invalid_priority"
	maxTitleLength     = 100
	maxDescriptionLength = 500
)

var fixedTime = time.Date(2026, 1, 28, 10, 30, 0, 0, time.UTC)

func buildLongString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}

func titleTooLong() string {
	return buildLongString(maxTitleLength + 1)
}

func descriptionTooLong() string {
	return buildLongString(maxDescriptionLength + 1)
}

func buildValidTodo() domain.Todo {
	return domain.Todo{
		ID:          validUUID,
		Title:       "Test Todo",
		Description: "This is a test description",
		Status:      domain.StatusPending,
		Priority:    domain.PriorityMedium,
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}
}

func newTestController() *controller.Todo {
	errHandler := web.NewErrorHandler(
		web.NewErrorHandlerValueMapper(domain.ErrTodoNotFound, http.StatusNotFound),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidStatus, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidPriority, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidTitle, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidID, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrEmptyUpdateRequest, http.StatusBadRequest),
	)
	return controller.New(nil, errHandler)
}

func TestTodoController_Get_InvalidStatusFilter(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithQuery("status", invalidStatus)

	response := ctrl.Get(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Get_InvalidPriorityFilter(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithQuery("priority", invalidPriority)

	response := ctrl.Get(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_GetByID_MissingParam(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest()

	response := ctrl.GetByID(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_GetByID_InvalidUUID(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithParam("id", invalidUUID)

	response := ctrl.GetByID(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Create_InvalidJSON(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithBody("invalid json")

	response := ctrl.Create(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Create_EmptyTitle(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithBody(`{"title": ""}`)

	response := ctrl.Create(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Create_TitleTooLong(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithBody(`{"title": "` + titleTooLong() + `"}`)

	response := ctrl.Create(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Create_InvalidStatus(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithBody(`{"title": "Test", "status": "` + invalidStatus + `"}`)

	response := ctrl.Create(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Create_InvalidPriority(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithBody(`{"title": "Test", "priority": "` + invalidPriority + `"}`)

	response := ctrl.Create(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Create_DescriptionTooLong(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithBody(`{"title": "Test", "description": "` + descriptionTooLong() + `"}`)

	response := ctrl.Create(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Update_MissingParam(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest()

	response := ctrl.Update(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Update_InvalidUUID(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithParam("id", invalidUUID)

	response := ctrl.Update(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Update_EmptyBody(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().
		WithParam("id", validUUID).
		WithBody(`{}`)

	response := ctrl.Update(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Update_EmptyTitle(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().
		WithParam("id", validUUID).
		WithBody(`{"title": ""}`)

	response := ctrl.Update(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Update_InvalidStatus(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().
		WithParam("id", validUUID).
		WithBody(`{"status": "` + invalidStatus + `"}`)

	response := ctrl.Update(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Delete_MissingParam(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest()

	response := ctrl.Delete(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestTodoController_Delete_InvalidUUID(t *testing.T) {
	ctrl := newTestController()
	req := test.NewMockRequest().WithParam("id", invalidUUID)

	response := ctrl.Delete(req)

	if response.Status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, response.Status)
	}
}

func TestMapTodoToResponse(t *testing.T) {
	todo := buildValidTodo()

	response := controller.MapTodoToResponse(todo)

	if response.ID != todo.ID {
		t.Errorf("expected ID %s, got %s", todo.ID, response.ID)
	}
	if response.Title != todo.Title {
		t.Errorf("expected title %s, got %s", todo.Title, response.Title)
	}
	if response.Description != todo.Description {
		t.Errorf("expected description %s, got %s", todo.Description, response.Description)
	}
	if response.Status != string(todo.Status) {
		t.Errorf("expected status %s, got %s", string(todo.Status), response.Status)
	}
	if response.Priority != string(todo.Priority) {
		t.Errorf("expected priority %s, got %s", string(todo.Priority), response.Priority)
	}
	expectedTimeStr := "2026-01-28T10:30:00Z"
	if response.CreatedAt != expectedTimeStr {
		t.Errorf("expected createdAt %s, got %s", expectedTimeStr, response.CreatedAt)
	}
	if response.UpdatedAt != expectedTimeStr {
		t.Errorf("expected updatedAt %s, got %s", expectedTimeStr, response.UpdatedAt)
	}
}

func TestMapTodosToResponse_Multiple(t *testing.T) {
	todo1 := buildValidTodo()
	todo1.ID = "1"
	todo2 := buildValidTodo()
	todo2.ID = "2"
	todos := []domain.Todo{todo1, todo2}

	responses := controller.MapTodosToResponse(todos)

	if len(responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(responses))
	}
	if responses[0].ID != "1" {
		t.Errorf("expected first ID 1, got %s", responses[0].ID)
	}
	if responses[1].ID != "2" {
		t.Errorf("expected second ID 2, got %s", responses[1].ID)
	}
}

func TestMapTodosToResponse_Empty(t *testing.T) {
	todos := []domain.Todo{}

	responses := controller.MapTodosToResponse(todos)

	if len(responses) != 0 {
		t.Errorf("expected 0 responses, got %d", len(responses))
	}
}
