package schemas

import "fmt"

type CollectionSchema struct {
	Schema         string   `json:"$schema"`
	CollectionName string   `json:"collectionName"`
	Type           string   `json:"type"`
	Title          string   `json:"title,omitempty"`
	Links          *Links   `json:"links,omitempty"`
	Properties     colProps `json:"properties"`
}

func NewCollectionSchema(itemSchema *ItemSchema) *CollectionSchema {
	return &CollectionSchema{
		Schema:         itemSchema.Schema,
		CollectionName: itemSchema.CollectionName,
		Type:           "object",
		Title:          itemSchema.CollectionTitle,
		Links:          itemSchema.CollectionLinks,
		Properties:     colProps{itemSchema.CollectionName},
	}
}

type colProps struct {
	ref string
}

func (c *colProps) MarshalJSON() ([]byte, error) {
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
