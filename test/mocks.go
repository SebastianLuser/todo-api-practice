package test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"todo-api/pkg/domain"
	"todo-api/pkg/service"
	"todo-api/web"
)

type MockTodoService struct {
	GetFn      func(ctx context.Context, filters service.Filters) ([]domain.Todo, error)
	GetByIDFn  func(ctx context.Context, id string) (domain.Todo, error)
	CreateFn   func(ctx context.Context, input service.CreateInput) (domain.Todo, error)
	UpdateFn   func(ctx context.Context, id string, input service.UpdateInput) (domain.Todo, error)
	DeleteFn   func(ctx context.Context, id string) error
}

func (m *MockTodoService) Get(ctx context.Context, filters service.Filters) ([]domain.Todo, error) {
	return m.GetFn(ctx, filters)
}

func (m *MockTodoService) GetByID(ctx context.Context, id string) (domain.Todo, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockTodoService) Create(ctx context.Context, input service.CreateInput) (domain.Todo, error) {
	return m.CreateFn(ctx, input)
}

func (m *MockTodoService) Update(ctx context.Context, id string, input service.UpdateInput) (domain.Todo, error) {
	return m.UpdateFn(ctx, id, input)
}

func (m *MockTodoService) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

type MockRequest struct {
	Ctx        context.Context
	ParamsMap  map[string]string
	QueriesMap map[string]string
	BodyStr    string
}

func NewMockRequest() *MockRequest {
	return &MockRequest{
		Ctx:        context.Background(),
		ParamsMap:  make(map[string]string),
		QueriesMap: make(map[string]string),
	}
}

func (m *MockRequest) WithParam(key, value string) *MockRequest {
	m.ParamsMap[key] = value
	return m
}

func (m *MockRequest) WithQuery(key, value string) *MockRequest {
	m.QueriesMap[key] = value
	return m
}

func (m *MockRequest) WithBody(body string) *MockRequest {
	m.BodyStr = body
	return m
}

func (m *MockRequest) Context() context.Context                           { return m.Ctx }
func (m *MockRequest) Raw() *http.Request                                 { return &http.Request{} }
func (m *MockRequest) DeclaredPath() string                               { return "" }
func (m *MockRequest) Params() []web.Param                                { return nil }
func (m *MockRequest) Queries() url.Values                                { return nil }
func (m *MockRequest) Headers() http.Header                               { return nil }
func (m *MockRequest) Body() io.ReadCloser                                { return io.NopCloser(bytes.NewBufferString(m.BodyStr)) }
func (m *MockRequest) Header(key string) ([]string, bool)                 { return nil, false }
func (m *MockRequest) FormFile(key string) (*multipart.FileHeader, error) { return nil, nil }
func (m *MockRequest) FormValue(key string) (string, bool)                { return "", false }
func (m *MockRequest) MultipartForm() (*multipart.Form, error)            { return nil, nil }

func (m *MockRequest) Param(key string) (string, bool) {
	v, ok := m.ParamsMap[key]
	return v, ok
}

func (m *MockRequest) Query(key string) (string, bool) {
	v, ok := m.QueriesMap[key]
	return v, ok
}
