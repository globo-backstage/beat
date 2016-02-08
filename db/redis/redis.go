package redis

import (
	"encoding/json"
	"fmt"
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var (
	DbPrefix          = "db"
	ErrNotImplemented = errors.New("Not Implemented for Redis", http.StatusNotImplemented)
)

type Redis struct {
	host     string
	password string
	db       int
	pool     *redis.Pool
}

func init() {
	viper.SetDefault("redis.host", "localhost:6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool.maxIdle", 10)
	viper.SetDefault("redis.pool.maxActive", 10)
	viper.SetDefault("redis.pool.wait", true)
	viper.SetDefault("redis.pool.idleTimeout", 180e9)
}

func New() (*Redis, error) {
	d := &Redis{}
	d.host = viper.GetString("redis.host")
	d.password = viper.GetString("redis.password")
	d.db = viper.GetInt("redis.db")

	d.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", d.host)
			if err != nil {
				return nil, err
			}
			if d.password != "" {
				_, err = conn.Do("AUTH", d.password)
				if err != nil {
					return nil, err
				}
			}
			_, err = conn.Do("SELECT", d.db)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
		MaxIdle:     viper.GetInt("redis.pool.maxIdle"),
		MaxActive:   viper.GetInt("redis.pool.maxActive"),
		Wait:        viper.GetBool("redis.pool.wait"),
		IdleTimeout: time.Duration(viper.GetInt("redis.pool.idleTimeout")),
	}

	return d, nil
}

func (r *Redis) CreateItemSchema(itemSchema *schemas.ItemSchema) errors.Error {
	return r.createResource(schemas.ItemSchemaCollectionName, itemSchema.CollectionName, itemSchema)
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
	if err == redis.ErrNil {
		return nil, db.ItemSchemaNotFound
	} else if err != nil {
		return nil, errors.Wraps(err, http.StatusInternalServerError)
	}
	return itemSchema, nil
}

func (r *Redis) DeleteItemSchemaByCollectionName(collectionName string) errors.Error {
	err := r.deleteResource(schemas.ItemSchemaCollectionName, collectionName)
	if err == redis.ErrNil {
		return db.ItemSchemaNotFound
	} else if err != nil {
		return errors.Wraps(err, http.StatusInternalServerError)
	}

	return nil
}

func (r *Redis) Flush() {
	redisConn := r.pool.Get()
	redisConn.Do("flushdb")
	redisConn.Close()
}

func (r *Redis) createResource(collectionName string, primaryKey string, result interface{}) errors.Error {
	buf, err := json.Marshal(result)
	if err != nil {
		return errors.Wraps(err, http.StatusBadRequest)
	}
	redisKey := r.key(collectionName, primaryKey)
	redisConn := r.pool.Get()
	_, err = redis.String(redisConn.Do("SET", redisKey, string(buf), "NX"))
	redisConn.Close()

	if err == redis.ErrNil {
		validationError := &errors.ValidationError{}
		validationError.Put("_all", "Duplicated resource")
		return validationError
	} else if err != nil {
		return errors.Wraps(err, http.StatusInternalServerError)
	}

	return nil
}

func (r *Redis) getResource(collectionName string, primaryKey string, result interface{}) error {
	redisKey := r.key(collectionName, primaryKey)
	redisConn := r.pool.Get()
	reply, err := redis.String(redisConn.Do("GET", redisKey))
	redisConn.Close()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(reply), result)
}

func (r *Redis) deleteResource(collectionName string, primaryKey string) error {
	redisKey := r.key(collectionName, primaryKey)
	redisConn := r.pool.Get()
	reply, err := redis.Int(redisConn.Do("DEL", redisKey))
	redisConn.Close()

	if err != nil {
		return err
	}

	if reply == 0 {
		return redis.ErrNil
	}

	return nil
}

func (r *Redis) key(collectionName string, primaryKey string) string {
	return fmt.Sprintf("%s:%s:%s", DbPrefix, collectionName, primaryKey)
}
