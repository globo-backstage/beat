package server

import (
	"net/http"

	"github.com/backstage/beat/db"
	"github.com/backstage/beat/errors"
	"github.com/backstage/beat/schemas"
	"github.com/backstage/beat/transaction"
)

func (s *Server) createItemSchema(t *transaction.Transaction) {
	itemSchema, err := schemas.NewItemSchemaFromReader(t.Req.Body)

	if err != nil {
		t.WriteError(err)
		return
	}
	itemSchema.DiscardDefaultLinks()
	err = s.DB.CreateItemSchema(itemSchema)

	if err != nil {
		t.WriteError(err)
		return
	}

	itemSchema.AttachDefaultLinks(t.BaseUrl())
	t.WriteResultWithStatusCode(http.StatusCreated, itemSchema)
}

func (s *Server) listItemSchemas(t *transaction.Transaction) {
	filter, err := db.NewFilterFromQueryString(t.Req.URL.RawQuery)

	if err != nil {
		t.WriteError(errors.Wraps(err, http.StatusBadRequest))
		return
	}

	reply, findErr := s.DB.FindItemSchema(filter)
	if findErr != nil {
		t.WriteError(findErr)
		return
	}

	baseUrl := t.BaseUrl()
	for _, itemSchema := range reply.Items {
		itemSchema.AttachDefaultLinks(baseUrl)
	}

	t.WriteResult(reply)
}

func (s *Server) findItemSchema(t *transaction.Transaction) {
	t.ItemSchema.AttachDefaultLinks(t.BaseUrl())
	t.WriteResult(t.ItemSchema)
}

func (s *Server) findOneItemSchema(t *transaction.Transaction) {
	filter, err := db.NewFilterFromQueryString(t.Req.URL.RawQuery)

	if err != nil {
		t.WriteError(errors.Wraps(err, http.StatusBadRequest))
		return
	}

	itemSchema, findErr := s.DB.FindOneItemSchema(filter)
	if findErr != nil {
		t.WriteError(findErr)
		return
	}

	itemSchema.AttachDefaultLinks(t.BaseUrl())
	t.WriteResult(itemSchema)
}

func (s *Server) deleteItemSchema(t *transaction.Transaction) {
	collectionName := t.Params["collectionName"]
	err := s.DB.DeleteItemSchemaByCollectionName(collectionName)
	if err != nil {
		t.WriteError(err)
		return
	}

	t.NoResultWithStatusCode(http.StatusNoContent)
}

func (s *Server) updateItemSchema(t *transaction.Transaction) {
}

func (s *Server) findCollectionSchema(t *transaction.Transaction) {
	collectionSchema := schemas.NewCollectionSchema(t.ItemSchema)
	collectionSchema.ApplyBaseUrl(t.BaseUrl())

	t.WriteResult(collectionSchema)
}
