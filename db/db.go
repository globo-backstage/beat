package db

import (
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
)

type Database interface {
	CreateItemSchema(*schemas.ItemSchema) errors.Error
	FindItemSchema(*Filter) (*ItemSchemasReply, errors.Error)
	FindOneItemSchema(*Filter) (*schemas.ItemSchema, errors.Error)
	FindItemSchemaByCollectionName(string) (*schemas.ItemSchema, errors.Error)
	DeleteItemSchemaByCollectionName(string) errors.Error
}

type ItemSchemasReply struct {
	Items []schemas.ItemSchema `json:"items"`
}
