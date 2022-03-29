package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type createWorkoutPlanBody struct {
	Name string `json:"name"`
}

type addNewWorkoutBody struct {
	WorkoutPlanID primitive.ObjectID `json:"workoutPlanId"`
}

type addExistingWorkoutBody struct {
	WorkoutPlanID primitive.ObjectID `json:"workoutPlanId"`
	WorkoutID     primitive.ObjectID `json:"workoutId"`
}

type addExerciseToWorkoutBody struct {
	WorkoutID   primitive.ObjectID `json:"workoutID"`
	Exercise    models.Exercise    `json:"exercise"`
	Sets        int16              `json:"sets"`
	Reps        int16              `json:"reps"`
	Weight      int16              `json:"weight"`
	Description string             `json:"description"`
}

func CreateWorkoutPlan(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var workoutPlanBody createWorkoutPlanBody
	var workoutPlan models.WorkoutPlan
	json.NewDecoder(req.Body).Decode(&workoutPlanBody)

	workoutPlan.ID = primitive.NewObjectID()
	workoutPlan.Name = workoutPlanBody.Name

	err := createWorkoutPlan(workoutPlan, res, req)
	if err != nil {
		fmt.Println(err)

		return
	}
	json.NewEncoder(res).Encode(workoutPlan)
}

