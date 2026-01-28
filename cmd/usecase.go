package main

import (
	"todo-api/database"
	"todo-api/pkg/usecase"
)

func NewTodoUsecase() *usecase.Todo {
	db := database.NewDatabase()
	svc := NewTodoService(db)
	return usecase.New(svc)
}
