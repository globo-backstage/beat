package server

import (
	"fmt"
	"github.com/backstage/beat/auth"
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
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

func (s *Server) createItemSchema(t *transaction.Transaction) {
	itemSchema, err := schemas.NewItemSchemaFromReader(t.Req.Body)

	if err != nil {
		t.WriteError(err)
		return
	}

	err = s.DB.CreateItemSchema(itemSchema)

	if err != nil {
		t.WriteError(err)
		return
	}

	t.WriteResultWithStatusCode(http.StatusCreated, itemSchema)
}

func (s *Server) findItemSchema(t *transaction.Transaction) {
	filter, err := db.NewFilterFromQueryString(t.Req.URL.RawQuery)

	if err != nil {
		t.WriteError(errors.Wraps(err, 400))
		return
	}

	reply, findErr := s.DB.FindItemSchema(filter)
	if findErr != nil {
		t.WriteError(findErr)
		return
	}

	t.WriteResult(reply)
}

func (s *Server) findItemSchemaByCollectionName(t *transaction.Transaction) {
	collectionName := t.Params["collectionName"]
	itemSchema, err := s.DB.FindItemSchemaByCollectionName(collectionName)
	if err != nil {
		t.WriteError(err)
		return
	}

	t.WriteResult(itemSchema)
}

func (s *Server) findOneItemSchema(t *transaction.Transaction) {
	filter, err := db.NewFilterFromQueryString(t.Req.URL.RawQuery)

	if err != nil {
		t.WriteError(errors.Wraps(err, 400))
		return
	}

	itemSchema, findErr := s.DB.FindOneItemSchema(filter)
	if findErr != nil {
		t.WriteError(findErr)
		return
	}

	t.WriteResult(itemSchema)
}

func (s *Server) deleteItemSchemaByCollectionName(t *transaction.Transaction) {
	collectionName := t.Params["collectionName"]
	err := s.DB.DeleteItemSchemaByCollectionName(collectionName)
	if err != nil {
		t.WriteError(err)
		return
	}

	t.NoResultWithStatusCode(http.StatusNoContent)
}
