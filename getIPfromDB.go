package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func getIPfromDB(user, realm string) string {

	//filter := bson.D{{"user_name", user}}
	filter := bson.M{"user_name": user}

	client := connectdb()
	defer client.Disconnect(context.TODO())
	coll := client.Database("ldapDB").Collection("user_history")

	result, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var data []LoginHistory
	result.All(context.TODO(), &data)

	if len(data) >= 1 {
		log.Println(data[0])
		log.Printf("[DB]Get FramedIp from DB for user: %v , IP: %v\n", user, data[0].FramedIP)
		return data[0].FramedIP
	} else {
		return "Not available"
	}
}
