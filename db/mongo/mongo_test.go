package mongo

import (
	simplejson "github.com/bitly/go-simplejson"
	"gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
)

var _ = check.Suite(&S{})

type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) TestNewMongoDBConfig(c *check.C) {
	os.Unsetenv("MONGO_URI")
	os.Unsetenv("MONGO_USER")
	os.Unsetenv("MONGO_PASSWORD")

	db, err := New()
	c.Assert(err, check.IsNil)
	c.Assert(db, check.Not(check.IsNil))
	c.Assert(db.config.Uri, check.Equals, "localhost:27017/backstage_beat_local")
	c.Assert(db.config.User, check.Equals, "")
	c.Assert(db.config.Password, check.Equals, "")

	os.Setenv("MONGO_URI", "localhost:27017/backstage_beat_test")

	db, err = New()
	c.Assert(err, check.IsNil)
	c.Assert(db, check.Not(check.IsNil))
	c.Assert(db.config.Uri, check.Equals, "localhost:27017/backstage_beat_test")
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
