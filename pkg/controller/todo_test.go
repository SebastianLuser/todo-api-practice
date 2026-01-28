package controller_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"

	"todo-api/pkg/controller"
	"todo-api/pkg/domain"
	"todo-api/test"
	"todo-api/web"
)

// mockRequest implements web.Request for testing
type mockRequest struct {
	ctx     context.Context
	params  map[string]string
	queries map[string]string
	body    string
}

func newMockRequest() *mockRequest {
	return &mockRequest{
		ctx:     context.Background(),
		params:  make(map[string]string),
		queries: make(map[string]string),
	}
}

func (m *mockRequest) withParam(key, value string) *mockRequest {
	m.params[key] = value
	return m
}

func (m *mockRequest) withQuery(key, value string) *mockRequest {
	m.queries[key] = value
	return m
}

func (m *mockRequest) withBody(body string) *mockRequest {
	m.body = body
	return m
}

func (m *mockRequest) Context() context.Context                           { return m.ctx }
func (m *mockRequest) Raw() *http.Request                                 { return &http.Request{} }
func (m *mockRequest) DeclaredPath() string                               { return "" }
func (m *mockRequest) Params() []web.Param                                { return nil }
func (m *mockRequest) Queries() url.Values                                { return nil }
func (m *mockRequest) Headers() http.Header                               { return nil }
func (m *mockRequest) Body() io.ReadCloser                                { return io.NopCloser(bytes.NewBufferString(m.body)) }
func (m *mockRequest) Header(key string) ([]string, bool)                 { return nil, false }
func (m *mockRequest) FormFile(key string) (*multipart.FileHeader, error) { return nil, nil }
func (m *mockRequest) FormValue(key string) (string, bool)                { return "", false }
func (m *mockRequest) MultipartForm() (*multipart.Form, error)            { return nil, nil }

func (m *mockRequest) Param(key string) (string, bool) {
	v, ok := m.params[key]
	return v, ok
}

func (m *mockRequest) Query(key string) (string, bool) {
	v, ok := m.queries[key]
	return v, ok
}

// newTestController creates a controller with error handler for testing
func newTestController() *controller.Todo {
	errHandler := web.NewErrorHandler(
		web.NewErrorHandlerValueMapper(domain.ErrTodoNotFound, http.StatusNotFound),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidStatus, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidPriority, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidTitle, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidID, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrEmptyUpdateRequest, http.StatusBadRequest),
	)
	// Note: We pass nil usecase because we're testing validation logic only
	// For full integration tests, you would inject a mock usecase
	return controller.New(nil, errHandler)
}

func TestTodoController_Get(t *testing.T) {
	t.Run("should return 400 for invalid status filter", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withQuery("status", test.InvalidStatus)

		response := ctrl.Get(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for invalid priority filter", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withQuery("priority", test.InvalidPriority)

		response := ctrl.Get(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})
}

func TestTodoController_GetByID(t *testing.T) {
	t.Run("should return 400 when id param is missing", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest()

		response := ctrl.GetByID(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withParam("id", test.InvalidUUID)

		response := ctrl.GetByID(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})
}

func TestTodoController_Create(t *testing.T) {
	t.Run("should return 400 for invalid JSON body", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withBody("invalid json")

		response := ctrl.Create(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for empty title", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withBody(`{"title": ""}`)

		response := ctrl.Create(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for title too long", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withBody(`{"title": "` + test.TitleTooLong() + `"}`)

		response := ctrl.Create(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for invalid status", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withBody(`{"title": "Test", "status": "` + test.InvalidStatus + `"}`)

		response := ctrl.Create(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for invalid priority", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withBody(`{"title": "Test", "priority": "` + test.InvalidPriority + `"}`)

		response := ctrl.Create(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for description too long", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withBody(`{"title": "Test", "description": "` + test.DescriptionTooLong() + `"}`)

		response := ctrl.Create(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})
}

func TestTodoController_Update(t *testing.T) {
	t.Run("should return 400 when id param is missing", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest()

		response := ctrl.Update(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withParam("id", test.InvalidUUID)

		response := ctrl.Update(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for empty update body", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().
			withParam("id", test.ValidUUID).
			withBody(`{}`)

		response := ctrl.Update(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for empty title in update", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().
			withParam("id", test.ValidUUID).
			withBody(`{"title": ""}`)

		response := ctrl.Update(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for invalid status in update", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().
			withParam("id", test.ValidUUID).
			withBody(`{"status": "` + test.InvalidStatus + `"}`)

		response := ctrl.Update(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})
}

func TestTodoController_Delete(t *testing.T) {
	t.Run("should return 400 when id param is missing", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest()

		response := ctrl.Delete(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		assert := test.NewAssert(t)
		ctrl := newTestController()
		req := newMockRequest().withParam("id", test.InvalidUUID)

		response := ctrl.Delete(req)

		assert.StatusCode(http.StatusBadRequest, response.Status)
	})
}

func TestMapTodoToResponse(t *testing.T) {
	t.Run("should map todo to response correctly", func(t *testing.T) {
		assert := test.NewAssert(t)
		todo := test.BuildValidTodo()

		response := controller.MapTodoToResponse(todo)

		assert.Equal(todo.ID, response.ID)
		assert.Equal(todo.Title, response.Title)
		assert.Equal(todo.Description, response.Description)
		assert.Equal(string(todo.Status), response.Status)
		assert.Equal(string(todo.Priority), response.Priority)
		assert.Equal(test.FixedTimeStr, response.CreatedAt)
		assert.Equal(test.FixedTimeStr, response.UpdatedAt)
	})
}

func TestMapTodosToResponse(t *testing.T) {
	t.Run("should map multiple todos to responses", func(t *testing.T) {
		assert := test.NewAssert(t)
		todo1 := test.BuildValidTodoWithID("1")
		todo2 := test.BuildValidTodoWithID("2")
		todos := []domain.Todo{todo1, todo2}

		responses := controller.MapTodosToResponse(todos)

		assert.Equal(2, len(responses))
		assert.Equal("1", responses[0].ID)
		assert.Equal("2", responses[1].ID)
	})

	t.Run("should return empty slice for empty input", func(t *testing.T) {
		assert := test.NewAssert(t)
		todos := []domain.Todo{}

		responses := controller.MapTodosToResponse(todos)

		assert.Equal(0, len(responses))
	})
}
