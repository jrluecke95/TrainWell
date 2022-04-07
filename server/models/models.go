package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// i think we will need to make copies of workout plans and store them in the client
// this allows for them to modify if they want and not effect coach workout as well as tracking weights etc. easier
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
	WorkoutPlans []primitive.ObjectID `json:"workoutPlans" bson:"workoutPlans"`
	PersonalInfo PersonalInfo         `json:"personalInfo" bson:"personalInfo"`
}

type PersonalInfo struct {
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	Password    string `gorm:"type:varchar(100) json:"password" bson:"password"`
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
	CoachID  primitive.ObjectID   `json:"coachID" bson:"coachID"`
	Workouts []primitive.ObjectID `json:"workouts" bson:"workouts"`
	Name     string               `json:"name" bson:"name"`
}
