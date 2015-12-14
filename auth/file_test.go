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

func (s *S) TestFileAuthenticationWith(c *check.C) {
	a, err := NewFileAuthentication("../examples/tokens.yml")
	c.Assert(err, check.IsNil)
	
	header := &http.Header{}
	header.Set("Token", "example1")
	

	user := a.GetUser(header)
	c.Assert(user, check.Not(check.IsNil))
	
	c.Assert(user.Email(), check.Equals, "admin@example.net")
}
