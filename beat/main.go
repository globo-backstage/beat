package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"

	"github.com/backstage/beat/config"
	"github.com/backstage/beat/server"

	_ "github.com/backstage/beat/auth/static"
	_ "github.com/backstage/beat/db/mongo"
	_ "github.com/backstage/beat/db/redis"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "./examples/config.yml", "Config file")
	flag.Parse()

	err := config.ReadConfigFile(configFile)

	if err != nil {
		log.Fatal(err)
	}

	config.LoadLogSettings()

	s, err := server.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	s.Run()
}
