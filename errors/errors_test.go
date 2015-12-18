package errors

import (
	"encoding/json"
	originalErrors "errors"
	simplejson "github.com/bitly/go-simplejson"
	"gopkg.in/check.v1"
	"testing"
)

var _ = check.Suite(&S{})

type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) TestWrapsNewError(c *check.C) {
	err := Wraps(originalErrors.New("test error 123"), 503)
	c.Assert(err, check.Not(check.IsNil))
	c.Assert(err.Error(), check.Equals, "test error 123")
	c.Assert(err.StatusCode(), check.Equals, 503)
}

func (s *S) TestNewError(c *check.C) {
	err := New("test error 123", 503)
	c.Assert(err, check.Not(check.IsNil))
	c.Assert(err.Error(), check.Equals, "test error 123")
	c.Assert(err.StatusCode(), check.Equals, 503)
}

func (s *S) TestNewfError(c *check.C) {
	err := Newf(503, "test error %s", "123")
	c.Assert(err, check.Not(check.IsNil))
	c.Assert(err.Error(), check.Equals, "test error 123")
	c.Assert(err.StatusCode(), check.Equals, 503)
}

func (s *S) TestMarshallJSONWrappedError(c *check.C) {
	errWrapped := Wraps(originalErrors.New("test error 123"), 503)

	data, err1 := json.Marshal(errWrapped)
	c.Assert(err1, check.IsNil)

	errJson, err2 := simplejson.NewJson(data)
	c.Assert(err2, check.IsNil)

	msg, err3 := errJson.Get("errors").GetIndex(0).Get("_all").GetIndex(0).String()
	c.Assert(err3, check.IsNil)
	c.Assert(msg, check.Equals, "test error 123")
}

func (s *S) TestValidationError(c *check.C) {
	err := &ValidationError{}
	c.Assert(err.StatusCode(), check.Equals, 422)
	c.Assert(err.Error(), check.Equals, "")
	c.Assert(err.Length(), check.Equals, 0)

	err.Put("name", "is required")
	c.Assert(err.Error(), check.Equals, "name: is required")
	c.Assert(err.Length(), check.Equals, 1)
}

func (s *S) TestValidationErrorMarshallJSON(c *check.C) {
	err := &ValidationError{}
	err.Put("name", "is required")

	data, err1 := json.Marshal(err)
	c.Assert(err1, check.IsNil)

	errJson, err2 := simplejson.NewJson(data)
	c.Assert(err2, check.IsNil)

	msg, err3 := errJson.Get("errors").GetIndex(0).Get("name").GetIndex(0).String()
	c.Assert(err3, check.IsNil)
	c.Assert(msg, check.Equals, "is required")
}
