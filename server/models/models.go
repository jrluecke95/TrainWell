package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Client struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id"`
	Coaches      []primitive.ObjectID `json:"coaches" bson:"coaches"`
	WorkoutPlans []primitive.ObjectID `json:"workoutPlans" bson:"workoutPlans"`
	PersonalInfo PersonalInfo         `json:"personalinfo" bson:"personalInfo"`
}
type Coach struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id"`
	Clients      []primitive.ObjectID `json:"clients" bson:"clients"`
	Workouts     []Workout            `json:"workouts" bson:"workouts"`
	WorkoutPlans []WorkoutPlan        `json:"workoutPlans" bson:"workoutPlans"`
	PersonalInfo PersonalInfo         `json:"personalInfo" bson:"personalInfo"`
}

type PersonalInfo struct {
	FirstName   string `json:"firstnName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	Password    string `json:"password" bson:"password"`
}

type Exercise struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type ExerciseDetails struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Exercise    Exercise           `json:"exercise" bson:"exercise"`
	Sets        int16              `json:"sets" bson:"sets"`
	Reps        int16              `json:"reps" bson:"reps"`
	Weight      int16              `json:"weight" bson:"weight"`
	Description string             `json:"description" bson:"description"`
}
type Workout struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
	ExercisesDetails []ExerciseDetails  `json:"exercisesDetails" bson:"exercisesDetails"`
}

type WorkoutPlan struct {
	ID       primitive.ObjectID   `json:"_id" bson:"_id"`
	Name     string               `json:"name" bson:"name"`
	Workouts []primitive.ObjectID `json:"workouts" bson:"workouts"`
}
