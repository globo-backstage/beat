package schemas

import (
	"fmt"
)

type CollectionSchema struct {
	Schema         string   `json:"$schema" bson:"%20schema"`
	CollectionName string   `json:"collectionName" bson:"_id"`
	Type           string   `json:"type"`
	Title          string   `json:"title,omitempty"`
	Properties     colProps `json:"properties"`
	Links          *Links   `json:"links,omitempty"`
}

func NewCollectionSchema(itemSchema *ItemSchema) *CollectionSchema {
	collectionSchema := &CollectionSchema{
		Schema:         itemSchema.Schema,
		CollectionName: itemSchema.CollectionName,
		Type:           "object",
		Title:          itemSchema.CollectionTitle,
		Links:          itemSchema.CollectionLinks,
		Properties:     colProps{itemSchema.url()},
	}

	customLinks := itemSchema.CollectionLinks
	collectionSchema.Links = collectionSchema.defaultLinks(itemSchema)

	if customLinks != nil {
		collectionSchema.Links = collectionSchema.Links.ConcatenateLinks(customLinks)
	}

	return collectionSchema
}

func (schema *CollectionSchema) ApplyBaseUrl(baseUrl string) {
	schema.Properties.ref = baseUrl + schema.Properties.ref
	schema.Links.ApplyBaseUrl(baseUrl)
}

func (schema *CollectionSchema) defaultLinks(itemSchema *ItemSchema) *Links {
	collectionUrl := itemSchema.collectionUrl()
	itemSchemaUrl := itemSchema.url()

	return &Links{
		&Link{Rel: "self", Href: collectionUrl},
		&Link{Rel: "list", Href: collectionUrl},
		&Link{Rel: "add", Method: "POST", Href: collectionUrl,
			Schema: map[string]interface{}{
				"$ref": itemSchemaUrl,
			},
		},
		&Link{
			Rel:  "previous",
			Href: fmt.Sprintf("%s?filter[perPage]={perPage}&filter[page]={previousPage}{&paginateQs*}", collectionUrl),
		},
		&Link{
			Rel:  "next",
			Href: fmt.Sprintf("%s?filter[perPage]={perPage}&filter[page]={nextPage}{&paginateQs*}", collectionUrl),
		},
		&Link{
			Rel:  "page",
			Href: fmt.Sprintf("%s?filter[perPage]={perPage}&filter[page]={page}{&paginateQs*}", collectionUrl),
		},
		&Link{
			Rel:  "order",
			Href: fmt.Sprintf("%s?filter[order]={orderAttribute}%s{orderDirection}{&orderQs*}", collectionUrl, "%20"),
		},
	}
}

type colProps struct {
	ref string
}

func (c colProps) MarshalJSON() ([]byte, error) {
	data := fmt.Sprintf(`{
    "items": {
      "items": {
        "$ref": "%s"
      },
      "type": "array"
    },
    "limit": {
      "type": "integer"
    },
    "previousOffset": {
      "type": "integer"
    },
    "nextOffset": {
      "type": "integer"
    },
    "perPage": {
      "type": "integer"
    },
    "previousPage": {
      "type": "integer"
    },
    "nextPage": {
      "type": "integer"
    },
    "itemCount": {
      "type": "integer"
    },
    "paginateQs": {
      "type": "object"
    },
    "orderQs": {
      "type": "object"
    }
}`, c.ref)
	return []byte(data), nil
}
