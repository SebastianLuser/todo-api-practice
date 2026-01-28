package main

import (
	"todo-api/pkg/controller"
	"todo-api/pkg/usecase"
	"todo-api/web"
)

func NewTodoController(uc usecase.Usecase, errHandler web.ErrorHandler) controller.Controller {
	return controller.New(uc, errHandler)
}
