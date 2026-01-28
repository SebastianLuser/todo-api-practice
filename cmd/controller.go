package main

import (
	"net/http"

	"todo-api/pkg/controller"
	"todo-api/pkg/domain"
	"todo-api/web"
)

func NewTodoController() *controller.Todo {
	uc := NewTodoUsecase()
	return controller.New(uc, newErrorHandler())
}

func newErrorHandler() web.ErrorHandler {
	return web.NewErrorHandler(
		web.NewErrorHandlerValueMapper(domain.ErrTodoNotFound, http.StatusNotFound),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidStatus, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidPriority, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidTitle, http.StatusBadRequest),
	)
}
