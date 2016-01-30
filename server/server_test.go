package server

import (
	"fmt"
	"gopkg.in/check.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var _ = check.Suite(&S{})

type S struct {
	server *Server
}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) SetUpSuite(c *check.C) {
	s.server = New(nil, nil)
}

func (s *S) TestHealthcheckInRootPath(c *check.C) {
	response := s.SimpleRequest("GET", "/")
	c.Assert(response.Code, check.Equals, http.StatusOK)
	c.Assert(response.Body.String(), check.Equals, "WORKING")
}

func (s *S) TestHealthcheck(c *check.C) {
	response := s.SimpleRequest("GET", "/healthcheck")
	c.Assert(response.Code, check.Equals, http.StatusOK)
	c.Assert(response.Body.String(), check.Equals, "WORKING")
}

func (s *S) SimpleRequest(method, path string) *httptest.ResponseRecorder {
	r, err := http.NewRequest(method, fmt.Sprintf("http://localhost%s", path), nil)
	if err != nil {
		log.Fatal(err)
	}
	return s.Request(r)
}

func (s *S) Request(r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	s.server.router.ServeHTTP(w, r)
	return w
}
