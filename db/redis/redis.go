package redis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backstage/beat/db"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
	"github.com/spf13/viper"
	"gopkg.in/redis.v4"
)

var (
	DbPrefix          = "db"
	ErrNotImplemented = errors.New("Not Implemented for Redis", http.StatusNotImplemented)
)

type Redis struct {
	*redis.Client
}

func init() {
	viper.SetDefault("redis.host", "localhost:6379")
	viper.SetDefault("redis.db", 0)

	db.Register("redis", func() (db.Database, error) {
		return New()
	})
}

func New() (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	return &Redis{client}, nil
}

func (r *Redis) CreateItemSchema(itemSchema *schemas.ItemSchema) errors.Error {
	return r.createResource(schemas.ItemSchemaCollectionName, itemSchema.CollectionName, itemSchema)
}

func (r *Redis) UpdateItemSchema(*schemas.ItemSchema) errors.Error {
	return ErrNotImplemented
}

func (r *Redis) FindItemSchema(*db.Filter) (*db.ItemSchemasReply, errors.Error) {
	return nil, ErrNotImplemented
}

func (r *Redis) FindOneItemSchema(*db.Filter) (*schemas.ItemSchema, errors.Error) {
	return nil, ErrNotImplemented
}

func (r *Redis) FindItemSchemaByCollectionName(collectionName string) (*schemas.ItemSchema, errors.Error) {
	itemSchema := &schemas.ItemSchema{}
	err := r.getResource(schemas.ItemSchemaCollectionName, collectionName, itemSchema)
	if err == redis.Nil {
		return nil, db.ErrItemSchemaNotFound
	} else if err != nil {
		return nil, errors.Wraps(err, http.StatusInternalServerError)
	}
	return itemSchema, nil
}

func (r *Redis) DeleteItemSchema(collectionName string) errors.Error {
	err := r.deleteResource(schemas.ItemSchemaCollectionName, collectionName)
	if err == redis.Nil {
		return db.ErrItemSchemaNotFound
	} else if err != nil {
		return errors.Wraps(err, http.StatusInternalServerError)
	}

	return nil
}

func (r *Redis) createResource(collectionName string, primaryKey string, result interface{}) errors.Error {
	buf, err := json.Marshal(result)
	if err != nil {
		return errors.Wraps(err, http.StatusBadRequest)
	}
	created, err := r.SetNX(r.key(collectionName, primaryKey), string(buf), 0).Result()

	if err == redis.Nil || !created {
		validationError := &errors.ValidationError{}
		validationError.Put("_all", "Duplicated resource")
		return validationError
	} else if err != nil {
		return errors.Wraps(err, http.StatusInternalServerError)
	}

	return nil
}

func (r *Redis) getResource(collectionName string, primaryKey string, result interface{}) error {
	reply, err := r.Get(r.key(collectionName, primaryKey)).Bytes()

	if err != nil {
		return err
	}

	return json.Unmarshal(reply, result)
}

func (r *Redis) deleteResource(collectionName string, primaryKey string) error {
	reply, err := r.Del(r.key(collectionName, primaryKey)).Result()

	if err != nil {
		return err
	}

	if reply == 0 {
		return redis.Nil
	}

	return nil
}

func (r *Redis) key(collectionName string, primaryKey string) string {
	return fmt.Sprintf("%s:%s:%s", DbPrefix, collectionName, primaryKey)
}
