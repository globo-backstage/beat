package schemas

import (
	"encoding/json"
	"fmt"
	"github.com/backstage/beat/errors"
	"io"
	"regexp"
)

const ItemSchemaCollectionName = "item-schemas"
const draft3Schema = "http://json-schema.org/draft-03/hyper-schema#"
const draft4Schema = "http://json-schema.org/draft-04/hyper-schema#"
const defaultSchema = draft4Schema

var CollectionNameRegex *regexp.Regexp = regexp.MustCompile(`^[a-z0-9-]+$`)
var CollectionNameSpaceRegex *regexp.Regexp = regexp.MustCompile(`^(\w+)-(.*)$`)

type Properties map[string]map[string]interface{}

// ItemSchema is the main struct for each collection, that describe
// the data contract and data services.
// This struct is based on json-schema specification,
// see more in: http://json-schema.org
type ItemSchema struct {
	Schema               string     `json:"$schema" bson:"%20schema"`
	CollectionName       string     `json:"collectionName" bson:"_id"`
	GlobalCollectionName bool       `json:"globalCollectionName" bson:"globalCollectionName"`
	AditionalProperties  *bool      `json:"aditionalProperties,omitempty" bson:"aditionalProperties"`
	Type                 string     `json:"type"`
	Properties           Properties `json:"properties,omitempty"`
	// used only in draft4
	Required []string `json:"required,omitempty"`
}

// NewItemSchemaFromReader return a new ItemSchema by an io.Reader.
// return a error if the buffer not is valid.
func NewItemSchemaFromReader(r io.Reader) (*ItemSchema, errors.Error) {
	itemSchema := &ItemSchema{}
	err := json.NewDecoder(r).Decode(itemSchema)
	if err != nil {
		return nil, errors.Wraps(err, 400)
	}
	itemSchema.fillDefaultValues()
	return itemSchema, itemSchema.validate()
}

func (schema *ItemSchema) fillDefaultValues() {
	if schema.Schema == "" {
		schema.Schema = defaultSchema
	}

	if schema.Type == "" {
		schema.Type = "object"
	}
}

func (schema *ItemSchema) String() string {
	return fmt.Sprintf(`<ItemSchema "%s">`, schema.CollectionName)
}

func (schema *ItemSchema) validate() errors.Error {
	validation := &errors.ValidationError{}

	if schema.Schema != draft3Schema && schema.Schema != draft4Schema {
		validation.Put("$schema", fmt.Sprintf(`must be "%s" or "%s"`, draft3Schema, draft4Schema))
	}

	if schema.Type != "object" {
		validation.Put("type", "Root type must be an object.")
	}

	isInvalidGlobalCollectionName := (!schema.GlobalCollectionName && !CollectionNameSpaceRegex.MatchString(schema.CollectionName))

	if schema.CollectionName == "" {
		validation.Put("collectionName", "must not be blank.")
	} else if isInvalidGlobalCollectionName || !CollectionNameRegex.MatchString(schema.CollectionName) {
		validation.Put("collectionName", "invalid format, use {namespace}-{name}, with characters a-z and 0-9, ex: backstage-users")
	}

	if validation.Length() > 0 {
		return validation
	}

	return nil
}
