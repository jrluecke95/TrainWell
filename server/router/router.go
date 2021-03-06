package router

import (
	"server/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	//coaches
	router.HandleFunc("/api/coach", middleware.CreateCoach).Methods("POST")
	router.HandleFunc("/api/coaches", middleware.GetCoaches).Methods("GET")

	//clients
	router.HandleFunc("/api/client", middleware.CreateClient).Methods("POST")
	router.HandleFunc("/api/clients", middleware.GetClients).Methods("GET")

	//assign and unassign coach
	router.HandleFunc("/api/assignCoach", middleware.AssignCoach).Methods("PUT")
	router.HandleFunc("/api/unassignCoach", middleware.UnassignCoach).Methods("PUT")

	//exercises
	router.HandleFunc("/api/exercise", middleware.CreateExercise).Methods("POST")
	router.HandleFunc("/api/exercises", middleware.GetExercises).Methods("GET")

	//exercise details
	router.HandleFunc("/api/exercise/details", middleware.CreateExerciseDetails).Methods("POST")

	return router
}