func createWorkoutPlan(workoutPlan models.WorkoutPlan, res http.ResponseWriter, req *http.Request) error {
	// checking to see if coach has valid jwt
	TokenCheck(res, req)
	// creating session if coach is logged in
	session, _ := store.Get(req, CoachSessionName)

	// converting id stored in session back to objevt id to search db
	value := session.Values["id"]
	str := fmt.Sprintf("%v", value)
	coachID, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return err
	}

	workoutPlan.CoachID = coachID
	var plan models.WorkoutPlan
	// search workout plan by coachid and name
	// if found reject
	duplicatePlanErr := workoutPlanCollection.FindOne(context.Background(), bson.M{"coachID": coachID, "name": workoutPlan.Name}).Decode((plan))

	if duplicatePlanErr != mongo.ErrNoDocuments {
		errString := "you already have a program with this name"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	// TODO deal with insert err
	_, insertErr := workoutPlanCollection.InsertOne(context.Background(), workoutPlan)
	if insertErr != nil {
		return nil
	}

	coachResult := &models.Coach{}
	coachErr := coachCollection.FindOne(context.Background(), bson.M{"_id": coachID}).Decode(&coachResult)

	if coachErr == mongo.ErrNoDocuments {
		errString := "coach was not found"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	newWorkoutPlans := append(coachResult.WorkoutPlans, workoutPlan.ID)
	coachUpdate := bson.M{
		"$set": bson.M{
			"workoutPlans": newWorkoutPlans,
		},
	}

	_, coachUpdateErr := coachCollection.UpdateByID(context.Background(), coachID, coachUpdate)

	if coachUpdateErr != nil {
		return coachUpdateErr
	}

	return nil
}

func AddNewWorkoutToPlan(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var addNewWorkoutBody addNewWorkoutBody
	var workout models.Workout
	json.NewDecoder(req.Body).Decode(&addNewWorkoutBody)

	workout.ID = primitive.NewObjectID()

	err := addNewWorkoutToPlan(addNewWorkoutBody.WorkoutPlanID, &workout, req, res)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("succesfully added new workout to plan")

	json.NewEncoder(res).Encode(addNewWorkoutBody)
}

func addNewWorkoutToPlan(workoutPlanID primitive.ObjectID, workout *models.Workout, req *http.Request, res http.ResponseWriter) error {
	// checking to see if coach has jwt
	TokenCheck(res, req)

	_, createWorkoutErr := workoutCollection.InsertOne(context.Background(), workout)

	if createWorkoutErr != nil {
		errString := "issue with inserting workout to workout document"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	var workoutPlan = &models.WorkoutPlan{}

	findWorkoutPlanErr := workoutPlanCollection.FindOne(context.Background(), bson.M{"_id": workoutPlanID}).Decode(&workoutPlan)

	if findWorkoutPlanErr == mongo.ErrNoDocuments {
		errString := "no workout plan found"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	workoutPlan.Workouts = append(workoutPlan.Workouts, workout.ID)

	workoutPlanUpdate := bson.M{
		"$set": bson.M{
			"Workouts": workoutPlan.Workouts,
		},
	}

	// checking for error in updating workoutplan and setting error equal to what is returned
	_, workoutPlanUpdateErr := workoutPlanCollection.UpdateByID(context.Background(), workoutPlan.ID, workoutPlanUpdate)

	if workoutPlanUpdateErr != nil {
		return workoutPlanUpdateErr
	}

	return nil
}

func AddExistingWorkoutToPlan(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var addExistingWorkoutBody addExistingWorkoutBody
	json.NewDecoder(req.Body).Decode(&addExistingWorkoutBody)

	err := addExistingWorkoutToPlan(addExistingWorkoutBody, req, res)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("succesfully added existing workout to plan")
}

func addExistingWorkoutToPlan(body addExistingWorkoutBody, req *http.Request, res http.ResponseWriter) error {
	// checking to see if coach has jwt
	TokenCheck(res, req)

	var workoutPlan = &models.WorkoutPlan{}

	findWorkoutPlanErr := workoutPlanCollection.FindOne(context.Background(), bson.M{"_id": body.WorkoutPlanID}).Decode(&workoutPlan)

	if findWorkoutPlanErr == mongo.ErrNoDocuments {
		errString := "no workout plan found"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	workoutPlan.Workouts = append(workoutPlan.Workouts, body.WorkoutID)

	workoutPlanUpdate := bson.M{
		"$set": bson.M{
			"Workouts": workoutPlan.Workouts,
		},
	}

	// checking for error in updating workoutplan and setting error equal to what is returned
	_, workoutPlanUpdateErr := workoutPlanCollection.UpdateByID(context.Background(), workoutPlan.ID, workoutPlanUpdate)

	if workoutPlanUpdateErr != nil {
		return workoutPlanUpdateErr
	}

	return nil

}

func AddExerciseToWorkout(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var addExerciseToWorkoutBody addExerciseToWorkoutBody
	json.NewDecoder(req.Body).Decode(&addExerciseToWorkoutBody)

	err := addExerciseToWorkout(addExerciseToWorkoutBody, req, res)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("succesfully added workout to plan")

	json.NewEncoder(res).Encode(exerciseDetails)
}

func addExerciseToWorkout(body addExerciseToWorkoutBody, req *http.Request, res http.ResponseWriter) error {
	// checking to see if coach has jwt
	TokenCheck(res, req)

	exerciseDetails := models.ExerciseDetails{
		ID:          primitive.NewObjectID(),
		Exercise:    body.Exercise,
		Sets:        body.Sets,
		Reps:        body.Reps,
		Weight:      body.Weight,
		Description: body.Description,
	}

	workout := &models.Workout{}

	//find the workout
	findWorkoutErr := workoutCollection.FindOne(context.Background(), bson.M{"_id": body.WorkoutID}).Decode(&workout)

	if findWorkoutErr == mongo.ErrNoDocuments {
		errString := "workout was not found"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	workout.ExercisesDetails = append(workout.ExercisesDetails, exerciseDetails)

	workoutUpdate := bson.M{
		"$set": bson.M{
			"exercisesDetails": workout.ExercisesDetails,
		},
	}

	_, updateWorkoutErr := workoutCollection.UpdateByID(context.Background(), workout.ID, workoutUpdate)

	if updateWorkoutErr != nil {
		errString := "issue updating workout with exercise details"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	return nil
}

func GetWorkoutPlanDetails(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	workoutPlanID := params["id"]

	workoutDetails, err := getWorkoutPlanDetails(req, res, workoutPlanID)

	if err != nil {
		return
	}

	json.NewEncoder(res).Encode(workoutDetails)
}

func getWorkoutPlanDetails(req *http.Request, res http.ResponseWriter, workoutPlanID string) ([]models.Workout, error) {
	TokenCheck(res, req)

	workouts := []models.Workout{}
	// turning string given from params back into objectID
	idPrimitive, _ := primitive.ObjectIDFromHex(workoutPlanID)

	workoutPlan := &models.WorkoutPlan{}
	err := workoutPlanCollection.FindOne(context.Background(), bson.M{"_id": idPrimitive}).Decode(&workoutPlan)

	if err != nil {
		fmt.Println(err)
		errString := "problem finding workout plan"
		return nil, errors.New(errString)
	}

	workoutIds := workoutPlan.Workouts
	//TODO improve error handling here
	for _, workoutId := range workoutIds {
		workout := &models.Workout{}
		workoutErr := workoutCollection.FindOne(context.Background(), bson.M{"_id": workoutId}).Decode(&workout)
		if workoutErr != nil {
			fmt.Println(workoutErr)
			return nil, workoutErr
		}
		workouts = append(workouts, *workout)
	}
	return workouts, nil

}
