package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var client *mongo.Client

func CreatePersonEndpoint(r http.ResponseWriter, req *http.Request) {
	r.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(req.Body).Decode(&person)
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(r).Encode(result)
}

func GetPeopleEndpoint(r http.ResponseWriter, req *http.Request) {
	r.Header().Add("content-type", "application/json")
	var people []Person
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		r.Write([]byte(`"message":` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}

	if err := cursor.Err(); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		r.Write([]byte(`"message":` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(r).Encode(people)

}

func GetPersonEndpoint(r http.ResponseWriter, req *http.Request) {

	r.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var person Person
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		r.Write([]byte(`"message":` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(r).Encode(person)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
