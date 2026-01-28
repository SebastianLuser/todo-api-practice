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

	GetByIDResponse struct {
		Data TodoResponse `json:"data"`
	}

	CreateRequest struct {
		Title       string  `json:"title"`
		Description *string `json:"description,omitempty"`
		Status      *string `json:"status,omitempty"`
		Priority    *string `json:"priority,omitempty"`
	}

	CreateResponse struct {
		Data TodoResponse `json:"data"`
	}

	UpdateRequest struct {
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		Status      *string `json:"status,omitempty"`
		Priority    *string `json:"priority,omitempty"`
	}

	UpdateResponse struct {
		Data TodoResponse `json:"data"`
	}
)

func New(uc *usecase.Todo, errHandler web.ErrorHandler) *Todo {
	return &Todo{
		usecase:    uc,
		errHandler: errHandler,
	}
}

func (c *Todo) Create(req web.Request) web.Response {
	var body CreateRequest
	if err := web.DecodeJSON(req.Body(), &body); err != nil {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, err),
		)
	}

	if len(body.Title) == 0 || len(body.Title) > 100 {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidTitle),
		)
	}

	if body.Description != nil && len(*body.Description) > 500 {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidDescription),
		)
	}

	input := usecase.CreateInput{
		Title:       body.Title,
		Description: body.Description,
	}

	if body.Status != nil {
		status := domain.Status(*body.Status)
		if !status.IsValid() {
			return web.NewJSONResponseFromError(
				web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidStatus),
			)
		}
		input.Status = &status
	}

	if body.Priority != nil {
		priority := domain.Priority(*body.Priority)
		if !priority.IsValid() {
			return web.NewJSONResponseFromError(
				web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidPriority),
			)
		}
		input.Priority = &priority
	}

	output, err := c.usecase.Create(req.Context(), input)
	if err != nil {
		return web.NewJSONResponseFromError(c.errHandler.Handle(err))
	}

	response := CreateResponse{
		Data: MapTodoToResponse(output.Todo),
	}

	return web.NewJSONResponse(http.StatusCreated, response)
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

func (c *Todo) GetByID(req web.Request) web.Response {
	id, ok := req.Param("id")
	if !ok {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidID),
		)
	}

	if err := domain.ValidateUUID(id); err != nil {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, err),
		)
	}

	output, err := c.usecase.GetByID(req.Context(), id)
	if err != nil {
		return web.NewJSONResponseFromError(c.errHandler.Handle(err))
	}

	response := GetByIDResponse{
		Data: MapTodoToResponse(output.Todo),
	}

	return web.NewJSONResponse(http.StatusOK, response)
}

func (c *Todo) Update(req web.Request) web.Response {
	id, ok := req.Param("id")
	if !ok {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidID),
		)
	}

	if err := domain.ValidateUUID(id); err != nil {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, err),
		)
	}

	var body UpdateRequest
	if err := web.DecodeJSON(req.Body(), &body); err != nil {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, err),
		)
	}

	if body.Title == nil && body.Description == nil && body.Status == nil && body.Priority == nil {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, domain.ErrEmptyUpdateRequest),
		)
	}

	if body.Title != nil && (len(*body.Title) == 0 || len(*body.Title) > 100) {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidTitle),
		)
	}

	if body.Description != nil && len(*body.Description) > 500 {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidDescription),
		)
	}

	input := usecase.UpdateInput{
		Title:       body.Title,
		Description: body.Description,
	}

	if body.Status != nil {
		status := domain.Status(*body.Status)
		if !status.IsValid() {
			return web.NewJSONResponseFromError(
				web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidStatus),
			)
		}
		input.Status = &status
	}

	if body.Priority != nil {
		priority := domain.Priority(*body.Priority)
		if !priority.IsValid() {
			return web.NewJSONResponseFromError(
				web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidPriority),
			)
		}
		input.Priority = &priority
	}

	output, err := c.usecase.Update(req.Context(), id, input)
	if err != nil {
		return web.NewJSONResponseFromError(c.errHandler.Handle(err))
	}

	response := UpdateResponse{
		Data: MapTodoToResponse(output.Todo),
	}

	return web.NewJSONResponse(http.StatusOK, response)
}

func (c *Todo) Delete(req web.Request) web.Response {
	id, ok := req.Param("id")
	if !ok {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, domain.ErrInvalidID),
		)
	}

	if err := domain.ValidateUUID(id); err != nil {
		return web.NewJSONResponseFromError(
			web.NewResponseError(http.StatusBadRequest, err),
		)
	}

	err := c.usecase.Delete(req.Context(), id)
	if err != nil {
		return web.NewJSONResponseFromError(c.errHandler.Handle(err))
	}

	return web.NewJSONResponse(http.StatusNoContent, nil)
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
