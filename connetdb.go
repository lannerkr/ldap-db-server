package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connectdb() (client *mongo.Client) {

	// Create a new client and connect to the server
	uri := "mongodb://" + configuration.MongoDBM + "/?retryWrites=true&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("[DB]Successfully connected and pinged.")

	return client
}

func connectdbS() (client *mongo.Client) {

	// Create a new client and connect to the server
	uri := "mongodb://" + configuration.MongoDBS + "/?retryWrites=true&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("[DB-S] Failed to be Successfully connected.")
		return nil
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println("[DB-S] Failed to be Successfully pinged.")
		return nil
	}
	log.Println("[DB-S]Successfully connected and pinged.")

	return client
}
