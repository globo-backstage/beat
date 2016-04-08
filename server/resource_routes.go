package server

import (
	"github.com/backstage/beat/db"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/transaction"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/dimfeld/httptreemux"
	"github.com/satori/go.uuid"
	"github.com/xeipuuv/gojsonschema"
	"io"
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
	resource, err := simplejson.NewFromReader(t.Req.Body)

	if err == io.EOF {
		t.WriteError(ErrEmptyResource)
		return
	} else if err != nil {
		t.WriteError(errors.Newf(http.StatusBadRequest, "Invalid json: %s", err.Error()))
		return
	}

	if _, err = resource.Map(); err != nil {
		t.WriteError(ErrResourceNotAnObject)
		return
	}
	schemaLoader := gojsonschema.NewGoLoader(t.ItemSchema)
	resourceLoader := gojsonschema.NewGoLoader(resource.MustMap())

	result, err := gojsonschema.Validate(schemaLoader, resourceLoader)
	if err != nil {
		t.WriteError(errors.Wraps(err, http.StatusInternalServerError))
		return
	}

	if !result.Valid() {
		validationError := &errors.ValidationError{}
		for _, fieldErr := range result.Errors() {
			validationError.Put(fieldErr.Field(), fieldErr.Description())
		}
		t.WriteError(validationError)
		return
	}

	primaryKey := resource.Get("id").MustString()
	if primaryKey == "" {
		resource.Set("id", uuid.NewV4().String())
	}
	if cErr := transaction.RunBeforeSave(t); cErr != nil {
		t.WriteError(cErr)
		return
	}
	if cErr := s.DB.CreateResource(t.CollectionName, resource); cErr != nil {
		t.WriteError(cErr)
		return
	}
	t.WriteResultWithStatusCode(http.StatusCreated, resource)
}

func (s *Server) findResource(t *transaction.Transaction) {
	filter, err := db.NewFilterFromQueryString(t.Req.URL.RawQuery)
	result, err := s.DB.FindResources(t.ItemSchema.CollectionName, filter)
	if err != nil {
		t.WriteError(errors.New("Error feching from database", http.StatusInternalServerError))
		return
	}
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
