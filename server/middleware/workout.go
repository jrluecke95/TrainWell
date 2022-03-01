package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type createWorkoutPlanBody struct {
	Name  string
	Coach primitive.ObjectID
}

func CreateWorkoutPlan(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var workoutPlanBody createWorkoutPlanBody
	var workoutPlan models.WorkoutPlan
	json.NewDecoder(req.Body).Decode(&workoutPlanBody)

	workoutPlan.ID = primitive.NewObjectID()
	workoutPlan.Name = workoutPlanBody.Name

	err := createWorkoutPlan(workoutPlan, res)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println("succesfully created workout")

	json.NewEncoder(res).Encode(workoutPlan)
}

func createWorkoutPlan(workoutPlan models.WorkoutPlan, res http.ResponseWriter) error {
	response := workoutPlanCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "name", Value: string(workoutPlan.Name)}})

	//TODO use session to get coach that is updating/creating plan and then lookup workout by name

	if response.Err() != mongo.ErrNoDocuments {
		errString := "workout plan with that name already exists"

		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	_, insertErr := workoutPlanCollection.InsertOne(context.Background(), workoutPlan)

	if insertErr != nil {
		return insertErr
	}

	// TODO should all be replaced with session info
	// var coachResult models.Coach

	// findCoachErr := coachCollection.FindOne(context.Background(), bson.M{"_id": workoutPlan.Coach}).Decode(&coachResult)

	// if findCoachErr != nil {
	// 	return findCoachErr
	// }

	// newWorkoutPlans := append(coachResult.WorkoutPlans, workoutPlan)

	// workoutPlanUpdate := bson.M{
	// 	"$set": bson.M{
	// 		"workoutPlans": newWorkoutPlans,
	// 	},
	// }

	// _, coachUpdateErr := coachCollection.UpdateByID(context.Background(), coachResult.ID, workoutPlanUpdate)

	// if coachUpdateErr != nil {
	// 	return coachUpdateErr
	// }

	return nil
}
