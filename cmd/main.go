package main

import (
	"context"
	"net/http"
	"todo-api/pkg/domain"

	"todo-api/boot"
	"todo-api/web"
	webgin "todo-api/web/gin"
)

func main() {
	boot.NewGin(
		boot.DefaultGinMiddlewareMapper(),
		routesMapper,
	).MustRun()
}

func routesMapper(ctx context.Context, conf boot.Config, router boot.GinRouter) {

	todoService := NewTodoService()

	todoUsecase := NewTodoUsecase(todoService)

	errHandler := web.NewErrorHandler(
		web.NewErrorHandlerValueMapper(domain.ErrTodoNotFound, http.StatusNotFound),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidStatus, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidPriority, http.StatusBadRequest),
		web.NewErrorHandlerValueMapper(domain.ErrInvalidTitle, http.StatusBadRequest),
	)

	todoController := NewTodoController(todoUsecase, errHandler)

	router.GET("/health", webgin.NewHandlerJSON(func(req web.Request) web.Response {
		return web.NewJSONResponse(http.StatusOK, map[string]string{"status": "healthy"})
	}))

	registerTodoRoutes(router, todoController)
}
