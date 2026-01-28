package main

import (
	"todo-api/pkg/service"
)

func NewTodoService() service.Service {
	return service.New()
}
