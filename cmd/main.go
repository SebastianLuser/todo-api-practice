package main

import (
	"context"

	"todo-api/boot"
)

func main() {
	boot.NewGin(
		boot.DefaultGinMiddlewareMapper(),
		setup,
	).MustRun()
}

func setup(_ context.Context, _ boot.Config, router boot.GinRouter) {
	registerTodoRoutes(router, NewTodoController())
}
