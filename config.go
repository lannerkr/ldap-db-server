package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	LdapServer        string
	LdapBaseDn        string
	RadServer         string
	RadSecret         string
	RadNASidS         string
	RadNASidP         string
	RadNASidE         string
	RadClientPassword string
	LogPath           string
	RadHCuser         string
	RadHCrealm        string
	Secondary         bool
	MongoDBM          string
	MongoDBS          string
}

var configuration Configuration

func config(conf string) {
	file, _ := os.Open(conf)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println("error:", err)
	}
	log.Println(configuration)
}
