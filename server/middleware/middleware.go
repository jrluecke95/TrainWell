package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"server/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

// collection object/instance
var coachCollection *mongo.Collection
var clientCollection *mongo.Collection
var exerciseCollection *mongo.Collection

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

	fmt.Println("Collection instance created!")
}

func CreateExercise(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var exercise models.Exercise
	json.NewDecoder(req.Body).Decode(&exercise)
	exercise.ID = primitive.NewObjectID()
	createExercise(exercise, res)
	json.NewEncoder(res).Encode(exercise)
}

func createExercise(exercise models.Exercise, res http.ResponseWriter) {
	var result = &models.Exercise{}
	duplicateErr := exerciseCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "name", Value: string(exercise.Name)}}).Decode(&result)

	if duplicateErr == nil {
		if duplicateErr == mongo.ErrNoDocuments {
			return
		}
		http.Error(res, "Exercise already exists", 400)
		return
	}

	_, err := exerciseCollection.InsertOne(context.Background(), exercise)

	if err != nil {
		log.Fatal(err)
	}
}

func GetExercises(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	payload := getAllExercises()
	json.NewEncoder(res).Encode(payload)
}

func getAllExercises() []primitive.M {
	cur, err := exerciseCollection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func CreateCoach(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var coach models.Coach
	json.NewDecoder(req.Body).Decode(&coach)
	coach.ID = primitive.NewObjectID()
	createCoach(coach, res)
	json.NewEncoder(res).Encode(coach)
}

func createCoach(coach models.Coach, res http.ResponseWriter) {
	var result = &models.Coach{}
	duplicateEmailErr := coachCollection.FindOne(context.Background(), bson.M{"email": string(coach.Email)}).Decode(&result)
	duplicatePhoneErr := coachCollection.FindOne(context.Background(), bson.M{"phonenumber": string(coach.PhoneNumber)}).Decode(&result)

	if duplicateEmailErr == nil {
		http.Error(res, "Email already in use for coach", 400)
		return
	}

	if duplicatePhoneErr == nil {
		http.Error(res, "Phonenumber already in use for coach", 400)
		return
	}

	_, err := coachCollection.InsertOne(context.Background(), coach)

	if err != nil {
		log.Fatal(err)
	}
}

func GetCoaches(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	payload := getAllCoaches()
	json.NewEncoder(res).Encode(payload)
}

func getAllCoaches() []primitive.M {
	cur, err := coachCollection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func CreateClient(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var client models.Client
	json.NewDecoder(req.Body).Decode(&client)
	client.ID = primitive.NewObjectID()
	createClient(client, res)
	json.NewEncoder(res).Encode(client)
}

func createClient(client models.Client, res http.ResponseWriter) {
	var result = &models.Client{}
	duplicateEmailErr := clientCollection.FindOne(context.Background(), bson.M{"email": string(client.Email)}).Decode(&result)
	duplicatePhoneErr := clientCollection.FindOne(context.Background(), bson.M{"phonenumber": string(client.PhoneNumber)}).Decode(&result)

	if duplicateEmailErr == nil {
		http.Error(res, "Email already in use for client", 400)
		return
	}

	if duplicatePhoneErr == nil {
		http.Error(res, "Phonenumber already in use for client", 400)
		return
	}

	_, err := clientCollection.InsertOne(context.Background(), client)

	if err != nil {
		log.Fatal(err)
	}
}

func GetClients(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	payload := getAllClients()
	json.NewEncoder(res).Encode(payload)
}

func getAllClients() []primitive.M {
	cur, err := clientCollection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func AssignCoach(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var client models.Client
	json.NewDecoder(req.Body).Decode(&client)
	client.ID = primitive.NewObjectID()
	assignCoach(client, res, req)
	json.NewEncoder(res).Encode(client)
}

func assignCoach(client models.Client, res http.ResponseWriter, req *http.Request) {
	// find client first
	// how to pull info from request body?
	// do i need to pass the whole client through?
	// or does it make more sense to pass email/phone/id and use that to find client?
	// db look up or passing larger object more efficient?
	//var client2 models.Client
	// json.NewDecoder(req.Body).Decode(&client2)
	// fmt.Println(client2.F)
	// append coach id to client coaches array
}
