package mongo

import (
	"fmt"
	_ "github.com/backstage/beat/config"
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func init() {
	viper.SetDefault("mongo.uri", "localhost:27017/backstage_beat_local")
	viper.SetDefault("mongo.failFast", true)

	db.Register("mongo", func() (db.Database, error) {
		return New()
	})
}

type MongoDB struct {
	dialInfo *mgo.DialInfo
	session  *mgo.Session
}

func New() (*MongoDB, error) {
	d := &MongoDB{}

	var err error
	d.dialInfo, err = mgo.ParseURL(viper.GetString("mongo.uri"))
	if err != nil {
		return nil, err
	}

	d.dialInfo.Username = viper.GetString("mongo.user")
	d.dialInfo.Password = viper.GetString("mongo.password")

	d.dialInfo.FailFast = viper.GetBool("mongo.failFast")
	session, err := mgo.DialWithInfo(d.dialInfo)

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

	if err != nil {
		return convertMongoError(err)
	}

	return nil
}

func (m *MongoDB) CreateResource(collectionName string, item *schemas.CollectionSchema) errors.Error {
	session := m.session.Clone()
	defer session.Close()
	// item.Set("_id", collectionName)
	// item.Del("collectionName")
	// bs = bson.M(item)

	fmt.Println(item)
	if err := session.DB("").C(collectionName).Insert(item); err != nil {
		return convertMongoError(err)
	}
	return nil
}

func (m *MongoDB) FindItemSchema(filter *db.Filter) (*db.ItemSchemasReply, errors.Error) {
	session := m.session.Clone()
	defer session.Close()
	where := BuildMongoWhere(filter.Where, schemas.ItemSchemaPrimaryKey)
	query := session.DB("").C(schemas.ItemSchemaCollectionName).Find(where)

	reply := &db.ItemSchemasReply{}
	reply.Items = []*schemas.ItemSchema{}
	err := query.Skip(filter.Skip()).Limit(filter.PerPage).Iter().All(&reply.Items)

	if err != nil {
		return nil, errors.Wraps(err, http.StatusInternalServerError)
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
		return nil, db.ItemSchemaNotFound
	} else if err != nil {
		return nil, errors.Wraps(err, http.StatusInternalServerError)
	}

	return itemSchema, nil
}

func (m *MongoDB) FindItemSchemaByCollectionName(collectionName string) (*schemas.ItemSchema, errors.Error) {
	session := m.session.Clone()
	defer session.Close()

	itemSchema := &schemas.ItemSchema{}
	err := session.DB("").C(schemas.ItemSchemaCollectionName).FindId(collectionName).One(&itemSchema)

	if err == mgo.ErrNotFound {
		return nil, db.ItemSchemaNotFound
	} else if err != nil {
		return nil, errors.Wraps(err, http.StatusInternalServerError)
	}

	return itemSchema, nil
}

func (m *MongoDB) DeleteItemSchemaByCollectionName(collectionName string) errors.Error {
	session := m.session.Clone()
	defer session.Close()

	err := session.DB("").C(schemas.ItemSchemaCollectionName).RemoveId(collectionName)
	if err == mgo.ErrNotFound {
		return db.ItemSchemaNotFound
	} else if err != nil {
		return errors.Wraps(err, http.StatusInternalServerError)
	}
	return nil
}

// TODO: change to FindResources
func (m *MongoDB) FindCollectionSchema(collectionName string, filter *db.Filter) (*db.ResourceReply, errors.Error) {
	session := m.session.Clone()
	defer session.Close()
	query := session.DB("").C(collectionName).Find(nil)

	reply := &db.ResourceReply{}
	reply.Items = []*schemas.Resources{}

	println("===============")
	println("Geting from: ", collectionName)
	println("===============")
	err := query.Limit(10).Iter().All(&reply.Items)

	if err != nil {
		return nil, errors.Wraps(err, http.StatusInternalServerError)
	}

	return reply, nil
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

func convertMongoError(err error) errors.Error {
	if mongoErr, ok := err.(*mgo.LastError); ok {
		if mongoErr.Code == 11000 {
			return buildMongoDuplicatedError()
		}
	}
	return errors.Wraps(err, http.StatusInternalServerError)
}

func buildMongoDuplicatedError() errors.Error {
	validationError := &errors.ValidationError{}
	validationError.Put("_all", "Duplicated resource")
	return validationError
}
