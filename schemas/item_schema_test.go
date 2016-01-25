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
		"$schema": "http://globo.com/invalid-schema",
		"collectionName" : "backstage-valid"
	}`

	_, err := NewItemSchemaFromReader(strings.NewReader(schema))

	c.Assert(err, check.Not(check.IsNil))
	c.Assert(err.Error(), check.Equals, `$schema: must be "http://json-schema.org/draft-03/hyper-schema#" or "http://json-schema.org/draft-04/hyper-schema#"`)

	schema = `{
		"collectionName": "backstage-users",
		"type": "array"
	}`
	_, err = NewItemSchemaFromReader(strings.NewReader(schema))

	c.Assert(err, check.Not(check.IsNil))
	c.Assert(err.Error(), check.Equals, "type: Root type must be an object.")
	c.Assert(err.StatusCode(), check.Equals, 422)

	schema = `{}`
	_, err = NewItemSchemaFromReader(strings.NewReader(schema))

	c.Assert(err, check.Not(check.IsNil))
	c.Assert(err.Error(), check.Equals, "collectionName: must not be blank.")
	c.Assert(err.StatusCode(), check.Equals, 422)

	schema = `{
                "collectionName": "123$!"
        }`
	_, err = NewItemSchemaFromReader(strings.NewReader(schema))

	c.Assert(err, check.Not(check.IsNil))
	c.Assert(err.Error(), check.Equals, "collectionName: invalid format, use {namespace}-{name}, with characters a-z and 0-9, ex: backstage-users")
	c.Assert(err.StatusCode(), check.Equals, 422)
}

func (s *S) TestNewItemSchemaWithoutNameSpace(c *check.C) {
	schema := `{
                "collectionName": "users"
        }`
	_, err := NewItemSchemaFromReader(strings.NewReader(schema))

	c.Assert(err, check.Not(check.IsNil))
	c.Assert(err.Error(), check.Equals, "collectionName: invalid format, use {namespace}-{name}, with characters a-z and 0-9, ex: backstage-users")
	c.Assert(err.StatusCode(), check.Equals, 422)
}

func (s *S) TestNewItemSchemaWithGlobalCollectionName(c *check.C) {
	schema := `{
                "collectionName": "users",
                "globalCollectionName": true
        }`
	itemSchema, err := NewItemSchemaFromReader(strings.NewReader(schema))

	c.Assert(err, check.IsNil)
	c.Assert(itemSchema.GlobalCollectionName, check.Equals, true)
}

var (
	DefaultLinks = Links{
		&Link{Rel: "self", Href: "http://api.mysite.com/backstage-users/{id}"},
		&Link{Rel: "item", Href: "http://api.mysite.com/backstage-users/{id}"},
		&Link{Rel: "create", Href: "http://api.mysite.com/backstage-users", Method: "POST",
			Schema: map[string]interface{}{
				"$ref": "http://api.mysite.com/item-schemas/backstage-users",
			},
		},
		&Link{Rel: "update", Href: "http://api.mysite.com/backstage-users/{id}", Method: "PUT"},
		&Link{Rel: "delete", Href: "http://api.mysite.com/backstage-users/{id}", Method: "DELETE"},
		&Link{Rel: "parent", Href: "http://api.mysite.com/backstage-users"},
	}
)

func (s *S) TestAttachDefaultLinks(c *check.C) {
	schema := `{
                "collectionName": "backstage-users"
        }`
	itemSchema, err := NewItemSchemaFromReader(strings.NewReader(schema))
	c.Assert(err, check.IsNil)
	itemSchema.AttachDefaultLinks("http://api.mysite.com")

	for i, expectedLink := range DefaultLinks {
		link := *(*itemSchema.Links)[i]
		c.Assert(link, check.DeepEquals, *expectedLink)
	}
}

func (s *S) TestAttachDefaultLinksWithCustomLinks(c *check.C) {
	schema := `{
                "collectionName": "backstage-users",
                "links": [
                    {"rel": "permissions", "href": "/backstage-permissions/{id}"}
                ]
        }`
	itemSchema, err := NewItemSchemaFromReader(strings.NewReader(schema))
	c.Assert(err, check.IsNil)
	itemSchema.AttachDefaultLinks("http://api.mysite.com")

	lenDefaultLinks := len(DefaultLinks)
	link := *(*itemSchema.Links)[lenDefaultLinks]
	c.Assert(link, check.DeepEquals, Link{Rel: "permissions", Href: "http://api.mysite.com/backstage-permissions/{id}"})
}

func (s *S) TestAttachDefaultLinksWithCustomLinksWithAbsoluteLink(c *check.C) {
	schema := `{
                "collectionName": "backstage-users",
                "links": [
                    {"rel": "logs", "href": "http://mylog-service/by-user/{id}"}
                ]
        }`
	itemSchema, err := NewItemSchemaFromReader(strings.NewReader(schema))
	c.Assert(err, check.IsNil)
	itemSchema.AttachDefaultLinks("http://api.mysite.com")

	lenDefaultLinks := len(DefaultLinks)
	link := *(*itemSchema.Links)[lenDefaultLinks]

	c.Assert(link, check.DeepEquals, Link{Rel: "logs", Href: "http://mylog-service/by-user/{id}"})
}

func (s *S) TestAttachDefaultLinksWithCustomLinksWithTemplateLink(c *check.C) {
	schema := `{
                "collectionName": "backstage-users",
                "links": [
                    {"rel": "view", "href": "{+url}"}
                ]
        }`
	itemSchema, err := NewItemSchemaFromReader(strings.NewReader(schema))
	c.Assert(err, check.IsNil)
	itemSchema.AttachDefaultLinks("http://api.mysite.com")

	lenDefaultLinks := len(DefaultLinks)
	link := *(*itemSchema.Links)[lenDefaultLinks]

	c.Assert(link, check.DeepEquals, Link{Rel: "view", Href: "{+url}"})
}

func (s *S) TestAttachDefaultLinksWithCustomLinksWithRefSchema(c *check.C) {
	schema := `{
                "collectionName": "backstage-users",
                "links": [
                    {
                      "rel": "view", "href": "/blah",
                      "schema": {"$ref": "/api/kaka1"},
                      "targetSchema": {"$ref": "/api/kaka2"}
                    }
                ]
        }`
	itemSchema, err := NewItemSchemaFromReader(strings.NewReader(schema))
	c.Assert(err, check.IsNil)
	itemSchema.AttachDefaultLinks("http://api.mysite.com")

	lenDefaultLinks := len(DefaultLinks)
	link := *(*itemSchema.Links)[lenDefaultLinks]

	c.Assert(link, check.DeepEquals, Link{
		Rel: "view", Href: "http://api.mysite.com/blah",
		Schema: map[string]interface{}{
			"$ref": "http://api.mysite.com/api/kaka1",
		},
		TargetSchema: map[string]interface{}{
			"$ref": "http://api.mysite.com/api/kaka2",
		},
	})
}

func (s *S) TestDiscardDefaultLinks(c *check.C) {
	schema := `{
                "collectionName": "backstage-users",
                "links": [
                    {
                      "rel": "self",
                      "href": "/hacked-url"
                    }
                ]
        }`
	itemSchema, err := NewItemSchemaFromReader(strings.NewReader(schema))
	c.Assert(err, check.IsNil)
	itemSchema.DiscardDefaultLinks()
	c.Assert(*itemSchema.Links, check.HasLen, 0)
}

func (s *S) TestDiscardDefaultLinksWithCustomLinks(c *check.C) {
	schema := `{
                "collectionName": "backstage-users",
                "links": [
                    {
                      "rel": "customLink1",
                      "href": "/hacked-url1"
                    },
                    {
                      "rel": "create",
                      "href": "/api/user"
                    },
                    {
                      "rel": "customLink2",
                      "href": "/hacked-url2"
                    }
                ]
        }`
	itemSchema, err := NewItemSchemaFromReader(strings.NewReader(schema))
	c.Assert(err, check.IsNil)
	itemSchema.DiscardDefaultLinks()

	c.Assert(*itemSchema.Links, check.HasLen, 2)

	link := *(*itemSchema.Links)[0]
	c.Assert(link, check.DeepEquals, Link{Rel: "customLink1", Href: "/hacked-url1"})

	link = *(*itemSchema.Links)[1]
	c.Assert(link, check.DeepEquals, Link{Rel: "customLink2", Href: "/hacked-url2"})
}
