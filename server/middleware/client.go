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

func ClientLogin(res http.ResponseWriter, req *http.Request) {
	loginInfo := &LoginBody{}
	json.NewDecoder(req.Body).Decode(loginInfo)
	err := clientLogin(loginInfo, req, res)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func clientLogin(loginInfo *LoginBody, req *http.Request, res http.ResponseWriter) error {
	session, _ := store.Get(req, ClientSessionName)

	// fetching coach from emailprovided so we can compare pw's
	var client = &models.Client{}
	clientErr := clientCollection.FindOne(context.Background(), bson.M{"personalInfo.email": string(loginInfo.Email)}).Decode(&client)

	// checking to see if anything found and throwing err if nothing found
	if clientErr == mongo.ErrNoDocuments {
		errString := "no client found with that email"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	// function returns true/false based on err or no err
	//comparing provided password along with pw from db
	match := CheckPasswordHash(loginInfo.Password, client.PersonalInfo.Password)

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := Claims{
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		client.PersonalInfo.Email,
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
		session.Values["email"] = client.PersonalInfo.Email
		session.Values["id"] = client.ID.Hex()
		session.Save(req, res)
		// this is done here to avoid sending headers prematurely ahead of cookie being set
		json.NewEncoder(res).Encode(client)
		return nil
	} else {
		errString := "invalid password"
		http.Error(res, errString, http.StatusForbidden)
		return errors.New(errString)
	}
}

func CreateClient(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var client models.Client
	json.NewDecoder(req.Body).Decode(&client)
	client.ID = primitive.NewObjectID()
	err := createClient(client, res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("added client:", client.PersonalInfo.FirstName)
	json.NewEncoder(res).Encode(client)
}

func createClient(client models.Client, res http.ResponseWriter) error {
	duplicateEmailErr := clientCollection.FindOne(context.Background(), bson.M{"email": string(client.PersonalInfo.Email)})
	duplicatePhoneErr := clientCollection.FindOne(context.Background(), bson.M{"phonenumber": string(client.PersonalInfo.PhoneNumber)})

	if duplicateEmailErr.Err() != mongo.ErrNoDocuments {
		errString := "email already in use for client"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	if duplicatePhoneErr.Err() != mongo.ErrNoDocuments {
		errString := "phone already in use for client"
		http.Error(res, errString, 400)
		return errors.New(errString)
	}

	// TODO do i need to do something with this error
	hash, _ := HashPassword(client.PersonalInfo.Password)
	client.PersonalInfo.Password = hash

	_, err := clientCollection.InsertOne(context.Background(), client)

	return err
}

func GetClients(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	payload, err := getAllClients(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(res).Encode(payload)
}

func getAllClients(res http.ResponseWriter) ([]primitive.M, error) {
	cur, err := clientCollection.Find(context.Background(), bson.D{{}})

	if err != nil {
		errString := "error fetching clients"
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
