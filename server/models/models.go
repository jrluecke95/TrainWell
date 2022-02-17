package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Client struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id"`
	Coaches      []primitive.ObjectID `json:"coaches"`
	WorkoutPlans []primitive.ObjectID `json:"workoutPlans"`
	FirstName    string               `json:"firstname"`
	LastName     string               `json:"lastname"`
	Email        string               `json:"email"`
	PhoneNumber  string               `json:"phonenumber"`
	Password     string               `json:"password"`
}
type Coach struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id"`
	Clients      []primitive.ObjectID `json:"clients"`
	Workouts     []primitive.ObjectID `json:"workouts"`
	WorkoutPlans []primitive.ObjectID `json:"workoutsplans"`
	FirstName    string               `json:"firstname"`
	LastName     string               `json:"lastname"`
	Email        string               `json:"email"`
	PhoneNumber  string               `json:"phonenumber"`
	Password     string               `json:"password"`
}

type Exercise struct {
	ID              primitive.ObjectID   `json:"_id" bson:"_id"`
	ExerciseDetails []primitive.ObjectID `json:"exercisedetails"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
}

type ExerciseDetails struct {
	ID       primitive.ObjectID   `json:"_id" bson:"_id"`
	Name_ID  primitive.ObjectID   `json:"nameid"`
	Workouts []primitive.ObjectID `json:"workouts"`
	Sets     int16                `json:"sets"`
	Reps     int16                `json:"reps"`
}

type Workout struct {
	ID               primitive.ObjectID   `json:"_id" bson:"_id"`
	ExercisesDetails []primitive.ObjectID `json:"exercisesdetails"`
	WorkoutPlans     []primitive.ObjectID `json:"workoutplans"`
}

type WorkoutPlan struct {
	ID       primitive.ObjectID   `json:"_id"`
	Name     string               `json:"name"`
	Coach    primitive.ObjectID   `json:"coach"`
	Athletes primitive.ObjectID   `json:"athletes"`
	Workouts []primitive.ObjectID `json:"workouts"`
}
