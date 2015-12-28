package db

import (
	"github.com/backstage/beat/schemas"
)

type Database interface {
	CreateItemSchema(*schemas.ItemSchema) error
}
