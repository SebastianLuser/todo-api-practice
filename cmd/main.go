package main

import (
	"context"

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
	router.GET("/health", webgin.NewHandlerJSON(func(req web.Request) web.Response {
		return web.NewJSONResponse(200, map[string]string{"status": "healthy"})
	}))
}
