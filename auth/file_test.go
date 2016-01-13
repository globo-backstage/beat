package auth

import (
	"github.com/backstage/beat/config"
	"gopkg.in/check.v1"
	"net/http"
	"testing"
)

var _ = check.Suite(&S{})

type S struct {
	authenticaton *FileAuthentication
}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) SetUpSuite(c *check.C) {
	err := config.ReadConfigFile("../examples/config.yml")
	c.Assert(err, check.IsNil)

	s.authenticaton = NewFileAuthentication()
}

func (s *S) TestFileAuthenticationWithUserFound(c *check.C) {
	header := &http.Header{}
	header.Set("Token", "example1")

	user := s.authenticaton.GetUser(header)
	c.Assert(user, check.NotNil)
	c.Assert(user.Email(), check.Equals, "admin@example.net")

	header = &http.Header{}
	header.Set("Token", "example2")

	user = s.authenticaton.GetUser(header)
	c.Assert(user, check.Not(check.IsNil))
	c.Assert(user.Email(), check.Equals, "guest@example.net")
}

func (s *S) TestFileAuthenticationWithUserNotFound(c *check.C) {
	header := &http.Header{}
	header.Set("Token", "not-found")

	user := s.authenticaton.GetUser(header)
	c.Assert(user, check.IsNil)
}

func (s *S) TestFileAuthenticationWithMissingToken(c *check.C) {
	header := &http.Header{}
	user := s.authenticaton.GetUser(header)
	c.Assert(user, check.IsNil)
}
