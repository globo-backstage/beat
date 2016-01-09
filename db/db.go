package db

import (
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
)

type Database interface {
	CreateItemSchema(*schemas.ItemSchema) errors.Error
	FindItemSchema(filter *Filter) (*ItemSchemasReply, errors.Error)
}

type ItemSchemasReply struct {
	Items []schemas.ItemSchema `json:"items"`
}
