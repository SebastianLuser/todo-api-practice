package main

import (
	"todo-api/boot"
	"todo-api/pkg/controller"
	webgin "todo-api/web/gin"
)

func registerTodoRoutes(router boot.GinRouter, ctrl *controller.Todo) {
	router.GET("/api/todos", webgin.NewHandlerJSON(ctrl.Get))
	router.GET("/api/todos/:id", webgin.NewHandlerJSON(ctrl.GetByID))
	router.POST("/api/todos", webgin.NewHandlerJSON(ctrl.Create))
	router.PATCH("/api/todos/:id", webgin.NewHandlerJSON(ctrl.Update))
	router.DELETE("/api/todos/:id", webgin.NewHandlerJSON(ctrl.Delete))
}
