package server

import (
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
	"github.com/backstage/beat/transaction"
	"github.com/dimfeld/httptreemux"
	"io"
	"log"
	"net/http"
)

var (
	ErrResourceNotAnObject = errors.New("Json root not is an object", http.StatusBadRequest)
	ErrEmptyResource       = errors.New("Empty resource", http.StatusBadRequest)
)

func (s *Server) collectionHandle(handler transaction.TransactionHandler) httptreemux.HandlerFunc {
	return transaction.CollectionHandle(func(t *transaction.Transaction) {
		itemSchema, err := s.DB.FindItemSchemaByCollectionName(t.CollectionName)

		if err != nil {
			t.WriteError(err)
			return
		}

		t.ItemSchema = itemSchema
		handler(t)
	})
}

func (s *Server) createResource(t *transaction.Transaction) {
	resource, err := schemas.NewCollectionSchemaFromReader(t.Req.Body)

	if err == io.EOF {
		t.WriteError(ErrEmptyResource)
		return
	} else if err != nil {
		t.WriteError(errors.Newf(http.StatusBadRequest, "Invalid json: %s", err.Error()))
		return
	}

	// if _, err = resource.Map(); err != nil {
	// 	t.WriteError(ErrResourceNotAnObject)
	// 	return
	// }
	if err = s.DB.CreateResource(t.CollectionName, resource); err != nil {
		t.WriteError(errors.Newf(http.StatusInternalServerError, "Could not save to database", err.Error()))
	}
	t.WriteResultWithStatusCode(http.StatusCreated, resource)
}

func (s *Server) findResource(t *transaction.Transaction) {
	result, err := s.DB.FindCollectionSchema(t.ItemSchema.CollectionName, nil)
	if err != nil {
		println("error trying to recover")
		log.Println(err)
		return
	}
	println("returning result", result)
	t.WriteResult(result)
}

func (s *Server) findOneResource(t *transaction.Transaction) {
	t.WriteError(errors.New("TODO: findOne resource", http.StatusNotImplemented))
}

func (s *Server) findResourceById(t *transaction.Transaction) {
	t.WriteError(errors.New("TODO: find resource by id", http.StatusNotImplemented))
}

func (s *Server) deleteResourceById(t *transaction.Transaction) {
	t.WriteError(errors.New("TODO: delete resource by id", http.StatusNotImplemented))
}
