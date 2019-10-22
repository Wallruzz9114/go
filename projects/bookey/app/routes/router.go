package routes

import (
	"github.com/Wallruzz9114/bookey/app/handler"
	"github.com/Wallruzz9114/bookey/app/server"

	"github.com/go-chi/chi"
)

// New ...
func New(server *server.Server) *chi.Mux {
	logger := server.Logger()
	router := chi.NewRouter()

	router.Method("GET", "/", handler.NewHandler(server.HandleIndex, logger))

	return router
}
