package router

import (
	"server/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	//exercises
	router.HandleFunc("/api/exercise", middleware.CreateExercise).Methods("POST")
	router.HandleFunc("/api/exercises", middleware.GetExercises).Methods("GET")

	//coaches
	router.HandleFunc("/api/coach", middleware.CreateCoach).Methods("POST")
	router.HandleFunc("/api/coaches", middleware.GetCoaches).Methods("GET")

	//clients
	router.HandleFunc("/api/client", middleware.CreateClient).Methods("POST")
	router.HandleFunc("/api/clients", middleware.GetClients).Methods("GET")

	//assign coach
	router.HandleFunc("/api/assignCoach", middleware.AssignCoach).Methods("PUT")

	return router
}
