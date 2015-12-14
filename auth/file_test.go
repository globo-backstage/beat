package auth


import (
	"testing"
	"gopkg.in/check.v1"
	"net/http"
)

var _ = check.Suite(&S{})
type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) TestFileAuthenticationWithUserFound(c *check.C) {
	a, err := NewFileAuthentication("../examples/config.yml")
	c.Assert(err, check.IsNil)

	header := &http.Header{}
	header.Set("Token", "example1")

	user := a.GetUser(header)
	c.Assert(user, check.Not(check.IsNil))
	c.Assert(user.Email(), check.Equals, "admin@example.net")

	header = &http.Header{}
	header.Set("Token", "example2")

	user = a.GetUser(header)
	c.Assert(user, check.Not(check.IsNil))
	c.Assert(user.Email(), check.Equals, "guest@example.net")
}

func (s *S) TestFileAuthenticationWithUserNotFound(c *check.C) {
	a, err := NewFileAuthentication("../examples/config.yml")
	c.Assert(err, check.IsNil)

	header := &http.Header{}
	header.Set("Token", "not-found")

	user := a.GetUser(header)
	c.Assert(user, check.IsNil)
}
