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

func (s *S) TestNewItemSchemaWhenOmmitAditionalProperties(c *check.C) {
	schema := `{
		"collectionName": "example-my-schema",
		"$schema": "http://json-schema.org/draft-03/hyper-schema#"
	}`
	reader := strings.NewReader(schema)
	itemSchema, err := NewItemSchemaFromReader(reader)

	c.Assert(err, check.IsNil)
	c.Assert(itemSchema.AditionalProperties, check.IsNil)
}

func (s *S) TestNewItemSchemaWithDefaultValues(c *check.C) {
	schema := `{
		"collectionName": "example-my-schema"
	}`
	reader := strings.NewReader(schema)
	itemSchema, err := NewItemSchemaFromReader(reader)

	c.Assert(err, check.IsNil)
	c.Assert(itemSchema.Schema, check.Equals, "http://json-schema.org/draft-04/hyper-schema#")
}

func (s *S) TestNewItemSchemaWithInvalidSchema(c *check.C) {
	schema := `{
		"$schema": "http://globo.com/invalid-schema"
	}`

	reader := strings.NewReader(schema)
	_, err := NewItemSchemaFromReader(reader)

	c.Assert(err, check.Not(check.IsNil))
	c.Assert(err.Error(), check.Equals, `$schema must be "http://json-schema.org/draft-03/hyper-schema#" or "http://json-schema.org/draft-04/hyper-schema#"`)
}
