package schemas

import (
	"gopkg.in/check.v1"
	"strings"
)

var _ = check.Suite(&CollectionSchemaSuite{})

type CollectionSchemaSuite struct {
	collectionSchema *CollectionSchema
}

func (s *CollectionSchemaSuite) SetUpTest(c *check.C) {
	schema := `{
		"$schema": "http://json-schema.org/draft-04/hyper-schema#",
		"collectionName": "backstage-users",
                "collectionTitle": "my collection Title",
		"type": "object",
                "collectionLinks": [
                    {"rel": "top10", "href": "http://github.com/jedi"},
                    {"rel": "history", "href": "/juniors"}
                ]
	}`
	reader := strings.NewReader(schema)
	itemSchema, err := NewItemSchemaFromReader(reader)
	c.Assert(err, check.IsNil)

	s.collectionSchema = NewCollectionSchema(itemSchema)
}

var (
	DefaultCollectionSchemaLinks = Links{
		&Link{Rel: "self", Href: "/backstage-users"},
		&Link{Rel: "list", Href: "/backstage-users"},
		&Link{Rel: "add", Href: "/backstage-users", Method: "POST",
			Schema: map[string]interface{}{
				"$ref": "/item-schemas/backstage-users",
			},
		},
		&Link{Rel: "previous", Href: "/backstage-users?filter[perPage]={perPage}&filter[page]={previousPage}{&paginateQs*}"},
		&Link{Rel: "next", Href: "/backstage-users?filter[perPage]={perPage}&filter[page]={nextPage}{&paginateQs*}"},
		&Link{Rel: "page", Href: "/backstage-users?filter[perPage]={perPage}&filter[page]={page}{&paginateQs*}"},
		&Link{Rel: "order", Href: "/backstage-users?filter[order]={orderAttribute}%20{orderDirection}{&orderQs*}"},
	}
)

func (s *CollectionSchemaSuite) TestNewCollectionSchemaProperties(c *check.C) {

	c.Assert(s.collectionSchema.Schema, check.Equals, "http://json-schema.org/draft-04/hyper-schema#")
	c.Assert(s.collectionSchema.CollectionName, check.Equals, "backstage-users")
	c.Assert(s.collectionSchema.Title, check.Equals, "my collection Title")
	c.Assert(s.collectionSchema.Type, check.Equals, "object")
}

func (s *CollectionSchemaSuite) TestNewCollectionSchemaLinks(c *check.C) {
	for i, expectedLink := range DefaultCollectionSchemaLinks {
		link := *(*s.collectionSchema.Links)[i]
		c.Assert(link, check.DeepEquals, *expectedLink)
	}

	lenDefaultLinks := len(DefaultCollectionSchemaLinks)

	link := *(*s.collectionSchema.Links)[lenDefaultLinks]
	c.Assert(link, check.DeepEquals, Link{Rel: "top10", Href: "http://github.com/jedi"})

	link = *(*s.collectionSchema.Links)[lenDefaultLinks+1]
	c.Assert(link, check.DeepEquals, Link{Rel: "history", Href: "/juniors"})

}

func (s *CollectionSchemaSuite) TestNewCollectionSchemaApplyBaseUrl(c *check.C) {
	s.collectionSchema.ApplyBaseUrl("https://my-beat.com/api")

	c.Assert(s.collectionSchema.Properties.ref, check.Equals, "https://my-beat.com/api/item-schemas/backstage-users")

	lenDefaultLinks := len(DefaultCollectionSchemaLinks)

	link := *(*s.collectionSchema.Links)[lenDefaultLinks]
	c.Assert(link, check.DeepEquals, Link{Rel: "top10", Href: "http://github.com/jedi"})

	link = *(*s.collectionSchema.Links)[lenDefaultLinks+1]
	c.Assert(link, check.DeepEquals, Link{Rel: "history", Href: "https://my-beat.com/api/juniors"})
}
