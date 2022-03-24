package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/models"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CoachLogin(res http.ResponseWriter, req *http.Request) {
	loginInfo := &LoginBody{}
	json.NewDecoder(req.Body).Decode(loginInfo)
	err := coachLogin(loginInfo, req, res)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func coachLogin(loginInfo *LoginBody, req *http.Request, res http.ResponseWriter) error {
	session, _ := store.Get(req, CoachSessionName)

	// fetching coach from emailprovided so we can compare pw's
	coach := &models.Coach{}
	coachErr := coachCollection.FindOne(context.Background(), bson.M{"personalInfo.email": string(loginInfo.Email)}).Decode(&coach)

	// checking to see if anything found and throwing err if nothing found
	if coachErr == mongo.ErrNoDocuments {
		errString := "no coach found with that email"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	// function returns true/false based on err or no err
	//comparing provided password along with pw from db
	match := CheckPasswordHash(loginInfo.Password, coach.PersonalInfo.Password)

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := Claims{
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		coach.PersonalInfo.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		errString := "issue with jwt token"
		http.Error(res, errString, http.StatusForbidden)
		return errors.New(errString)
	}

	// creating session or throwing error if pw match/doesn't match
	if match {
		http.SetCookie(res, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
		session.Values["email"] = coach.PersonalInfo.Email
		session.Values["id"] = coach.ID.Hex()
		session.Save(req, res)
		json.NewEncoder(res).Encode(coach)
		return nil
	} else {
		errString := "invalid password"
		http.Error(res, errString, http.StatusForbidden)
		return errors.New(errString)
	}
}

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
	duplicateEmailErr := coachCollection.FindOne(context.Background(), bson.M{"personalInfo.email": string(coach.PersonalInfo.Email)})
	duplicatePhoneErr := coachCollection.FindOne(context.Background(), bson.M{"personalInfo.phoneNumber": string(coach.PersonalInfo.PhoneNumber)})

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

	// TODO do i need to do something with this error
	hash, _ := HashPassword(coach.PersonalInfo.Password)
	coach.PersonalInfo.Password = hash

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
	clientErr := clientCollection.FindOne(context.Background(), bson.M{"personalInfo.email": string(body.ClientEmail)}).Decode(&clientResult)

	coachResult := &models.Coach{}
	coachErr := coachCollection.FindOne(context.Background(), bson.M{"personalInfo.email": string(body.CoachEmail)}).Decode(&coachResult)

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
