package mongo

import (
	"fmt"
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/schemas"
	simplejson "github.com/bitly/go-simplejson"
	"gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
)

var _ = check.Suite(&S{})

type S struct {
	Db *MongoDB
}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) SetUpSuite(c *check.C) {
	var err error

	os.Setenv("MONGO_URI", "localhost:27017/backstage_beat_test")
	s.Db, err = New()
	c.Assert(err, check.IsNil)

	session := s.Db.session.Clone()
	defer session.Close()
	session.DB("").DropDatabase()
}

func (s *S) TestNewMongoDBConfigWithEnviromentVariables(c *check.C) {
	os.Unsetenv("MONGO_URI")
	os.Unsetenv("MONGO_USER")
	os.Unsetenv("MONGO_PASSWORD")

	db, err := New()
	c.Assert(err, check.IsNil)
	c.Assert(db, check.Not(check.IsNil))
	c.Assert(db.dialInfo.Addrs, check.DeepEquals, []string{"localhost:27017"})
	c.Assert(db.dialInfo.Database, check.Equals, "backstage_beat_local")
	c.Assert(db.dialInfo.Username, check.Equals, "")
	c.Assert(db.dialInfo.Password, check.Equals, "")
}

func (s *S) TestNewMongoDBConfigWithDefaultVariables(c *check.C) {
	c.Assert(s.Db.dialInfo.Addrs, check.DeepEquals, []string{"localhost:27017"})
	c.Assert(s.Db.dialInfo.Database, check.Equals, "backstage_beat_test")
}

func (s *S) TestGetFromRegister(c *check.C) {
	db, err := db.New("mongo")
	c.Assert(err, check.IsNil)
	c.Assert(db, check.FitsTypeOf, &MongoDB{})
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

func (s *S) TestFindItemSchema(c *check.C) {
	for i := 0; i < 3; i++ {
		dbErr := s.Db.CreateItemSchema(&schemas.ItemSchema{CollectionName: fmt.Sprintf("find-%d", i)})
		c.Assert(dbErr, check.IsNil)
		i++
	}
	filter, err := db.NewFilterFromQueryString("")
	c.Assert(err, check.IsNil)

	reply, dbErr := s.Db.FindItemSchema(filter)
	c.Assert(dbErr, check.IsNil)
	c.Assert(len(reply.Items) > 3, check.Equals, true)
}

func (s *S) TestFindItemSchemaWithExactPattern(c *check.C) {
	for i := 0; i < 3; i++ {
		dbErr := s.Db.CreateItemSchema(&schemas.ItemSchema{CollectionName: fmt.Sprintf("find-exact-%d", i)})
		c.Assert(dbErr, check.IsNil)
	}
	filter, err := db.NewFilterFromQueryString("filter[where][collectionName]=find-exact-1")
	c.Assert(err, check.IsNil)

	reply, dbErr := s.Db.FindItemSchema(filter)
	c.Assert(dbErr, check.IsNil)
	c.Assert(len(reply.Items), check.Equals, 1)
	c.Assert(reply.Items[0].CollectionName, check.Equals, "find-exact-1")
}

func (s *S) TestFindOneItemSchemaWithExactPattern(c *check.C) {
	for i := 0; i < 3; i++ {
		dbErr := s.Db.CreateItemSchema(&schemas.ItemSchema{CollectionName: fmt.Sprintf("find-one-exact-%d", i)})
		c.Assert(dbErr, check.IsNil)
	}
	filter, err := db.NewFilterFromQueryString("filter[where][collectionName]=find-one-exact-1")
	c.Assert(err, check.IsNil)

	itemSchema, dbErr := s.Db.FindOneItemSchema(filter)
	c.Assert(dbErr, check.IsNil)
	c.Assert(itemSchema.CollectionName, check.Equals, "find-one-exact-1")
}

func (s *S) TestFindOneItemSchemaWithNotFound(c *check.C) {
	filter, err := db.NewFilterFromQueryString("filter[where][collectionName]=not-found")
	c.Assert(err, check.IsNil)

	_, dbErr := s.Db.FindOneItemSchema(filter)
	c.Assert(dbErr, check.NotNil)
	c.Assert(dbErr.StatusCode(), check.Equals, 404)
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

func (s *S) TestMongoBuildWhereSimple(c *check.C) {
	where, _ := simplejson.NewJson([]byte(`{"name": "r2"}`))
	mongoWhere := BuildMongoWhere(where, "id")
	c.Assert(mongoWhere, check.DeepEquals, bson.M{"name": "r2"})
}

func (s *S) TestMongoBuildWhereAndQuery(c *check.C) {
	where, _ := simplejson.NewJson([]byte(`{"and": [{"name": "wilson"}, {"tenantId": "globocom"}]}`))
	mongoWhere := BuildMongoWhere(where, "id")
	c.Assert(mongoWhere, check.DeepEquals, bson.M{
		"$and": []bson.M{
			bson.M{"name": "wilson"},
			bson.M{"tenantId": "globocom"},
		},
	})

}

func (s *S) TestMongoBuildWhereWithPrimaryKey(c *check.C) {
	where, _ := simplejson.NewJson([]byte(`{"tenantId": "globocom"}`))
	mongoWhere := BuildMongoWhere(where, "tenantId")
	c.Assert(mongoWhere, check.DeepEquals, bson.M{"_id": "globocom"})
}
