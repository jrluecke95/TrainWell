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
	Name string `json:"name"`
}

type addWorkoutBody struct {
	WorkoutPlanID primitive.ObjectID `json:"workoutPlanId"`
}

type addExerciseToWorkoutBody struct {
	WorkoutPlanID primitive.ObjectID `json:"workoutPlanId`
	WorkoutID     primitive.ObjectID `json:"workoutID"`
	Exercise      models.Exercise    `json:"exercise"`
	Sets          int16              `json:"sets"`
	Reps          int16              `json:"reps"`
	Weight        int16              `json:"weight"`
	Description   string             `json:"description"`
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

	fmt.Println("succesfully created workout plan")

	json.NewEncoder(res).Encode(workoutPlan)
}

func createWorkoutPlan(workoutPlan models.WorkoutPlan, res http.ResponseWriter, req *http.Request) error {
	// checking to see if coach is logged in
	if !CheckLogin(res, req, CoachSessionName) {
		errString := "not logged in"

		http.Error(res, errString, http.StatusForbidden)
		return errors.New(errString)
	}
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

// func AddWorkoutToPlan(res http.ResponseWriter, req *http.Request) {
// 	res.Header().Add("content-type", "application/json")
// 	var addWorkoutBody addWorkoutBody
// 	var workout models.Workout
// 	json.NewDecoder(req.Body).Decode(&addWorkoutBody)

// 	workout.ID = primitive.NewObjectID()

// 	err := addWorkoutToPlan(addWorkoutBody.WorkoutPlanID, &workout, req, res)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println("succesfully added workout to plan")

// 	json.NewEncoder(res).Encode(addWorkoutBody)
// }

// func addWorkoutToPlan(workoutPlanID primitive.ObjectID, workout *models.Workout, req *http.Request, res http.ResponseWriter) error {
// 	session, _ := store.Get(req, CoachSessionName)

// 	// finding coach to access workout plans
// 	var coachResult models.Coach

// 	findCoachErr := coachCollection.FindOne(context.Background(), bson.M{"personalInfo.email": session.Values["email"]}).Decode(&coachResult)

// 	if findCoachErr != nil {
// 		return findCoachErr
// 	}

// 	// find workout plan
// 	var findWorkoutPlan *models.WorkoutPlan
// 	for _, workoutPlan := range coachResult.WorkoutPlans {
// 		if workoutPlan.ID == workoutPlanID {
// 			findWorkoutPlan = &workoutPlan
// 		}
// 	}

// 	findWorkoutPlan.Workouts = append(findWorkoutPlan.Workouts, *workout)

// 	workoutPlanUpdate := bson.M{
// 		"$set": bson.M{
// 			"workoutPlans": &findWorkoutPlan,
// 		},
// 	}

// 	// checking for error in updating workoutplan and setting error equal to what is returned
// 	_, coachUpdateErr := coachCollection.UpdateByID(context.Background(), coachResult.ID, workoutPlanUpdate)

// 	if coachUpdateErr != nil {
// 		return coachUpdateErr
// 	}
// 	return nil
// }

// func AddExerciseToWorkout(res http.ResponseWriter, req *http.Request) {
// 	res.Header().Add("content-type", "application/json")
// 	var addExerciseToWorkoutBody addExerciseToWorkoutBody
// 	var exerciseDetails models.ExerciseDetails
// 	json.NewDecoder(req.Body).Decode(&addExerciseToWorkoutBody)

// 	exerciseDetails.ID = primitive.NewObjectID()

// 	err := addExerciseToWorkout(addExerciseToWorkoutBody, req, res)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println("succesfully added workout to plan")

// 	json.NewEncoder(res).Encode(exerciseDetails)
// }

// func addExerciseToWorkout(body addExerciseToWorkoutBody, req *http.Request, res http.ResponseWriter) error {
// 	// checking to see if coach is logged in
// 	// TODO this through the loggin coacherr should probably be put into a function becasue it is going to be reuse dalot
// 	if !CheckLogin(res, req, CoachSessionName) {
// 		errString := "not logged in"

// 		http.Error(res, errString, http.StatusForbidden)
// 		return errors.New(errString)
// 	}
// 	// creating session if coach is logged in
// 	session, _ := store.Get(req, CoachSessionName)

// 	// converting id stored in session back to objevt id to search db
// 	value := session.Values["id"]
// 	str := fmt.Sprintf("%v", value)
// 	coachID, err := primitive.ObjectIDFromHex(str)
// 	if err != nil {
// 		return err
// 	}

// 	// finding coach with session info
// 	var coach = &models.Coach{}
// 	coachErr := coachCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: coachID}}).Decode(&coach)

// 	if coachErr == mongo.ErrNoDocuments {
// 		errString := "no coach found"
// 		http.Error(res, errString, 400)
// 		return errors.New(errString)
// 	}

// 	//TODO rework db

// 	//find the workout plan
// 	//update the workout
// 	//put the workout back in
// 	// update whole workout plan array

// 	// need to locate workout within workout plan and update that array
// 	// if workoutErr == mongo.ErrNoDocuments {
// 	// 	errString := "no workout found"
// 	// 	http.Error(res, errString, 400)
// 	// 	return errors.New(errString)
// 	// }

// 	// fmt.Println("result", workoutResult)

// 	return nil
// }
