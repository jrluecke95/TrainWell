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
	Workouts     []Workout            `json:"workouts"`
	WorkoutPlans []WorkoutPlan        `json:"workoutplans"`
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
}

type ExerciseDetails struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Exercise_ID primitive.ObjectID `json:"exerciseid"`
	Sets        int16              `json:"sets"`
	Reps        int16              `json:"reps"`
	Weight      int16              `json:"weight"`
	Description string             `json:"description"`
}

type Workout struct {
	ID               primitive.ObjectID   `json:"_id" bson:"_id"`
	ExercisesDetails []ExerciseDetails    `json:"exercisesdetails"`
	WorkoutPlans     []primitive.ObjectID `json:"workoutplans"`
}

// do i only store workout plan within the coach document?
// as in no need to update this workoutplan/create its own document within mongo?
// i think that is the case but i really have no idea
// maybe i can just store name and coach id in workoutplan document and use those to reference everything else
type WorkoutPlan struct {
	ID       primitive.ObjectID   `json:"_id"`
	Name     string               `json:"name"`
	Coach    primitive.ObjectID   `json:"coach"`
	Athletes []primitive.ObjectID `json:"athletes"`
	Workouts []primitive.ObjectID `json:"workouts"`
}
