package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func writedb(user, realm, FramedIPstring string) {

	loc := time.FixedZone("KST", +9*60*60)
	currentTime := time.Now().In(loc)

	filter := bson.M{"user_name": user}
	login := LoginHistory{
		Realm:     realm,
		UserName:  user,
		LastLogin: currentTime,
		Enabled:   "True",
		FramedIP:  FramedIPstring,
	}
	var update bson.M
	switch rdDown {
	case false:
		// update every entry
		update = bson.M{"$set": login}
	case true:
		// update except FramedIP
		update = bson.M{"$set": bson.M{"last_login": currentTime, "realm": realm, "enabled": "True"}}
	}
	opts := options.Update().SetUpsert(true)

	client := connectdb()
	defer client.Disconnect(context.TODO())
	coll := client.Database("ldapDB").Collection("user_history")

	result, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}

	log.Printf("[DB]Number of documents updated: %v\n", result.ModifiedCount)
	log.Printf("[DB]Number of documents upserted: %v\n", result.UpsertedCount)

	if !secondary {
		if client := connectdbS(); client != nil {
			defer client.Disconnect(context.TODO())
			coll := client.Database("ldapDB").Collection("user_history")

			result, err := coll.UpdateOne(context.TODO(), filter, update, opts)
			if err != nil {
				panic(err)
			}

			log.Printf("[DB-S]Number of documents updated: %v\n", result.ModifiedCount)
			log.Printf("[DB-S]Number of documents upserted: %v\n", result.UpsertedCount)
		}
	}
}
