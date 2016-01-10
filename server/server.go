package server

import (
	"encoding/json"
	"fmt"
	"github.com/backstage/beat/auth"
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Server struct {
	Authentication auth.Authable
	DB             db.Database
	router         *httprouter.Router
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
	s.router = httprouter.New()
	s.router.GET("/", s.healthCheck)
	s.router.GET("/healthcheck", s.healthCheck)
	s.router.POST("/api/:collectionName", s.createResource)
	s.router.GET("/api/:collectionName", s.findResource)
	s.router.GET("/api/:collectionName/:resourceId", s.findResourceById)
	s.router.DELETE("/api/:collectionName/:resourceId", s.deleteResourceById)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "WORKING")
}

func (s *Server) createResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collectionName := ps.ByName("collectionName")
	if collectionName == schemas.ItemSchemaCollectionName {
		s.createItemSchema(w, r, ps)
		return
	}

	fmt.Fprintf(w, "Created")
}

func (s *Server) findResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collectionName := ps.ByName("collectionName")
	if collectionName == schemas.ItemSchemaCollectionName {
		s.findItemSchema(w, r, ps)
		return
	}

	fmt.Fprintf(w, "Find")
}

func (s *Server) findResourceById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collectionName := ps.ByName("collectionName")

	if collectionName == schemas.ItemSchemaCollectionName {
		s.findItemSchemaByCollectionName(w, r, ps)
		return
	}

	fmt.Fprintf(w, "FindById")
}

func (s *Server) deleteResourceById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	collectionName := ps.ByName("collectionName")

	if collectionName == schemas.ItemSchemaCollectionName {
		s.deleteItemSchemaByCollectionName(w, r, ps)
		return
	}

	fmt.Fprintf(w, "deleteById")
}

func (s *Server) createItemSchema(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (s *Server) findItemSchema(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (s *Server) findItemSchemaByCollectionName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resourceId := ps.ByName("resourceId")

	if resourceId == "findOne" {
		s.findOneItemSchema(w, r, ps)
		return
	}

	itemSchema, err := s.DB.FindItemSchemaByCollectionName(resourceId)
	if err != nil {
		s.writeError(w, err)
		return
	}

	json.NewEncoder(w).Encode(itemSchema)
}

func (s *Server) findOneItemSchema(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func (s *Server) deleteItemSchemaByCollectionName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resourceId := ps.ByName("resourceId")
	err := s.DB.DeleteItemSchemaByCollectionName(resourceId)
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
