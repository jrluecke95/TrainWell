package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	json.NewDecoder(req.Body).Decode((&exercise))
	exercise.ID = primitive.NewObjectID()
	createExercise(exercise)
	json.NewEncoder(res).Encode((exercise))
}

func GetExercises(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	payload := getAllExercises()
	json.NewEncoder(res).Encode(payload)
}

func createExercise(exercise models.Exercise) {

	var result = &models.Exercise{}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	duplicateErr := exerciseCollection.FindOne(ctx, bson.D{primitive.E{Key: "name", Value: string(exercise.Name)}}).Decode(&result)

	if duplicateErr == nil {
		if duplicateErr == mongo.ErrNoDocuments {
			return
		}
		fmt.Println("400 error exercise alrady exists")
		return
	}

	_, err := exerciseCollection.InsertOne(context.Background(), exercise)

	if err != nil {
		log.Fatal(err)
	}
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

// get all task from the DB and return it
// func getAllTask() []primitive.M {
// 	cur, err := collection.Find(context.Background(), bson.D{{}})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var results []primitive.M
// 	for cur.Next(context.Background()) {
// 		var result bson.M
// 		e := cur.Decode(&result)
// 		if e != nil {
// 			log.Fatal(e)
// 		}
// 		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
// 		results = append(results, result)

// 	}

// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	cur.Close(context.Background())
// 	return results
// }

// // Insert one task in the DB
// func insertOneTask(task models.ToDoList) {
// 	insertResult, err := collection.InsertOne(context.Background(), task)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
// }

// // task complete method, update task's status to true
// func taskComplete(task string) {
// 	fmt.Println(task)
// 	id, _ := primitive.ObjectIDFromHex(task)
// 	filter := bson.M{"_id": id}
// 	update := bson.M{"$set": bson.M{"status": true}}
// 	result, err := collection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("modified count: ", result.ModifiedCount)
// }

// // task undo method, update task's status to false
// func undoTask(task string) {
// 	fmt.Println(task)
// 	id, _ := primitive.ObjectIDFromHex(task)
// 	filter := bson.M{"_id": id}
// 	update := bson.M{"$set": bson.M{"status": false}}
// 	result, err := collection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("modified count: ", result.ModifiedCount)
// }

// // delete one task from the DB, delete by ID
// func deleteOneTask(task string) {
// 	fmt.Println(task)
// 	id, _ := primitive.ObjectIDFromHex(task)
// 	filter := bson.M{"_id": id}
// 	d, err := collection.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Deleted Document", d.DeletedCount)
// }

// // delete all the tasks from the DB
// func deleteAllTask() int64 {
// 	d, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Deleted Document", d.DeletedCount)
// 	return d.DeletedCount
// }
