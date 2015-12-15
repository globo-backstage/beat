package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/backstage/beat/auth"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Authentication auth.Authable
	router         *httprouter.Router
}

func New(authentication auth.Authable) *Server {
	server := &Server{
		Authentication: authentication,
	}
	server.initRoutes()
	return server
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":3000", s.router))
}

func (s *Server) initRoutes() {
	s.router = httprouter.New()
	s.router.GET("/", s.healthCheck)
	s.router.GET("/healthcheck", s.healthCheck)
	s.router.POST("/api/:collectionName", s.createResource)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "WORKING")
}

func (s *Server) createResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collectionName := ps.ByName("collectionName")
	if collectionName == "item-schemas" {
		s.createItemSchema(w, r, ps)
		return
	}

	fmt.Fprintf(w, "Created")
}

func (s *Server) createItemSchema(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Created Item schema ")
}
