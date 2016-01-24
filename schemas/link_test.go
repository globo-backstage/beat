package schemas

import (
	"gopkg.in/check.v1"
)

func (s *S) TestBuildDefaultLinks(c *check.C) {
	links := BuildDefaultLinks("backstage-users")
	c.Assert(links, check.DeepEquals, Links{
		&Link{Rel: "self", Href: "/backstage-users/{id}"},
		&Link{Rel: "item", Href: "/backstage-users/{id}"},
		&Link{Rel: "create", Method: "POST", Href: "/backstage-users"},
		&Link{Rel: "update", Method: "PUT", Href: "/backstage-users/{id}"},
		&Link{Rel: "delete", Method: "DELETE", Href: "/backstage-users/{id}"},
		&Link{Rel: "parent", Href: "/backstage-users"},
	})
}
