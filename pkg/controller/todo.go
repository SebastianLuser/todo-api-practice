package controller

import (
	"net/http"

	"todo-api/pkg/domain"
	"todo-api/pkg/usecase"
	"todo-api/web"
)

type (
	Todo struct {
		usecase    *usecase.Todo
		errHandler web.ErrorHandler
	}

	TodoResponse struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
		Status      string `json:"status"`
		Priority    string `json:"priority"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}

	GetResponse struct {
		Data  []TodoResponse `json:"data"`
		Total int            `json:"total"`
	}
)

func New(uc *usecase.Todo, errHandler web.ErrorHandler) *Todo {
	return &Todo{
		usecase:    uc,
		errHandler: errHandler,
	}
}

func (c *Todo) Get(req web.Request) web.Response {
	input := usecase.ListInput{}

	if statusStr, ok := req.Query("status"); ok {
		status := domain.Status(statusStr)
		if !status.IsValid() {
			return web.NewJSONResponseFromError(
				web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidStatus),
			)
		}
		input.Status = &status
	}

	if priorityStr, ok := req.Query("priority"); ok {
		priority := domain.Priority(priorityStr)
		if !priority.IsValid() {
			return web.NewJSONResponseFromError(
				web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidPriority),
			)
		}
		input.Priority = &priority
	}

	output, err := c.usecase.Get(req.Context(), input)
	if err != nil {
		return web.NewJSONResponseFromError(c.errHandler.Handle(err))
	}

	response := GetResponse{
		Data:  MapTodosToResponse(output.Todos),
		Total: output.Total,
	}

	return web.NewJSONResponse(http.StatusOK, response)
}

func MapTodoToResponse(todo domain.Todo) TodoResponse {
	return TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      string(todo.Status),
		Priority:    string(todo.Priority),
		CreatedAt:   todo.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   todo.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func MapTodosToResponse(todos []domain.Todo) []TodoResponse {
	result := make([]TodoResponse, len(todos))
	for i, todo := range todos {
		result[i] = MapTodoToResponse(todo)
	}
	return result
}
