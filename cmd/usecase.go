package main

import (
	"todo-api/pkg/service"
	"todo-api/pkg/usecase"
)

func NewTodoUsecase(svc service.Service) usecase.Usecase {
	return usecase.New(svc)
}
