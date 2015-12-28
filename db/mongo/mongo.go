package mongo

import (
	"github.com/backstage/beat/schemas"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/mgo.v2"
)

type MongoConfig struct {
	Uri      string `default:"localhost:27017/backstage_beat_local"`
	User     string
	Password string
}

type MongoDB struct {
	config  MongoConfig
	session *mgo.Session
}

func New() (*MongoDB, error) {
	d := &MongoDB{}
	err := envconfig.Process("mongo", &d.config)

	if err != nil {
		return nil, err
	}

	dialInfo, err := mgo.ParseURL(d.config.Uri)
	if err != nil {
		return nil, err
	}

	if d.config.User != "" {
		dialInfo.Username = d.config.User
	}

	if d.config.Password != "" {
		dialInfo.Password = d.config.Password
	}

	dialInfo.FailFast = true
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		return nil, err
	}
	d.session = session
	return d, nil
}

func (m *MongoDB) CreateItemSchema(itemSchema *schemas.ItemSchema) error {
	session := m.session.Clone()
	defer session.Close()
	return session.DB("").C(schemas.ItemSchemaCollectionName).Insert(itemSchema)
}
