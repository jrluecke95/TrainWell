package middleware

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AssignCoachBody struct {
	ClientEmail string `json:"clientEmail" bson:"clientEmail"`
	CoachEmail  string `json:"coachEmail" bson:"coachEmail"`
}
type ExerciseDetailsBody struct {
	WorkoutId    primitive.ObjectID `json:"workoutId"`
	ExerciseName string             `json:"exerciseName"`
	Sets         int16              `json:"sets"`
	Reps         int16              `json:"reps"`
	Weight       int16              `json:"weight"`
	Description  string             `json:"description"`
}

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// DB connection string
// for localhost mongoDB
// const connectionString = "mongodb://localhost:27017"
var value = GoDotEnvVariable("mongodbConnectString")
var connectionString = value

// Database Name
const dbName = "trainwell"

// Collection names
const coaches = "Coaches"
const clients = "Clients"
const exercises = "Exercises"
const exerciseDetails = "ExerciseDetails"
const workoutPlans = "WorkoutPlans"
const workouts = "Workouts"

// collection object/instance
var coachCollection *mongo.Collection
var clientCollection *mongo.Collection
var exerciseCollection *mongo.Collection
var exerciseDetailsCollection *mongo.Collection
var workoutPlanCollection *mongo.Collection
var workoutCollection *mongo.Collection

// create connection with mongo db
func init() {

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	coachCollection = client.Database(dbName).Collection(coaches)
	clientCollection = client.Database(dbName).Collection(clients)
	exerciseCollection = client.Database(dbName).Collection((exercises))
	exerciseDetailsCollection = client.Database(dbName).Collection(exerciseDetails)
	workoutCollection = client.Database(dbName).Collection(workouts)
	workoutPlanCollection = client.Database(dbName).Collection(workoutPlans)

	fmt.Println("Collection instance created!")
}

func removeIDFromArray(initArray []primitive.ObjectID, badId primitive.ObjectID) []primitive.ObjectID {
	var finalArr []primitive.ObjectID
	for _, s := range initArray {
		if s != badId {
			finalArr = append(finalArr, s)
		}
	}
	return finalArr
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
