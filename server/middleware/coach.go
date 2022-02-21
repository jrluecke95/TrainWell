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

func CreateCoach(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var coach models.Coach
	json.NewDecoder(req.Body).Decode(&coach)
	coach.ID = primitive.NewObjectID()
	err := createCoach(coach, res)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(res).Encode(coach)
}

func createCoach(coach models.Coach, res http.ResponseWriter) error {
	duplicateEmailErr := coachCollection.FindOne(context.Background(), bson.M{"email": string(coach.Email)})
	duplicatePhoneErr := coachCollection.FindOne(context.Background(), bson.M{"phonenumber": string(coach.PhoneNumber)})

	if duplicateEmailErr.Err() != mongo.ErrNoDocuments {
		errString := "email already in use for coach"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	if duplicatePhoneErr.Err() != mongo.ErrNoDocuments {
		errString := "phone number already in use for coach"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	_, err := coachCollection.InsertOne(context.Background(), coach)

	return err
}

func GetCoaches(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	payload, err := getAllCoaches(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(res).Encode(payload)
}

func getAllCoaches(res http.ResponseWriter) ([]primitive.M, error) {
	cur, err := coachCollection.Find(context.Background(), bson.D{{}})

	if err != nil {
		errString := "error fetching coaches"
		http.Error(res, errString, 400)
		return nil, errors.New(errString)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			return nil, e
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
	}

	cur.Close(context.Background())
	return results, nil
}

func AssignCoach(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var body AssignCoachBody
	json.NewDecoder(req.Body).Decode(&body)
	err := assignCoach(body, res)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(res).Encode(body)
}

func assignCoach(body AssignCoachBody, res http.ResponseWriter) error {
	clientResult := &models.Client{}
	clientErr := clientCollection.FindOne(context.Background(), bson.M{"email": string(body.ClientEmail)}).Decode(&clientResult)

	coachResult := &models.Coach{}
	coachErr := coachCollection.FindOne(context.Background(), bson.M{"email": string(body.CoachEmail)}).Decode(&coachResult)

	// if client is not found throw error
	if clientErr == mongo.ErrNoDocuments {
		errString := "client was not found"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}
	if coachErr == mongo.ErrNoDocuments {
		errString := "coach was not found"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	// initiating error to then set equal to any error that pops up in following code blocks
	var err error

	newCoaches := append(clientResult.Coaches, coachResult.ID)
	coachUpdate := bson.M{
		"$set": bson.M{
			"coaches": newCoaches,
		},
	}

	// checking for error in finding client and setting error equal to what is returned
	_, clientUpdateErr := clientCollection.UpdateByID(context.Background(), clientResult.ID, coachUpdate)

	if clientUpdateErr != nil {
		err = clientUpdateErr
	}

	newClients := append(coachResult.Clients, clientResult.ID)
	clientUpdate := bson.M{
		"$set": bson.M{
			"clients": newClients,
		},
	}

	_, coachUpdateErr := coachCollection.UpdateByID(context.Background(), coachResult.ID, clientUpdate)

	if coachUpdateErr != nil {
		err = coachUpdateErr
	}

	return err
}

func UnassignCoach(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var body AssignCoachBody
	json.NewDecoder(req.Body).Decode(&body)
	err := unassignCoach(body, res)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(res).Encode("Coach removal successful")
}

func unassignCoach(body AssignCoachBody, res http.ResponseWriter) error {
	clientResult := &models.Client{}
	clientErr := clientCollection.FindOne(context.Background(), bson.M{"email": string(body.ClientEmail)}).Decode(&clientResult)

	coachResult := &models.Coach{}
	coachErr := coachCollection.FindOne(context.Background(), bson.M{"email": string(body.CoachEmail)}).Decode(&coachResult)

	// if client is not found throw error
	if clientErr == mongo.ErrNoDocuments {
		errString := "client was not found"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}
	if coachErr == mongo.ErrNoDocuments {
		errString := "coach was not found"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	var err error

	newCoaches := removeIDFromArray(clientResult.Coaches, coachResult.ID)
	coachUpdate := bson.M{
		"$set": bson.M{
			"coaches": newCoaches,
		},
	}

	_, clientUpdateErr := clientCollection.UpdateByID(context.Background(), clientResult.ID, coachUpdate)

	if clientUpdateErr != nil {
		err = clientUpdateErr
	}

	newClients := removeIDFromArray(coachResult.Clients, clientResult.ID)
	clientUpdate := bson.M{
		"$set": bson.M{
			"clients": newClients,
		},
	}

	_, coachUpdateErr := coachCollection.UpdateByID(context.Background(), coachResult.ID, clientUpdate)

	if coachUpdateErr != nil {
		err = coachUpdateErr
	}

	return err
}
