package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

type middlewareType func(http.Handler) http.Handler

func applyMiddleware(h http.Handler, middlewares ...middlewareType) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	todoService := service.NewTODOService(todoDB)

	helthzHandler := handler.NewHealthzHandler()
	todoHandler := handler.NewTODOHandler(todoService)
	doPanicHandler := handler.NewDoPanicHandler()

	middlewares := []middlewareType{
		middleware.OSParser,
		middleware.AccessLog,
		middleware.Recovery,
	}

	todosMiddlewares := append(middlewares, middleware.BasicAuth)

	mux.Handle("/healthz", applyMiddleware(helthzHandler, middlewares...))
	mux.Handle("/todos", applyMiddleware(todoHandler, todosMiddlewares...))
	mux.Handle("/do_panic", applyMiddleware(doPanicHandler, middlewares...))

	return mux
}
