package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/backstage/beat/auth"
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/transaction"
	"github.com/dimfeld/httptreemux"
	"github.com/spf13/viper"
)

type Server struct {
	*httptreemux.TreeMux
	Authentication auth.Authable
	DB             db.Database
}

type ServerOpts struct {
	Authentication auth.Authable
	DB             db.Database
}

func init() {
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", 3000)
	viper.SetDefault("database", "mongo")
	viper.SetDefault("authentication", "static")
}

func New() (*Server, error) {
	db, err := db.New(viper.GetString("database"))

	if err != nil {
		return nil, err
	}

	auth, err := auth.New(viper.GetString("authentication"))

	if err != nil {
		return nil, err
	}
	return NewWithOpts(&ServerOpts{
		Authentication: auth,
		DB:             db,
	}), nil

}

func NewWithOpts(opts *ServerOpts) *Server {
	router := httptreemux.New()
	server := &Server{
		TreeMux:        router,
		Authentication: opts.Authentication,
		DB:             opts.DB,
	}
	server.initRoutes()
	return server
}

func (s *Server) Run() {
	bind := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	log.Infof("Backstage Beat is running on http://%s/", bind)
	log.Fatal(http.ListenAndServe(bind, s))
}

func (s *Server) initRoutes() {
	s.GET("/", s.healthCheck)
	s.GET("/healthcheck", s.healthCheck)

	s.POST("/api/item-schemas", transaction.Handle(s.createItemSchema))
	s.GET("/api/item-schemas", transaction.Handle(s.findItemSchema))
	s.GET("/api/item-schemas/findOne", transaction.Handle(s.findOneItemSchema))
	s.GET("/api/item-schemas/:collectionName", transaction.Handle(s.findItemSchemaByCollectionName))
	s.DELETE("/api/item-schemas/:collectionName", transaction.Handle(s.deleteItemSchemaByCollectionName))

	s.GET("/api/collection-schemas/:collectionName", transaction.Handle(s.findCollectionSchemaByCollectionName))

	s.POST("/api/:collectionName", s.collectionHandle(s.createResource))
	s.GET("/api/:collectionName", s.collectionHandle(s.findResource))
	s.GET("/api/:collectionName/findOne", s.collectionHandle(s.findOneResource))
	s.GET("/api/:collectionName/:resourceId", s.collectionHandle(s.findResourceById))
	s.DELETE("/api/:collectionName/:resourceId", s.collectionHandle(s.deleteResourceById))
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "WORKING")
}
