package mongo

import (
	"github.com/kelseyhightower/envconfig"
)

type MongoConfig struct {
	Uri      string `default:"localhost:27017/backstage_beat_local"`
	User     string
	Password string
}

type MongoDB struct {
	config MongoConfig
}

func New() (*MongoDB, error) {
	d := &MongoDB{}
	err := envconfig.Process("mongo", &d.config)
	return d, err
}
