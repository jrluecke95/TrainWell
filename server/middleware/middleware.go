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
)

type AssignCoachBody struct {
	ClientEmail string `json:"clientEmail" bson:"clientEmail"`
	CoachEmail  string `json:"coachEmail" bson:"coachEmail"`
}
type ExerciseDetailsBody struct {
	ID   primitive.ObjectID `json:"id"`
	Sets int16              `json:"sets"`
	Reps int16              `json:"reps"`
}

func goDotEnvVariable(key string) string {

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
var value = goDotEnvVariable("mongodbConnectString")
var connectionString = value

// Database Name
const dbName = "trainwell"

// Collection names
const coaches = "Coaches"
const clients = "Clients"
const exercises = "Exercises"
const exerciseDetails = "ExerciseDetails"

// collection object/instance
var coachCollection *mongo.Collection
var clientCollection *mongo.Collection
var exerciseCollection *mongo.Collection
var exerciseDetailsCollection *mongo.Collection

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
