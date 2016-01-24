package schemas

import (
	"gopkg.in/check.v1"
	//"strings"
	//"testing"
)

func (s *S) TestBuildDefaultLinks(c *check.C) {
	links := BuildDefaultLinks("backstage-users")
	c.Assert(links, check.DeepEquals, Links{
		Link{Rel: "self", Href: "/api/backstage-users/{id}"},
		Link{Rel: "item", Href: "/api/backstage-users/{id}"},
		Link{Rel: "create", Method: "POST", Href: "/api/backstage-users"},
		Link{Rel: "update", Method: "PUT", Href: "/api/backstage-users/{id}"},
		Link{Rel: "delete", Method: "DELETE", Href: "/api/backstage-users/{id}"},
		Link{Rel: "parent", Href: "/api/backstage-users"},
	})
}
