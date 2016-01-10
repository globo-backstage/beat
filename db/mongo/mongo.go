package mongo

import (
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func (m *MongoDB) CreateItemSchema(itemSchema *schemas.ItemSchema) errors.Error {
	session := m.session.Clone()
	defer session.Close()
	err := session.DB("").C(schemas.ItemSchemaCollectionName).Insert(itemSchema)

	if err == nil {
		return nil
	} else {
		return errors.Wraps(err, 500)
	}
}

func (m *MongoDB) FindItemSchema(filter *db.Filter) (*db.ItemSchemasReply, errors.Error) {
	session := m.session.Clone()
	defer session.Close()
	where := BuildMongoWhere(filter.Where, schemas.ItemSchemaPrimaryKey)
	query := session.DB("").C(schemas.ItemSchemaCollectionName).Find(where)

	reply := &db.ItemSchemasReply{}
	reply.Items = []schemas.ItemSchema{}
	err := query.Skip(filter.Skip()).Limit(filter.PerPage).Iter().All(&reply.Items)

	if err != nil {
		return nil, errors.Wraps(err, 500)
	}

	return reply, nil
}

func (m *MongoDB) FindOneItemSchema(filter *db.Filter) (*schemas.ItemSchema, errors.Error) {
	session := m.session.Clone()
	defer session.Close()
	where := BuildMongoWhere(filter.Where, schemas.ItemSchemaPrimaryKey)
	query := session.DB("").C(schemas.ItemSchemaCollectionName).Find(where)

	itemSchema := &schemas.ItemSchema{}
	err := query.One(&itemSchema)

	if err == mgo.ErrNotFound {
		return nil, errors.New("item-schema not found", 404)
	} else if err != nil {
		return nil, errors.Wraps(err, 500)
	}

	return itemSchema, nil
}

func (m *MongoDB) FindItemSchemaByCollectionName(collectionName string) (*schemas.ItemSchema, errors.Error) {
	session := m.session.Clone()
	defer session.Close()

	itemSchema := &schemas.ItemSchema{}
	err := session.DB("").C(schemas.ItemSchemaCollectionName).FindId(collectionName).One(&itemSchema)

	if err == mgo.ErrNotFound {
		return nil, errors.Newf(404, `item-schema "%s" not found`, collectionName)
	} else if err != nil {
		return nil, errors.Wraps(err, 500)
	}

	return itemSchema, nil
}

func (m *MongoDB) DeleteItemSchemaByCollectionName(collectionName string) errors.Error {
	session := m.session.Clone()
	defer session.Close()

	err := session.DB("").C(schemas.ItemSchemaCollectionName).RemoveId(collectionName)
	if err == mgo.ErrNotFound {
		return errors.Newf(404, `item-schema "%s" not found`, collectionName)
	} else if err != nil {
		return errors.Wraps(err, 500)
	}
	return nil
}

func BuildMongoWhere(where *simplejson.Json, primaryKey string) bson.M {
	mongoWhere := bson.M{}
	for key, value := range where.MustMap() {
		switch key {
		case "and", "or", "nor":
			mongoWhere["$"+key] = buildMongoWhereByArray(
				where.Get(key),
				primaryKey,
			)
			continue

		case primaryKey:
			mongoWhere["_id"] = value
			continue
		}
		mongoWhere[key] = value
	}
	return mongoWhere
}

func buildMongoWhereByArray(wheres *simplejson.Json, primaryKey string) []bson.M {
	mongoWheres := []bson.M{}
	for key, _ := range wheres.MustArray() {
		mongoWhere := BuildMongoWhere(wheres.GetIndex(key), primaryKey)
		mongoWheres = append(mongoWheres, mongoWhere)
	}
	return mongoWheres
}
