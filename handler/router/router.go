package router

import (
	"database/sql"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
	"net/http"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	helthzHandler := handler.NewHealthzHandler()
	TODOService := service.NewTODOService(todoDB)
	TODOHandler := handler.NewTODOHandler(TODOService)
	mux.Handle("/healthz", helthzHandler)
	mux.Handle("/todos", TODOHandler)

	return mux
}
