package server

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/backstage/beat/auth"
	"github.com/backstage/beat/db"
	"gopkg.in/check.v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/backstage/beat/auth/static"
	_ "github.com/backstage/beat/db/mongo"
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
	logrus.SetOutput(ioutil.Discard)
}

func (s *S) TestNewWithConfigurableSettingsWithInvalidDatabase(c *check.C) {
	os.Setenv("DATABASE", "not-found")
	os.Setenv("AUTHENTICATION", "static")

	server, err := NewWithConfigurableSettings()
	c.Assert(server, check.IsNil)
	c.Assert(err, check.FitsTypeOf, db.ErrNotFound{})
}

func (s *S) TestNewWithConfigurableSettingsWithInvalidAuthentication(c *check.C) {
	os.Setenv("DATABASE", "mongo")
	os.Setenv("AUTHENTICATION", "not-found")

	server, err := NewWithConfigurableSettings()
	c.Assert(server, check.IsNil)
	c.Assert(err, check.FitsTypeOf, auth.ErrNotFound{})
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
