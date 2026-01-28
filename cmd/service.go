package main

import (
	"database/sql"

	"todo-api/pkg/service"
)

func NewTodoService(db *sql.DB) service.Todo {
	return service.New(db)
}
