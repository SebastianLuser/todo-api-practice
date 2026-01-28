package main

import (
	"todo-api/boot"
	"todo-api/pkg/controller"
	webgin "todo-api/web/gin"
)

func registerTodoRoutes(router boot.GinRouter, ctrl *controller.Todo) {
	router.GET("/api/todos", webgin.NewHandlerJSON(ctrl.Get))
}
