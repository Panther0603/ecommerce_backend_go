package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBset()

func DBset() *mongo.Client {

	// one way of the connect
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// client, err := mongo.Connect(context.Background(), clientOptions)

	// another way to connect

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("Error occured while creating dataabse")
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Println("Error occured while cheking ping dataabse")
		log.Fatal(err)
		return nil
	}
	log.Println("sucessfully connected")
	return client

}

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {

	collection := client.Database("ecommerce").Collection(collectionName)
	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("ecommerce").Collection(collectionName)
	return collection
}
