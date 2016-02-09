package server

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/backstage/beat/auth"
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/transaction"
	"github.com/dimfeld/httptreemux"
	"github.com/spf13/viper"
	"net/http"
)

type Server struct {
	Authentication auth.Authable
	DB             db.Database
	router         *httptreemux.TreeMux
}

func init() {
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", 3000)
	viper.SetDefault("database", "mongo")
	viper.SetDefault("authentication", "static")
}

func New(authentication auth.Authable, db db.Database) *Server {
	server := &Server{
		Authentication: authentication,
		DB:             db,
	}
	server.initRoutes()
	return server
}

func NewWithConfigurableSettings() (*Server, error) {
	db, err := db.New(viper.GetString("database"))

	if err != nil {
		return nil, err
	}

	auth, err := auth.New(viper.GetString("authentication"))

	if err != nil {
		return nil, err
	}

	return New(auth, db), nil
}

func (s *Server) Run() {
	bind := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	log.Infof("Backstage Beat is running on http://%s/", bind)
	log.Fatal(http.ListenAndServe(bind, s.router))
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

	s.router.GET("/api/collection-schemas/:collectionName", transaction.Handle(s.findCollectionSchemaByCollectionName))

	s.router.POST("/api/:collectionName", s.collectionHandle(s.createResource))
	s.router.GET("/api/:collectionName", s.collectionHandle(s.findResource))
	s.router.GET("/api/:collectionName/findOne", s.collectionHandle(s.findOneResource))
	s.router.GET("/api/:collectionName/:resourceId", s.collectionHandle(s.findResourceById))
	s.router.DELETE("/api/:collectionName/:resourceId", s.collectionHandle(s.deleteResourceById))
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintf(w, "WORKING")
}
