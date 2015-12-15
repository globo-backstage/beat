package schemas

import (
	"encoding/json"
	"io"
)

type Properties map[string]map[string]interface{}

type ItemSchema struct {
	Schema               string     `json:"$schema" bson:"%20schema"`
	CollectionName       string     `json:"collectionName" bson:"_id"`
	GlobalCollectionName bool       `json:"globalCollectionName" bson:"globalCollectionName"`
	AditionalProperties  *bool      `json:"aditionalProperties,omitempty"`
	VersionId            string     `json:"versionId" bson:"versionId"`
	Type                 string     `json:"type"`
	Properties           Properties `json:"properties,omitempty"`
	// used only in draft4
	Required []string `json:"required,omitempty"`
}

func NewItemSchemaFromReader(r io.Reader) (*ItemSchema, error) {
	itemSchema := &ItemSchema{}
	err := json.NewDecoder(r).Decode(itemSchema)
	if err != nil {
		return nil, err
	}
	return itemSchema, nil
}
