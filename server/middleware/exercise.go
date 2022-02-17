package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateExercise(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var exercise models.Exercise
	json.NewDecoder(req.Body).Decode(&exercise)
	exercise.ID = primitive.NewObjectID()
	createExercise(exercise, res)
	json.NewEncoder(res).Encode(exercise)
}

func createExercise(exercise models.Exercise, res http.ResponseWriter) error {
	var result = &models.Exercise{}
	duplicateErr := exerciseCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "name", Value: string(exercise.Name)}}).Decode(&result)
	fmt.Println(duplicateErr)
	fmt.Println(result)

	if duplicateErr == nil {
		http.Error(res, "Exercise already exists", 400)
		fmt.Println(duplicateErr)
		return duplicateErr
	}

	_, err := exerciseCollection.InsertOne(context.Background(), exercise)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
