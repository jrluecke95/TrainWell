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
