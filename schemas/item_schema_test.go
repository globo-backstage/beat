package schemas

import (
	"gopkg.in/check.v1"
	"strings"
	"testing"
)

var _ = check.Suite(&S{})

type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) TestNewItemSchemaFromReader(c *check.C) {
	schema := `{
		"collectionName": "example-my-schema",
		"$schema": "http://json-schema.org/draft-03/hyper-schema#",
		"globalCollectionName": true,
		"aditionalProperties": true,
		"type": "object",
		"properties": {
			"name": {
				"type": "string"
			}
		}
	}`
	reader := strings.NewReader(schema)
	itemSchema, err := NewItemSchemaFromReader(reader)

	c.Assert(err, check.IsNil)
	c.Assert(itemSchema.CollectionName, check.Equals, "example-my-schema")
	c.Assert(itemSchema.Schema, check.Equals, "http://json-schema.org/draft-03/hyper-schema#")
	c.Assert(*itemSchema.AditionalProperties, check.Equals, true)
	c.Assert(itemSchema.Type, check.Equals, "object")
	c.Assert(itemSchema.Properties["name"]["type"], check.Equals, "string")
}
