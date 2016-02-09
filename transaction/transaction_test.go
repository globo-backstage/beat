package transaction

import (
	"github.com/Sirupsen/logrus"
	"github.com/backstage/beat/errors"
	simplejson "github.com/bitly/go-simplejson"
	"gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var _ = check.Suite(&S{})

type S struct {
	Writer *httptest.ResponseRecorder
	Req    *http.Request
	T      *Transaction
}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) SetUpSuite(c *check.C) {
	logrus.SetOutput(ioutil.Discard)
}

func (s *S) SetUpTest(c *check.C) {
	var err error
	s.Writer = httptest.NewRecorder()
	s.Req, err = http.NewRequest("GET", "http://localhost/path", nil)
	c.Assert(err, check.IsNil)

	s.T = &Transaction{
		Req:    s.Req,
		writer: s.Writer,
	}
}

func (s *S) TestHandle(c *check.C) {
	var capturedTransaction *Transaction

	handler := Handle(func(t *Transaction) {
		capturedTransaction = t
	})

	handler(s.Writer, s.Req, map[string]string{"collectionName": "users"})

	c.Assert(capturedTransaction.Id, check.HasLen, 22)
	c.Assert(capturedTransaction.writer, check.Equals, s.Writer)
	c.Assert(capturedTransaction.Req, check.Equals, s.Req)
	c.Assert(capturedTransaction.Params, check.DeepEquals, map[string]string{"collectionName": "users"})
}

func (s *S) TestWriteError(c *check.C) {
	s.T.WriteError(errors.New("my error", http.StatusInternalServerError))
	c.Assert(s.T.statusCode, check.Equals, http.StatusInternalServerError)
	c.Assert(s.Writer.Code, check.Equals, http.StatusInternalServerError)

	json, err := simplejson.NewFromReader(s.Writer.Body)
	c.Assert(err, check.IsNil)

	msg := json.Get("errors").GetIndex(0).Get("_all").GetIndex(0).MustString()
	c.Assert(msg, check.Equals, "my error")
}

func (s *S) TestNoResultWithStatusCode(c *check.C) {
	s.T.NoResultWithStatusCode(http.StatusCreated)
	c.Assert(s.T.statusCode, check.Equals, http.StatusCreated)
	c.Assert(s.Writer.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestWriteResult(c *check.C) {
	result := map[string]string{
		"test": "ok",
	}
	s.T.WriteResult(&result)

	c.Assert(s.T.statusCode, check.Equals, http.StatusOK)
	c.Assert(s.Writer.Code, check.Equals, http.StatusOK)

	json, err := simplejson.NewFromReader(s.Writer.Body)
	c.Assert(err, check.IsNil)
	msg := json.Get("test").MustString()
	c.Assert(msg, check.Equals, "ok")
}

func (s *S) TestWriteResultWithStatusCode(c *check.C) {
	result := map[string]string{
		"test": "with-status-code",
	}
	s.T.WriteResultWithStatusCode(http.StatusMethodNotAllowed, &result)

	c.Assert(s.T.statusCode, check.Equals, http.StatusMethodNotAllowed)
	c.Assert(s.Writer.Code, check.Equals, http.StatusMethodNotAllowed)

	json, err := simplejson.NewFromReader(s.Writer.Body)
	c.Assert(err, check.IsNil)

	msg := json.Get("test").MustString()
	c.Assert(msg, check.Equals, "with-status-code")
}

func (s *S) TestBaseUrl(c *check.C) {
	r, err := http.NewRequest("GET", "http://my-host.com/healtcheck", nil)
	c.Assert(err, check.IsNil)
	s.T.Req = r

	c.Assert(s.T.BaseUrl(), check.Equals, "http://my-host.com/api")
}

func (s *S) TestIdFromRequestWithEmptyHeader(c *check.C) {
	r, err := http.NewRequest("GET", "http://localhost", nil)
	c.Assert(err, check.IsNil)

	id := IdFromRequest(r)
	c.Assert(id, check.HasLen, 22)
}

func (s *S) TestIdFromRequestWithFilledHeader(c *check.C) {
	r, err := http.NewRequest("GET", "http://localhost", nil)
	c.Assert(err, check.IsNil)

	r.Header.Set("Backstage-Transaction", "BBBBBBBBBBBBBBBBBBBBBZ")
	id := IdFromRequest(r)
	c.Assert(id, check.Equals, "BBBBBBBBBBBBBBBBBBBBBZ")
}

func (s *S) TestIdFromRequestWithBigHeader(c *check.C) {
	r, err := http.NewRequest("GET", "http://localhost", nil)
	c.Assert(err, check.IsNil)

	r.Header.Set("Backstage-Transaction", "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
	id := IdFromRequest(r)
	c.Assert(id, check.HasLen, 22)
}
