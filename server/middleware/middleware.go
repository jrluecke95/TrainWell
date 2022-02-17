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

	// need to encrypt/decrypt password
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
	var body AssignCoachBody
	json.NewDecoder(req.Body).Decode(&body)
	assignCoach(body, res, req)
	json.NewEncoder(res).Encode(body)
}

func assignCoach(body AssignCoachBody, res http.ResponseWriter, req *http.Request) {
	clientResult := &models.Client{}
	clientCollection.FindOne(context.Background(), bson.M{"email": string(body.ClientEmail)}).Decode(&clientResult)

	coachResult := &models.Coach{}
	coachCollection.FindOne(context.Background(), bson.M{"email": string(body.CoachEmail)}).Decode(&coachResult)

	newCoaches := append(clientResult.Coaches, coachResult.ID)
	coachUpdate := bson.M{
		"$set": bson.M{
			"coaches": newCoaches,
		},
	}

	clientCollection.UpdateByID(context.Background(), clientResult.ID, coachUpdate)

	// need to add error handling for non-existent client
	// if needed

	newClients := append(coachResult.Clients, clientResult.ID)
	clientUpdate := bson.M{
		"$set": bson.M{
			"clients": newClients,
		},
	}

	coachCollection.UpdateByID(context.Background(), coachResult.ID, clientUpdate)
}

func UnassignCoach(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var body AssignCoachBody
	json.NewDecoder(req.Body).Decode(&body)
	unassignCoach(body, res, req)
	json.NewEncoder(res).Encode("Coach removal successful")
}

func unassignCoach(body AssignCoachBody, res http.ResponseWriter, req *http.Request) {
	clientResult := &models.Client{}
	clientCollection.FindOne(context.Background(), bson.M{"email": string(body.ClientEmail)}).Decode(&clientResult)

	coachResult := &models.Coach{}
	coachCollection.FindOne(context.Background(), bson.M{"email": string(body.CoachEmail)}).Decode(&coachResult)

	newCoaches := removeIDFromArray(clientResult.Coaches, coachResult.ID)
	coachUpdate := bson.M{
		"$set": bson.M{
			"coaches": newCoaches,
		},
	}

	clientCollection.UpdateByID(context.Background(), clientResult.ID, coachUpdate)

	newClients := removeIDFromArray(coachResult.Clients, clientResult.ID)
	clientUpdate := bson.M{
		"$set": bson.M{
			"clients": newClients,
		},
	}

	coachCollection.UpdateByID(context.Background(), coachResult.ID, clientUpdate)
}

func CreateExerciseDetails(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var body ExerciseDetailsBody
	var exerciseDetails models.ExerciseDetails
	json.NewDecoder(req.Body).Decode(&body)
	createExerciseDetails(body, exerciseDetails, res, req)
}

func createExerciseDetails(body ExerciseDetailsBody, exerciseDetails models.ExerciseDetails, res http.ResponseWriter, req *http.Request) {
	exercise := models.Exercise{}
	exerciseCollection.FindOne(context.Background(), bson.M{"_id": body.ID}).Decode(&exercise)

	exerciseDetails.ID = primitive.NewObjectID()
	exerciseDetails.Reps = body.Reps
	exerciseDetails.Sets = body.Sets
	exerciseDetails.Name_ID = body.ID

	details := append(exercise.ExerciseDetails, exerciseDetails.ID)

	exerciseUpdate := bson.M{
		"$set": bson.M{
			"exercisedetails": details,
		},
	}

	exerciseCollection.UpdateByID(context.Background(), body.ID, exerciseUpdate)

	_, err := exerciseDetailsCollection.InsertOne(context.Background(), exerciseDetails)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(res).Encode(exerciseDetails)
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
