package server

import (
	"encoding/json"
	"fmt"
	"github.com/backstage/beat/auth"
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
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

	s.router.POST("/api/item-schemas", s.createItemSchema)
	s.router.GET("/api/item-schemas", s.findItemSchema)
	s.router.GET("/api/item-schemas/findOne", s.findOneItemSchema)
	s.router.GET("/api/item-schemas/:collectionName", s.findItemSchemaByCollectionName)
	s.router.DELETE("/api/item-schemas/:collectionName", s.deleteItemSchemaByCollectionName)

	s.router.POST("/api/:collectionName", s.createResource)
	s.router.GET("/api/:collectionName", s.findResource)
	s.router.GET("/api/:collectionName/findOne", s.findOneResource)
	s.router.GET("/api/:collectionName/:resourceId", s.findResourceById)
	s.router.DELETE("/api/:collectionName/:resourceId", s.deleteResourceById)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "WORKING")
}

func (s *Server) createResource(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: Create resource")
}

func (s *Server) findResource(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: Find resource")
}

func (s *Server) findOneResource(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: findOne resource")
}

func (s *Server) findResourceById(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: Find resource by id")
}

func (s *Server) deleteResourceById(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "TODO: delete resource By Id")
}

func (s *Server) createItemSchema(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	itemSchema, err := schemas.NewItemSchemaFromReader(r.Body)

	if err != nil {
		s.writeError(w, err)
		return
	}

	dbErr := s.DB.CreateItemSchema(itemSchema)

	if dbErr != nil {
		s.writeError(w, errors.Wraps(dbErr, 500))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(itemSchema)
}

func (s *Server) findItemSchema(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := db.NewFilterFromQueryString(r.URL.RawQuery)

	if err != nil {
		s.writeError(w, errors.Wraps(err, 400))
		return
	}

	reply, findErr := s.DB.FindItemSchema(filter)
	if findErr != nil {
		s.writeError(w, findErr)
		return
	}

	json.NewEncoder(w).Encode(reply)
}

func (s *Server) findItemSchemaByCollectionName(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	collectionName := ps["collectionName"]
	itemSchema, err := s.DB.FindItemSchemaByCollectionName(collectionName)
	if err != nil {
		s.writeError(w, err)
		return
	}

	json.NewEncoder(w).Encode(itemSchema)
}

func (s *Server) findOneItemSchema(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := db.NewFilterFromQueryString(r.URL.RawQuery)

	if err != nil {
		s.writeError(w, errors.Wraps(err, 400))
		return
	}

	itemSchema, findErr := s.DB.FindOneItemSchema(filter)
	if findErr != nil {
		s.writeError(w, findErr)
		return
	}

	json.NewEncoder(w).Encode(itemSchema)
}

func (s *Server) deleteItemSchemaByCollectionName(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	collectionName := ps["collectionName"]
	err := s.DB.DeleteItemSchemaByCollectionName(collectionName)
	if err != nil {
		s.writeError(w, err)
		return
	}

	w.WriteHeader(204)
}

func (s *Server) writeError(w http.ResponseWriter, err errors.Error) {
	w.WriteHeader(err.StatusCode())
	json.NewEncoder(w).Encode(err)
}
