package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Client struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Coaches     []*Coach           `json:"coach`
	FirstName   string             `json:"firstname"`
	LastName    string             `json:"lastname"`
	Email       string             `json:"email"`
	PhoneNumber int32              `json:"phonenumber"`
	Password    string             `json:"password"`
}
type Coach struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName   string             `json:"firstname"`
	LastName    string             `json:"lastname"`
	Email       string             `json:"email"`
	PhoneNumber int32              `json:"phonenumber"`
	Password    string             `json:"password"`
}

type Exercise struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}
