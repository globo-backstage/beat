package mongo

import (
	"gopkg.in/check.v1"
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
