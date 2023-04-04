package util

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://sainuthanreddy264:Woovb6N7LCbtewq5@cluster0.jfcmmg5.mongodb.net/?retryWrites=true&w=majority"

const colName = "Users"
const DbName = "GoDB"

var Db *mongo.Collection

func ConnectDb() {

	// Replace the placeholder with your Atlas connection string

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	Db = client.Database(DbName).Collection(colName)

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
