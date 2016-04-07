package db

import (
	"fmt"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
	// simplejson "github.com/bitly/go-simplejson"
	"net/http"
	"strings"
)

var (
	ItemSchemaNotFound       = errors.New("item-schema not found", http.StatusNotFound)
	CollectionSchemaNotFound = errors.New("collection-schema not found", http.StatusNotFound)
	databases                = map[string]RegisterFunc{}
)

type RegisterFunc func() (Database, error)

type Database interface {
	CreateItemSchema(*schemas.ItemSchema) errors.Error
	FindItemSchema(*Filter) (*ItemSchemasReply, errors.Error)
	FindOneItemSchema(*Filter) (*schemas.ItemSchema, errors.Error)
	FindItemSchemaByCollectionName(string) (*schemas.ItemSchema, errors.Error)
	DeleteItemSchemaByCollectionName(string) errors.Error

	CreateResource(string, *schemas.CollectionSchema) errors.Error
}

type ItemSchemasReply struct {
	Items []*schemas.ItemSchema `json:"items"`
}

// Register inserts a implementation of `Database` in the register, is useful
// to auto discover implementations and change it without changing the code.
func Register(name string, fn RegisterFunc) {
	databases[name] = fn
}

// New returns a implementation of `Database` found in the register, if not found
// return an error.
func New(name string) (Database, error) {
	fn := databases[name]
	if fn == nil {
		return nil, ErrNotFound{name: name}
	}
	db, err := fn()

	if err != nil {
		return nil, databaseError{name: name, originalErr: err}
	}

	return db, nil
}

type ErrNotFound struct {
	name string
}

func (d ErrNotFound) Error() string {
	availableDatabases := make([]string, 0, len(databases))
	for db := range databases {
		availableDatabases = append(availableDatabases, db)
	}

	return fmt.Sprintf(`Database "%s" not found, are available: %s.`, d.name, strings.Join(availableDatabases, ", "))
}

type databaseError struct {
	name        string
	originalErr error
}

func (d databaseError) Error() string {
	return fmt.Sprintf(`[db][%s] %s`, d.name, d.originalErr.Error())
}
