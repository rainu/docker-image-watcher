package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rainu/docker-image-watcher/internal/database"
	"github.com/rainu/docker-image-watcher/internal/server/handler"
	"net/http"
	"os"
)

func NewRouter(repo database.Repository) http.Handler {
	router := mux.NewRouter()

	// RESTful API
	router.Handle("/api/v1/register", handler.NewAddObservationHandler(repo)).Methods(http.MethodPost)

	return handlers.LoggingHandler(os.Stdout, router)
}
