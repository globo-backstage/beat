package server

import (
	"fmt"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/transaction"
	simplejson "github.com/bitly/go-simplejson"
	"io"
	"net/http"
)

var (
	ErrResourceNotAnObject = errors.New("Json root not is an object", http.StatusBadRequest)
	ErrEmptyResource       = errors.New("Empty resource", http.StatusBadRequest)
)

func (s *Server) createResource(t *transaction.Transaction) {
	resource, err := simplejson.NewFromReader(t.Req.Body)

	if err == io.EOF {
		t.WriteError(ErrEmptyResource)
		return
	} else if err != nil {
		t.WriteError(errors.Newf(http.StatusBadRequest, "Invalid json: %s", err.Error()))
		return
	}
	_, err = resource.Map()

	if err != nil {
		t.WriteError(ErrResourceNotAnObject)
		return
	}

	t.WriteResultWithStatusCode(http.StatusCreated, resource)
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
