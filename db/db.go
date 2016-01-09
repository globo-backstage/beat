package db

import (
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
)

type Database interface {
	CreateItemSchema(*schemas.ItemSchema) errors.Error
	FindItemSchema(filter *Filter) (interface{}, errors.Error)
}
