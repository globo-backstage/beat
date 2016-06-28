package server

import (
	"io"
	"net/http"

	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/transaction"
	simplejson "github.com/bitly/go-simplejson"
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

func (s *Server) findResource(t *transaction.Transaction) {
	t.WriteError(errors.New("TODO: Find resource", http.StatusNotImplemented))
}

func (s *Server) findOneResource(t *transaction.Transaction) {
	t.WriteError(errors.New("TODO: findOne resource", http.StatusNotImplemented))
}

func (s *Server) findResourceByID(t *transaction.Transaction) {
	t.WriteError(errors.New("TODO: find resource by id", http.StatusNotImplemented))
}

func (s *Server) deleteResourceByID(t *transaction.Transaction) {
	t.WriteError(errors.New("TODO: delete resource by id", http.StatusNotImplemented))
}
