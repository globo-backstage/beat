package main

import (
	"flag"
	"log"

	"github.com/backstage/beat/auth"
	"github.com/backstage/beat/db/mongo"
	"github.com/backstage/beat/server"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "./examples/config.yml", "Config file")
	flag.Parse()

	authentication, err := auth.NewFileAuthentication(configFile)

	if err != nil {
		log.Fatal(err)
	}

	db, err := mongo.New()

	if err != nil {
		log.Fatalf("[mongodb] %s", err)
	}

	s := server.New(authentication, db)
	s.Run()
}
