package schemas

import (
	"encoding/json"
	"fmt"
	"github.com/backstage/beat/errors"
	"io"
	"net/http"
	"regexp"
)

const ItemSchemaPrimaryKey = "collectionName"
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
	Title                string     `json:"title,omitempty"`
	CollectionTitle      string     `json:"collectionTitle,omitempty"`
	GlobalCollectionName bool       `json:"globalCollectionName,omitempty"" bson:"globalCollectionName"`
	AditionalProperties  *bool      `json:"aditionalProperties,omitempty" bson:"aditionalProperties"`
	Type                 string     `json:"type"`
	Properties           Properties `json:"properties,omitempty"`
	Required             []string   `json:"required,omitempty"` // used only in draft4
	Links                *Links     `json:"links,omitempty"`
	CollectionLinks      *Links     `json:"collectionLinks,omitempty"`
}

// NewItemSchemaFromReader return a new ItemSchema by an io.Reader.
// return a error if the buffer not is valid.
func NewItemSchemaFromReader(r io.Reader) (*ItemSchema, errors.Error) {
	itemSchema := &ItemSchema{}
	err := json.NewDecoder(r).Decode(itemSchema)
	if err != nil {
		return nil, errors.Wraps(err, http.StatusBadRequest)
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

func (schema *ItemSchema) AttachDefaultLinks(baseUrl string) {
	customLinks := schema.Links
	schema.Links = schema.defaultLinks()

	if customLinks != nil {
		schema.Links = schema.Links.ConcatenateLinks(customLinks)
	}
	schema.Links.ApplyBaseUrl(baseUrl)
}

func (schema *ItemSchema) DiscardDefaultLinks() {
	if schema.Links != nil {
		schema.Links = schema.Links.DiscardDefaultLinks()
	}
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

func (schema *ItemSchema) collectionUrl() string {
	return fmt.Sprintf("/%s", schema.CollectionName)
}

func (schema *ItemSchema) url() string {
	return fmt.Sprintf("/%s/%s", ItemSchemaCollectionName, schema.CollectionName)
}

func (schema *ItemSchema) defaultLinks() *Links {
	collectionUrl := schema.collectionUrl()
	schemaUrl := schema.url()
	itemUrl := fmt.Sprintf("/%s/{id}", schema.CollectionName)

	return &Links{
		&Link{Rel: "self", Href: itemUrl},
		&Link{Rel: "item", Href: itemUrl},
		&Link{Rel: "create", Method: "POST", Href: collectionUrl,
			Schema: map[string]interface{}{
				"$ref": schemaUrl,
			},
		},
		&Link{Rel: "update", Method: "PUT", Href: itemUrl},
		&Link{Rel: "delete", Method: "DELETE", Href: itemUrl},
		&Link{Rel: "parent", Href: collectionUrl},
	}
}
