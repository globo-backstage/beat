package server

import (
	"fmt"
	"github.com/backstage/beat/auth"
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/transaction"
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
)

type Server struct {
	Authentication auth.Authable
	DB             db.Database
	router         *httptreemux.TreeMux
}

func New(authentication auth.Authable, db db.Database) *Server {
	server := &Server{
		Authentication: authentication,
		DB:             db,
	}
	server.initRoutes()
	return server
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":3000", s.router))
}

func (s *Server) initRoutes() {
	s.router = httptreemux.New()
	s.router.GET("/", s.healthCheck)
	s.router.GET("/healthcheck", s.healthCheck)

	s.router.POST("/api/item-schemas", transaction.Handle(s.createItemSchema))
	s.router.GET("/api/item-schemas", transaction.Handle(s.findItemSchema))
	s.router.GET("/api/item-schemas/findOne", transaction.Handle(s.findOneItemSchema))
	s.router.GET("/api/item-schemas/:collectionName", transaction.Handle(s.findItemSchemaByCollectionName))
	s.router.DELETE("/api/item-schemas/:collectionName", transaction.Handle(s.deleteItemSchemaByCollectionName))

	s.router.POST("/api/:collectionName", s.createResource)
	s.router.GET("/api/:collectionName", s.findResource)
	s.router.GET("/api/:collectionName/findOne", s.findOneResource)
	s.router.GET("/api/:collectionName/:resourceId", s.findResourceById)
	s.router.DELETE("/api/:collectionName/:resourceId", s.deleteResourceById)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "WORKING")
}
