package redis

import (
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/schemas"
	"gopkg.in/check.v1"
	"net/http"
	"testing"
)

var _ = check.Suite(&S{})

type S struct {
	Db *Redis
}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) SetUpSuite(c *check.C) {
	var err error

	s.Db, err = New()
	s.Db.Flush()
	c.Assert(err, check.IsNil)
}

func (s *S) TestImplementInterface(c *check.C) {
	var dbType db.Database
	c.Assert(s.Db, check.Implements, &dbType)
}

func (s *S) TestCreateItemSchema(c *check.C) {
	itemSchema := &schemas.ItemSchema{CollectionName: "test-schema"}
	dbErr := s.Db.CreateItemSchema(itemSchema)
	c.Assert(dbErr, check.IsNil)

	itemSchema, dbErr = s.Db.FindItemSchemaByCollectionName("test-schema")
	c.Assert(dbErr, check.IsNil)
	c.Assert(itemSchema.CollectionName, check.Equals, "test-schema")
}

func (s *S) TestCreateItemSchemaDuplicated(c *check.C) {
	itemSchema := &schemas.ItemSchema{CollectionName: "duplicated-schema"}
	dbErr := s.Db.CreateItemSchema(itemSchema)

	c.Assert(dbErr, check.IsNil)

	dbErr = s.Db.CreateItemSchema(itemSchema)
	c.Assert(dbErr, check.NotNil)
	c.Assert(dbErr.StatusCode(), check.Equals, 422)
	c.Assert(dbErr.Error(), check.Equals, "_all: Duplicated resource")
}

func (s *S) TestFindItemSchemaByCollectionNameWithNotFound(c *check.C) {
	_, dbErr := s.Db.FindItemSchemaByCollectionName("not-found")
	c.Assert(dbErr, check.NotNil)
	c.Assert(dbErr.StatusCode(), check.Equals, 404)
}

func (s *S) TestDeleteItemSchemaByCollectionNameWithNotFound(c *check.C) {
	dbErr := s.Db.DeleteItemSchemaByCollectionName("not-found")
	c.Assert(dbErr, check.NotNil)
	c.Assert(dbErr.StatusCode(), check.Equals, 404)
}

func (s *S) TestDeleteItemSchemaByCollectionName(c *check.C) {
	dbErr := s.Db.CreateItemSchema(&schemas.ItemSchema{CollectionName: "to-be-deleted"})
	c.Assert(dbErr, check.IsNil)

	dbErr = s.Db.DeleteItemSchemaByCollectionName("to-be-deleted")
	c.Assert(dbErr, check.IsNil)

	_, dbErr = s.Db.FindItemSchemaByCollectionName("to-be-deleted")
	c.Assert(dbErr, check.NotNil)
	c.Assert(dbErr.StatusCode(), check.Equals, 404)
}

func (s *S) TestFindItemSchema(c *check.C) {
	reply, dbErr := s.Db.FindItemSchema(nil)

	c.Assert(reply, check.IsNil)
	c.Assert(dbErr, check.NotNil)

	c.Assert(dbErr.StatusCode(), check.Equals, http.StatusNotImplemented)
	c.Assert(dbErr.Error(), check.Equals, "Not Implemented for Redis")
}

func (s *S) TestFindOneItemSchema(c *check.C) {
	reply, dbErr := s.Db.FindOneItemSchema(nil)

	c.Assert(reply, check.IsNil)
	c.Assert(dbErr, check.NotNil)

	c.Assert(dbErr.StatusCode(), check.Equals, http.StatusNotImplemented)
	c.Assert(dbErr.Error(), check.Equals, "Not Implemented for Redis")
}