package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateExercise(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var exercise models.Exercise
	json.NewDecoder(req.Body).Decode(&exercise)
	exercise.ID = primitive.NewObjectID()

	// Now that createExercise returns an error, this CreateExercise func can know about it
	err := createExercise(exercise, res)
	if err != nil {
		// We can print the error, or if you use a logging tool you'll log it here
		fmt.Println(err)

		// Since we encountered an error we can return early
		return
	}

	// Since we made it to this point we can log out the success message
	fmt.Println("added exercise:", exercise.Name)

	// This does not actually return the Exercise. It just does an HTTP response to the caller.
	json.NewEncoder(res).Encode(exercise)
}

func createExercise(exercise models.Exercise, res http.ResponseWriter) error {
	// response is the SingleResult that mongo sends back if it finds a document that matches the exercise.Name
	response := exerciseCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "name", Value: string(exercise.Name)}})

	// Here we check if the response has an Err
	// This is where the reverse expectation stuff is
	// The response is SingleResult type and it has a function Err() that returns the error.
	// Here we check if that err is NOT equal to ErrNoDocuments...
	// ... if that is true then we found a document (not equal to no documents means yes we found one)
	// If response was equal to ErrNoDocument then we move past this block and add it to mongo.
	if response.Err() != mongo.ErrNoDocuments {
		// We can use the errString in both the HTTP response and the returned err
		errString := "exercise already exists"

		// This will send the HTTP response back, but it does not return the error
		http.Error(res, errString, 400)

		// This returns the Err in case the function that calls createExercise wants to know about the error.
		return errors.New(errString)
	}

	// Insert the exercise into mongo. The err will be nil if successful.
	_, err := exerciseCollection.InsertOne(context.Background(), exercise)

	// Return the err.
	return err
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
	exerciseDetails.Exercise_ID = body.ID

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
