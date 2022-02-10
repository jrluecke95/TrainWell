package router

import (
	"server/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/exercise", middleware.CreateExercise).Methods("POST")
	router.HandleFunc("/api/exercises", middleware.GetExercises).Methods(("GET"))

	return router
}
